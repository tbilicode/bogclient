package print

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/effective-security/x/slices"
	"github.com/effective-security/x/values"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// JSON prints value to out
func JSON(w io.Writer, value any) error {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return errors.WithMessage(err, "failed to encode")
	}
	_, _ = w.Write(json)
	_, _ = w.Write([]byte{'\n'})
	return nil
}

// Yaml prints value  to out
func Yaml(w io.Writer, value any) error {
	y, err := yaml.Marshal(value)
	if err != nil {
		return errors.WithMessage(err, "failed to encode")
	}
	_, _ = w.Write(y)
	return nil
}

// Object prints value to out in format
func Object(w io.Writer, format string, value any) error {
	if format == "yaml" {
		return Yaml(w, value)
	}
	if format == "json" {
		return JSON(w, value)
	}
	Print(w, value)
	return nil
}

// Print value
func Print(w io.Writer, value any) {
	switch t := value.(type) {
	case map[string]string:
		Map(w, []string{"Key", "Value"}, t)
	case []string:
		Strings(w, t)

	default:
		_ = JSON(w, value)
	}
}

// Map prints map
func Map(w io.Writer, header []string, vals map[string]string) {
	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetHeader(header)
	table.SetHeaderLine(true)

	for _, k := range values.OrderedMapKeys(vals) {
		table.Append([]string{k, slices.StringUpto(vals[k], 80)})
	}

	table.Render()
	fmt.Fprintln(w)
}

// Strings prints strings
func Strings(w io.Writer, res []string) {
	for _, r := range res {
		fmt.Fprintln(w, r)
	}
}
