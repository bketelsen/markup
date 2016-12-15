package markup

import (
	"bytes"
	"encoding/json"
	"text/template"
)

func render(m string, data interface{}) (rendered string, err error) {
	var b bytes.Buffer

	fnmap := template.FuncMap{
		"json": convertToJSON,
	}

	tmpl, err := template.New("").Funcs(fnmap).Parse(m)
	if err != nil {
		return
	}

	if err = tmpl.Execute(&b, data); err != nil {
		return
	}

	rendered = b.String()
	return
}

func convertToJSON(v interface{}) (out string) {
	b, _ := json.Marshal(v)
	out = template.HTMLEscapeString(string(b))
	return
}
