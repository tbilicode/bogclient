package translate

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"regexp"

	"cloud.google.com/go/translate"
	"github.com/effective-security/x/fileutil"
	"github.com/effective-security/x/values"
	"github.com/effective-security/xlog"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

var logger = xlog.NewPackageLogger("github.com/tbilicode/bogclient/pkg", "translate")

type Provider interface {
	io.Closer
	Extract(doc any) []string
	Translate(ctx context.Context, provider string, texts []string) (map[string]string, error)
	// Update replaces Georgian text with translated text, and returns a map of replaced text
	Update(ctx context.Context, doc any) (map[string]string, error)
}

type Translator struct {
	Provider

	translated map[string]string
	file       string
	override   bool
	apiKey     string
}

func NewTranslator() *Translator {
	return &Translator{
		translated: make(map[string]string),
	}
}

func (t *Translator) WithAPIKey(apiKey string) *Translator {
	t.apiKey = apiKey
	return t
}

func (t *Translator) Close() error {
	if t.file != "" && t.override {
		return t.SaveDictionary(t.file)
	}
	return nil
}

func (t *Translator) SaveDictionary(file string) error {
	jsonData, err := json.MarshalIndent(t.translated, "", "  ")
	if err != nil {
		return errors.WithMessage(err, "failed to marshal data")
	}
	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		return errors.WithMessage(err, "failed to write file")
	}
	return nil
}

func (t *Translator) LoadDictionary(file string, override bool) error {
	t.file = file
	t.override = override

	if fileutil.FileExists(file) == nil {
		data, err := os.ReadFile(file)
		if err != nil {
			return errors.WithMessage(err, "failed to read file")
		}

		var res map[string]string
		err = json.Unmarshal(data, &res)
		if err != nil {
			return errors.WithMessage(err, "failed to unmarshal data")
		}
		t.translated = res
	}
	return nil
}

func (t *Translator) Extract(doc any) (map[string]string, error) {
	texts := make(map[string]string)
	err := t.extractOrUpdate(reflect.ValueOf(doc), texts, false)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse document")
	}
	return texts, nil
}

func (t *Translator) Translate(ctx context.Context, provider string, texts map[string]string) error {
	if len(texts) == 0 {
		return nil
	}
	toTranslate := make(map[string]string, len(texts))
	for k, v := range texts {
		if t.translated[k] == "" {
			toTranslate[k] = v
		}
	}

	switch provider {
	case "google":
		err := t.GoogleTranslate(ctx, toTranslate)
		if err != nil {
			return errors.WithMessage(err, "failed to translate")
		}
	case "ai":
		err := t.OpenAITranslateJSON(ctx, toTranslate)
		if err != nil {
			return errors.WithMessage(err, "failed to translate")
		}
	}

	for k, v := range toTranslate {
		t.translated[k] = v
	}
	return nil
}

func (t *Translator) Update(ctx context.Context, doc any) (map[string]string, error) {
	replaced := make(map[string]string)
	err := t.extractOrUpdate(reflect.ValueOf(doc), replaced, true)
	if err != nil {
		return replaced, errors.WithMessage(err, "failed to update document")
	}

	return replaced, nil
}

func (t *Translator) extractOrUpdate(v reflect.Value, texts map[string]string, update bool) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := t.extractOrUpdate(v.Field(i), texts, update); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := t.extractOrUpdate(v.Index(i), texts, update); err != nil {
				return err
			}
		}
	case reflect.String:
		if v.CanSet() {
			str := v.String()
			if IsGeorgian(str) {
				if update {
					if tr := t.translated[str]; tr != "" && tr != str {
						texts[str] = tr
						v.SetString(tr)
					}
				} else if texts[str] == "" && t.translated[str] == "" {
					texts[str] = ""
				}
			}
		}
	case reflect.Ptr:
		if err := t.extractOrUpdate(v.Elem(), texts, update); err != nil {
			return err
		}
	case reflect.Map:
		if v.Type().String() == "map[string]string" {
			m := v.Interface().(map[string]string)
			for _, str := range m {
				if IsGeorgian(str) {
					if update {
						if tr := t.translated[str]; tr != "" && tr != str {
							texts[str] = tr
							v.SetString(tr)
						}
					} else if texts[str] == "" && t.translated[str] == "" {
						texts[str] = ""
					}
				}
			}
		} else {
			iter := v.MapRange()
			for iter.Next() {
				if err := t.extractOrUpdate(iter.Value(), texts, update); err != nil {
					return err
				}
			}
		}
	default:
	}
	return nil
}

var georgianRegex = regexp.MustCompile("[\u10A0-\u10FF]+")

// IsGeorgian detects Georgian text
func IsGeorgian(text string) bool {
	return georgianRegex.MatchString(text)
}

func (t *Translator) GoogleTranslate(ctx context.Context, texts map[string]string) error {
	apiKey := values.Coalesce(t.apiKey, os.Getenv("BOG_GOOGLE_APIKEY"))
	if apiKey == "" {
		return errors.New("missing GOOGLE API KEY")
	}

	keys := make([]string, 0, len(texts))
	for k := range texts {
		keys = append(keys, k)
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return errors.WithMessage(err, "failed to create client")
	}
	defer client.Close()

	sourceLang, _ := language.Parse("ka") // Georgian
	targetLang, _ := language.Parse("en")
	resp, err := client.Translate(ctx, keys, targetLang, &translate.Options{
		Format: "text", // Ensures text mode (not HTML)
		Source: sourceLang,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to translate")
	}

	for i, t := range resp {
		key := keys[i]
		texts[key] = t.Text
	}

	return nil
}
