package ml

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

func (d *decoder) Decode() (*Element, error) {
	if err := d.next(); err != nil {
		return nil, err
	}

	if d.root == nil {
		return nil, errors.New("empty markup")
	}

	return d.root, nil
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

func decode(r io.Reader) (*Element, error) {
	d := newDecoder(r)
	return d.Decode()
}

func decodeString(s string) (*Element, error) {
	r := bytes.NewBufferString(s)
	return decode(r)
}
