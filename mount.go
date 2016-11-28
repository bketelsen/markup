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
	var componentValue reflect.Value
	var rendered string
	var mounted bool
	var isMounter bool
	var mounter Mounter

	if v := reflect.ValueOf(c); v.Kind() != reflect.Ptr {
		err = fmt.Errorf("\033[33m%T\033[00m should be a pointer", c)
		return
	}

	if componentValue = reflect.Indirect(reflect.ValueOf(c)); componentValue.NumField() == 0 {
		err = fmt.Errorf("\033[33m%T\033[00m must have at least 1 field", c)
		return
	}

	if _, mounted = compoElements[c]; mounted {
		err = fmt.Errorf("component already mounted: %T %+v", c, c)
		return
	}

	if rendered, err = render(c.Render(), c); err != nil {
		return
	}

	if root, err = Decode(rendered); err != nil {
		return
	}

	if root.tagType != htmlTag {
		err = fmt.Errorf("component root must be a standard HTML tag: %T %+v", c, c)
		return
	}

	compoElements[c] = root

	if err = mount(root, c, ctx); err != nil {
		return
	}

	if mounter, isMounter = c.(Mounter); isMounter {
		mounter.OnMount()
	}

	return
}

func mount(e *Element, c Componer, ctx uid.ID) (err error) {
	switch e.tagType {
	case htmlTag:
		return mountElement(e, c, ctx)

	case componentTag:
		return mountComponent(e, ctx)
	}

	return
}

func mountElement(e *Element, c Componer, ctx uid.ID) (err error) {
	e.ID = uid.Elem()
	e.ContextID = ctx
	e.Component = c
	elements[e.ID] = e

	for _, child := range e.Children {
		if err = mount(child, c, ctx); err != nil {
			return
		}
	}

	return
}

func mountComponent(e *Element, ctx uid.ID) (err error) {
	var c Componer
	var root *Element

	if c, err = createComponent(e.Name); err != nil {
		return
	}

	if err = updateComponentFields(c, e.Attributes); err != nil {
		return
	}

	if root, err = Mount(c, ctx); err != nil {
		return
	}

	root.Parent = e.Parent
	e.ContextID = ctx
	e.Component = c

	return
}

// Dismount dismounts a component.
func Dismount(c Componer) {
	var rootElem *Element
	var dismounter Dismounter
	var mounted bool
	var isDismounter bool

	if rootElem, mounted = compoElements[c]; !mounted {
		log.Warnf("%#v is already dismounted", c)
		return
	}

	dismount(rootElem)
	delete(compoElements, c)

	if dismounter, isDismounter = c.(Dismounter); isDismounter {
		dismounter.OnDismount()
	}

	return
}

func dismount(e *Element) {
	switch e.tagType {
	case htmlTag:
		dismountElement(e)

	case componentTag:
		Dismount(e.Component)
	}
}

func dismountElement(e *Element) {
	for _, child := range e.Children {
		dismount(child)
	}

	delete(elements, e.ID)
}
