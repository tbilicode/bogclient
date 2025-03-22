package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/effective-security/porto/xhttp/header"
	"github.com/effective-security/x/values"
	"github.com/effective-security/xlog"
	"github.com/pkg/errors"
	"github.com/spaolacci/murmur3"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

// OpenAI API request struct
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI API response struct
type OpenAIResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

// OpenAI API error response
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code,omitempty"`
}

// OpenAITranslateJSON returns translated JSON using OpenAI
func (t *Translator) OpenAITranslateJSON(ctx context.Context, texts map[string]string) error {
	// TODO: chunks
	return t.openAITranslateJSON(ctx, texts)
}

func tag(key string) uint64 {
	h := murmur3.New64()
	h.Write([]byte(key))
	return uint64(h.Sum64())
}

func (t *Translator) openAITranslateJSON(ctx context.Context, texts map[string]string) error {
	apiKey := values.Coalesce(t.apiKey, os.Getenv("BOG_OPENAI_APIKEY"))
	if apiKey == "" {
		return errors.New("missing OpenAI API key, use BOG_OPENAI_APIKEY environment variable")
	}

	vals := values.MapAny{}
	for k := range texts {
		id := tag(k)
		vals[strconv.FormatUint(id, 10)] = k
	}

	toTranslate := vals.JSON()

	// Define the prompt to instruct GPT
	prompt := "Translate the following JSON from Georgian to English while preserving the structure and only translating values (not keys)." +
		" Return only the translated JSON without any extra text.\n```json\n" + toTranslate + "\n```"

	requestBody, _ := json.Marshal(map[string]any{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional translator specialized in accounting and finance."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	})

	req, _ := http.NewRequest("POST", openAIURL, bytes.NewBuffer(requestBody))
	req.Header.Set(header.Authorization, "Bearer "+apiKey)
	req.Header.Set(header.ContentType, header.ApplicationJSON)
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "openai request failed")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read OpenAI response")
	}

	// Parse response into struct
	var openAIResp OpenAIResponse
	err = json.Unmarshal(data, &openAIResp)
	if err != nil {
		logger.ContextKV(ctx, xlog.ERROR,
			"err", err.Error(),
			"data", string(data))
		return errors.Wrap(err, "failed to decode OpenAI response")
	}

	// Check for API errors
	if openAIResp.Error != nil {
		return errors.Errorf("OpenAI API error: %s (type: %s, code: %d)", openAIResp.Error.Message, openAIResp.Error.Type, openAIResp.Error.Code)
	}

	var content string
	// Extract translated JSON from the structured response
	if len(openAIResp.Choices) > 0 {
		content = openAIResp.Choices[0].Message.Content
	}

	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if content == "" {
		return errors.New("unexpected response format")
	}

	var resValues values.MapAny
	err = json.Unmarshal([]byte(content), &resValues)
	if err != nil {
		return errors.Wrap(err, "failed to decode translated JSON")
	}

	for k := range texts {
		id := tag(k)
		if v := resValues.String(strconv.FormatUint(id, 10)); v != "" {
			texts[k] = v
		}
	}

	return nil
}
