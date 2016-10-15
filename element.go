package ml

import (
	"encoding/xml"
	"fmt"

	"github.com/murlokswarm/uid"
)

const (
	htmlTag tagType = iota
	componentTag
	textTag
)

var (
	elements         = map[uid.ID]*Element{}
	voidElementNames = map[string]bool{
		"area":    true,
		"base":    true,
		"br":      true,
		"col":     true,
		"command": true,
		"embed":   true,
		"hr":      true,
		"img":     true,
		"input":   true,
		"link":    true,
		"meta":    true,
		"param":   true,
		"source":  true,
	}
)

type tagType uint8

// Element represents a ML element.
type Element struct {
	Name       string
	ID         uid.ID
	Context    uid.ID
	Attributes AttrList
	Parent     *Element
	Children   []*Element
	Component  Componer
	tagType    tagType
}

// Markup returns the markup representation of the element.
func (e *Element) Markup() string {
	return e.markup(0)
}

func (e *Element) markup(level int) (m string) {
	indt := indent(level)
	m = fmt.Sprintf("%v<\033[35m%v\033[00m", indt, e.Name)

	for _, attr := range e.Attributes {
		m += fmt.Sprintf(" \033[36m%v\033[00m=\"%v\"", attr.Name, attr.Value)
	}

	if len(e.ID) != 0 {
		m += fmt.Sprintf(" data-murlok-id=\"%v\"", e.ID)
	}

	if len(e.Children) == 0 {
		m += " />"
		return
	}

	m += ">"

	for _, c := range e.Children {
		m += "\n" + c.markup(level+1)
	}

	m += fmt.Sprintf("\n%v</\033[35m%v\033[00m>", indt, e.Name)
	return
}

// HTML returns the HTML representation of the element.
func (e *Element) HTML() string {
	return e.html(0)
}

func (e *Element) html(level int) (m string) {
	indt := indent(level)

	if e.tagType == textTag {
		text := e.Attributes[0].Value
		m += fmt.Sprintf("%v%v", indt, text)
		return
	}

	if e.tagType == componentTag {
		if e.Component == nil {
			m += fmt.Sprintf("%v<!-- %v -->", indt, e.Name)
			return
		}

		compoRoot := compoElements[e.Component]
		m += compoRoot.html(level)
		return
	}

	m = fmt.Sprintf("%v<%v", indt, e.Name)

	for _, attr := range e.Attributes {
		attrValue := attr.Value

		if attrEventVal, ok := attrEventValue(attr.Value); ok {
			attrValue = fmt.Sprintf("CallEvent('%v', '%v', 'event')", e.ID, attrEventVal)
		}

		m += fmt.Sprintf(" %v=\"%v\"", attr.Name, attrValue)
	}

	if len(e.ID) != 0 {
		m += fmt.Sprintf(" data-murlok-id=\"%v\"", e.ID)
	}

	if len(e.Children) == 0 {
		if _, isVoidElement := voidElementNames[e.Name]; isVoidElement || isComponentName(e.Name) {
			m += " />"
			return
		}

		m += fmt.Sprintf("></%v>", e.Name)
		return
	}

	m += ">"

	for _, c := range e.Children {
		m += "\n" + c.html(level+1)
	}

	m += fmt.Sprintf("\n%v</%v>", indt, e.Name)
	return
}

func newElement(token xml.StartElement, parent *Element) *Element {
	name := token.Name.Local
	tagType := htmlTag

	if isComponentName(name) {
		tagType = componentTag
	}

	return &Element{
		Name:       name,
		Attributes: makeAttrList(token.Attr),
		Parent:     parent,
		tagType:    tagType,
	}
}

func newTextElement(text string, parent *Element) *Element {
	return &Element{
		Name: "text",
		Attributes: AttrList{
			Attr{Name: "value", Value: text},
		},
		Parent:  parent,
		tagType: textTag,
	}
}

func indent(level int) (ret string) {
	for i := 0; i < level; i++ {
		ret += "  "
	}

	return
}
