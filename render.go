package markup

import (
	"bytes"
	"encoding/json"
	"text/template"
	"time"
)

func render(c Componer) (rendered string, err error) {
	var b bytes.Buffer

	fnmap := template.FuncMap{
		"json": convertToJSON,
		"time": formatTime,
	}
	tmpl := template.Must(template.New("Render").Funcs(fnmap).Parse(c.Render()))

	if err = tmpl.Execute(&b, c); err != nil {
		return
	}

	rendered = b.String()
	return
}

func convertToJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return template.HTMLEscapeString(string(b))
}

func formatTime(t time.Time, layout string) string {
	return t.Format(layout)
}
