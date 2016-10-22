package markup

import (
	"bytes"
	"text/template"
)

func render(m string, data interface{}) (rendered string, err error) {
	var tmpl *template.Template
	var b bytes.Buffer

	if tmpl, err = template.New("").Parse(m); err != nil {
		return
	}

	if err = tmpl.Execute(&b, data); err != nil {
		return
	}

	rendered = b.String()
	return
}
