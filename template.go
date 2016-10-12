package ml

import (
	"bytes"
	"text/template"
)

func parseTemplate(tpl string, data interface{}) (string, error) {
	var b bytes.Buffer

	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}

	if err = tmpl.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
