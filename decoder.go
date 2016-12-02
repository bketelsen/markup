package markup

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

type decoder struct {
	xmlDecoder *xml.Decoder
	root       *Element
	current    *Element
}

func newDecoder(r io.Reader) *decoder {
	return &decoder{
		xmlDecoder: xml.NewDecoder(r),
	}
}

func (d *decoder) Decode() (rootElem *Element, err error) {
	if err = d.next(); err != nil {
		return
	}

	if d.root == nil {
		err = errors.New("empty markup")
		return
	}

	rootElem = d.root
	return
}

func (d *decoder) next() error {
	token, err := d.xmlDecoder.Token()
	if err != nil {
		if err == io.EOF {
			return nil
		}

		return err
	}

	switch t := token.(type) {
	case xml.StartElement:
		elem := newElement(t, d.current)

		if d.root == nil {
			d.root = elem
		}

		if d.current != nil {
			d.current.Children = append(d.current.Children, elem)
		}

		d.current = elem

	case xml.EndElement:
		if d.current.Parent != nil {
			d.current = d.current.Parent
		}

	case xml.CharData:
		text := strings.TrimSpace(string(t))
		if len(text) == 0 {
			break
		}

		if d.current == nil {
			return fmt.Errorf("\"%v\" must be surrounded by HTML tags", t)
		}

		elem := newTextElement(text, d.current)
		d.current.Children = append(d.current.Children, elem)
	}
	return d.next()
}

// Decode translates an markup string to an element tree.
// Returns the root element of the tree.
func Decode(s string) (*Element, error) {
	r := bytes.NewBufferString(s)
	d := newDecoder(r)
	return d.Decode()
}
