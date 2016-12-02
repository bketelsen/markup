package markup

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestElementHTML(t *testing.T) {
	elem, err := Decode(fooXML)
	if err != nil {
		t.Fatal(err)
	}

	elem.ID = uid.Elem()
	t.Log(elem.HTML())
}

func TestElementString(t *testing.T) {
	elem, err := Decode(fooXML)
	if err != nil {
		t.Fatal(err)
	}

	elem.ID = uid.Elem()
	t.Log(elem.String())
}
