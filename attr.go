package ml

import (
	"encoding/xml"
	"strings"
)

// Attr represents an attribute in an ML element (Name=Value).
type Attr struct {
	Name  string
	Value string
}

func makeAttr(a xml.Attr) Attr {
	return Attr{
		Name:  a.Name.Local,
		Value: a.Value,
	}
}

// AttrList represents a list of Attr.
type AttrList []Attr

func makeAttrList(atrributes []xml.Attr) AttrList {
	attrs := make(AttrList, len(atrributes))

	for i, a := range atrributes {
		attrs[i] = makeAttr(a)
	}

	return attrs
}

func (l AttrList) equals(other AttrList) bool {
	if len(l) != len(other) {
		return false
	}

	for i, attr := range l {
		if attr != other[i] {
			return false
		}
	}

	return true
}

func attrEventValue(v string) (string, bool) {
	if len(v) > 0 && v[0] == '@' {
		return strings.TrimLeft(v, "@"), true
	}

	return "", false
}
