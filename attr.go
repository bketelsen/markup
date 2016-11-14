package markup

import "encoding/xml"

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

func (a Attr) isEvent() bool {
	if len(a.Name) > 0 && a.Name[0] == '_' {
		return true
	}

	return false
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
