package markup

import (
	"bytes"
	"text/template"
)

func render(m string, data interface{}) (rendered string, err error) {
	var b bytes.Buffer

	tmpl, err := template.New("").Parse(m)
	if err != nil {
		return
	}

	if err = tmpl.Execute(&b, data); err != nil {
		return
	}

	rendered = b.String()
	return
}
