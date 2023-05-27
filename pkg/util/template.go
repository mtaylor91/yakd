package util

import (
	"bytes"
	"text/template"
)

// TemplateString renders a template string
func TemplateString(tmpl string, data interface{}) (string, error) {
	var b bytes.Buffer
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
