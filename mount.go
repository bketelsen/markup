package markup

import (
	"fmt"
	"reflect"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

// Mounter is the interface that wraps OnMount method.
// OnMount si called when a component is mounted.
type Mounter interface {
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
// OnDismount si called when a component is dismounted.
type Dismounter interface {
	OnDismount()
}

// Mount maps a component and its underlying elements.
// It enable bidirectional communication between a component and the
// underlying driver.
func Mount(c Componer, ctx uid.ID) (root *Element, err error) {
	if v := reflect.ValueOf(c); v.Kind() != reflect.Ptr {
		err = fmt.Errorf("Mount accepts only pointers: \033[31m%T\033[00m", c)
		return
	}

	if compo, mounted := components[c]; mounted {
		compo.Count++
		return compo.Root, nil
	}

	rendered, err := render(c.Render(), c)
	if err != nil {
		return
	}

	if root, err = Decode(rendered); err != nil {
		return
	}

	if root.Type != HTML {
		err = fmt.Errorf("component root must be a standard HTML tag: %T %+v", c, c)
		return
	}

	if err = mount(root, c, ctx); err != nil {
		return
	}

	components[c] = &component{
		Count: 1,
		Root:  root,
	}

	if mounter, isMounter := c.(Mounter); isMounter {
		mounter.OnMount()
	}
	return
}

func mount(e *Element, c Componer, ctx uid.ID) (err error) {
	switch e.Type {
	case HTML:
		return mountElement(e, c, ctx)

	case Component:
		return mountComponent(e, ctx)
	}
	return
}

func mountElement(e *Element, c Componer, ctx uid.ID) error {
	e.ID = uid.Elem()
	e.ContextID = ctx
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
	c, err := createComponent(e.Name)
	if err != nil {
		return err
	}

	if err = updateComponentFields(c, e.Attributes); err != nil {
		return err
	}

	root, err := Mount(c, ctx)
	if err != nil {
		return err
	}

	root.Parent = e.Parent
	e.ContextID = ctx
	e.Component = c
	return err
}

// Dismount dismounts a component.
func Dismount(c Componer) {
	compo, mounted := components[c]
	if !mounted {
		log.Warnf("%#v is already dismounted", c)
		return
	}

	if compo.Count--; compo.Count == 0 {
		dismount(compo.Root)
		delete(components, c)
	}

	if dismounter, isDismounter := c.(Dismounter); isDismounter {
		dismounter.OnDismount()
	}
	return
}

func dismount(e *Element) {
	switch e.Type {
	case HTML:
		dismountElement(e)

	case Component:
		Dismount(e.Component)
	}
}

func dismountElement(e *Element) {
	for _, child := range e.Children {
		dismount(child)
	}
	delete(elements, e.ID)
}
