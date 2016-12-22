package markup

import (
	"bytes"
	"encoding/json"
	"text/template"
)

func render(c Componer) (rendered string, err error) {
	var b bytes.Buffer

	fnmap := template.FuncMap{
		"json": convertToJSON,
	}
	tmpl := template.Must(template.New("Render").Funcs(fnmap).Parse(c.Render()))

	if err = tmpl.Execute(&b, c); err != nil {
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
