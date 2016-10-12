package ml

import (
	"fmt"
	"reflect"

	"github.com/murlokswarm/uid"
)

// Mounter is the interface that wraps OnMount method.
// OnMount si called when a component is mounted.
type Mounter interface {
	OnMount() error
}

// Dismounter is the interface that wraps OnDismount method.
// OnDismount si called when a component is dismounted.
type Dismounter interface {
	OnDismount() error
}

// Mount mounts a component.
// It enable bidirectional communication between le component and the
// underlying driver.
func Mount(c Componer, ctx uid.ID) error {
	compoVal := reflect.Indirect(reflect.ValueOf(c))
	if compoVal.NumField() == 0 {
		return fmt.Errorf("\033[33m%T\033[00m must have at least 1 field", c)
	}

	rendered, err := parseTemplate(c.Render(), c)
	if err != nil {
		return err
	}

	elem, err := decodeString(rendered)
	if err != nil {
		return err
	}

	compoElements[c] = elem

	if err := mount(elem, c, ctx); err != nil {
		return err
	}

	if mounter, ok := c.(Mounter); ok {
		if err := mounter.OnMount(); err != nil {
			return err
		}
	}

	return nil
}

func mount(e *Element, c Componer, ctx uid.ID) error {
	switch e.tagType {
	case htmlTag:
		return mountElement(e, c, ctx)

	case componentTag:
		return mountComponent(e, ctx)
	}

	return nil
}

func mountElement(e *Element, c Componer, ctx uid.ID) error {
	e.ID = uid.Elem()
	e.Context = ctx
	e.Component = c
	elements[e.ID] = e

	for _, child := range e.Children {
		if err := mount(child, c, ctx); err != nil {
			return err
		}
	}

	return nil
}

func mountComponent(e *Element, ctx uid.ID) error {
	compo, err := createComponent(e.Name)
	if err != nil {
		return err
	}

	if err := updateComponentFields(compo, e.Attributes); err != nil {
		return err
	}

	if err := Mount(compo, ctx); err != nil {
		return err
	}

	e.Component = compo

	return nil
}

func Dismount(c Componer) error {
	elem, ok := compoElements[c]
	if !ok {
		return fmt.Errorf("%#v is already dismounted", c)
	}

	if err := dismount(elem); err != nil {
		return err
	}

	delete(compoElements, c)

	if dismounter, ok := c.(Dismounter); ok {
		if err := dismounter.OnDismount(); err != nil {
			return err
		}
	}

	return nil
}

func dismount(e *Element) error {
	switch e.tagType {
	case htmlTag:
		return dismountElement(e)

	case componentTag:
		return Dismount(e.Component)
	}

	return nil
}

func dismountElement(e *Element) error {
	for _, child := range e.Children {
		if err := dismount(child); err != nil {
			return err
		}
	}

	delete(elements, e.ID)
	return nil
}
