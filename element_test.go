package ml

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestElementMarkup(t *testing.T) {
	elem, err := decodeString(fooXML)
	if err != nil {
		t.Fatal(err)
	}

	elem.ID = uid.Elem()
	t.Log(elem.Markup())
}

func TestElementHTML(t *testing.T) {
	elem, err := decodeString(fooXML)
	if err != nil {
		t.Fatal(err)
	}

	elem.ID = uid.Elem()
	t.Log(elem.HTML())
}
