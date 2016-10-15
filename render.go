package ml

import (
	"bytes"
	"text/template"
)

func renderMarkup(markup string, data interface{}) (rendered string, err error) {
	var tmpl *template.Template
	var b bytes.Buffer

	if tmpl, err = template.New("").Parse(markup); err != nil {
		return
	}

	if err = tmpl.Execute(&b, data); err != nil {
		return
	}

	rendered = b.String()
	return
}
