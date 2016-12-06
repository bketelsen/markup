package markup

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"

	"github.com/murlokswarm/uid"
)

const (
	// HTML represents a standard HTML element.
	HTML ElementType = iota

	// Component represents a component element.
	Component

	// Text represents a text element.
	Text
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

// ElementType represents the type of an element.
type ElementType uint8

// Element represents a markup element.
type Element struct {
	Name       string
	Type       ElementType
	ID         uid.ID
	ContextID  uid.ID
	Attributes AttrList
	Parent     *Element
	Children   []*Element
	Component  Componer
}

func newElement(token xml.StartElement, parent *Element) *Element {
	n := token.Name.Local
	t := HTML

	if IsComponentName(n) {
		t = Component
	}

	return &Element{
		Name:       n,
		Type:       t,
		Attributes: makeAttrList(token.Attr),
		Parent:     parent,
	}
}

func newTextElement(text string, parent *Element) *Element {
	return &Element{
		Name: "text",
		Type: Text,
		Attributes: AttrList{
			Attr{Name: "value", Value: text},
		},
		Parent: parent,
	}
}

// HTML returns the HTML representation of the element.
func (e *Element) HTML() string {
	return e.html(0)
}

func (e *Element) html(level int) (markup string) {
	indt := indent(level)

	if e.Type == Text {
		text := e.Attributes[0].Value
		text = html.EscapeString(text)
		markup += fmt.Sprintf("%v%v", indt, text)
		return
	}

	if e.Type == Component {
		if e.Component == nil {
			markup += fmt.Sprintf("%v<!-- %v -->", indt, e.Name)
			return
		}

		compoRoot := components[e.Component].Root
		markup += compoRoot.html(level)
		return
	}

	markup = fmt.Sprintf("%v<%v", indt, e.Name)

	for _, attr := range e.Attributes {
		if attr.isEvent() {
			markup += fmt.Sprintf(" %v=\"CallEvent('%v', '%v', this, event)\"",
				strings.TrimLeft(attr.Name, "_"),
				e.ID,
				attr.Value)
			continue
		}

		markup += fmt.Sprintf(" %v=\"%v\"", attr.Name, attr.Value)
	}

	if len(e.ID) != 0 {
		markup += fmt.Sprintf(" data-murlok-id=\"%v\"", e.ID)
	}

	if len(e.Children) == 0 {
		if _, isVoidElement := voidElementNames[e.Name]; isVoidElement || IsComponentName(e.Name) {
			markup += " />"
			return
		}

		markup += fmt.Sprintf("></%v>", e.Name)
		return
	}

	markup += ">"

	for _, c := range e.Children {
		markup += "\n" + c.html(level+1)
	}

	markup += fmt.Sprintf("\n%v</%v>", indt, e.Name)
	return
}

func (e *Element) String() string {
	return fmt.Sprintf("[\033[36m%v\033[00m \033[33m%v\033[00m]", e.Name, e.ID)
}

func indent(level int) (ret string) {
	for i := 0; i < level; i++ {
		ret += "  "
	}
	return
}
