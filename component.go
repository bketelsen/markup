package markup

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/murlokswarm/log"
)

var (
	compoBuilders = map[string]ComponentBuilder{}
	compoElements = map[Componer]*Element{}
)

// Componer is the interface that describes a component.
type Componer interface {
	// Render should returns a markup.
	// The markup can be a template string following the text/template standard
	// package rules.
	Render() string
}

// ComponentBuilder defines a function that build a component.
type ComponentBuilder func() Componer

// RegisterComponent registers a component builder. It allow to know which
// component should be built when name is found into a markup.
// Should be called in a init() function, one time per component.
func RegisterComponent(name string, b ComponentBuilder) {
	if !isComponentName(name) {
		log.Panicf("\"%v\" is an invalid component name. must not be empty and should have its first letter capitalized", name)
	}

	if _, ok := compoBuilders[name]; ok {
		log.Infof("component builder for \033[35m%v\033[00m is overloaded with %T", name, b)
	}

	compoBuilders[name] = b
}

// ComponentToHTML returns the HTML representation of the component.
// returns an error if c is not mounted.
func ComponentToHTML(c Componer) (HTML string, err error) {
	var rootElem *Element
	var mounted bool

	if rootElem, mounted = compoElements[c]; !mounted {
		err = fmt.Errorf("%#v is not mounted", c)
		return
	}

	HTML = rootElem.HTML()
	return
}

func isComponentName(v string) bool {
	if len(v) == 0 {
		return false
	}

	return v[0] >= 'A' && v[0] <= 'Z'
}

func createComponent(name string) (c Componer, err error) {
	var registered bool
	var builder ComponentBuilder

	if builder, registered = compoBuilders[name]; !registered {
		err = fmt.Errorf("component %v is not registered", name)
		return
	}

	c = builder()
	return
}

func updateComponentFields(c Componer, attrs AttrList) (err error) {
	compo := reflect.Indirect(reflect.ValueOf(c))

	for _, attr := range attrs {
		if err = updateComponentField(compo, attr); err != nil {
			return
		}
	}

	return
}

func updateComponentField(compo reflect.Value, attr Attr) error {
	field := compo.FieldByName(attr.Name)
	if !field.IsValid() {
		return fmt.Errorf("no field %v in %T", attr.Name, compo.Interface())
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(attr.Value)

	case reflect.Bool:
		if attr.Value != "true" && attr.Value != "false" {
			return fmt.Errorf("boolean attributes in a component must be set to true or false: %v", attr.Value)
		}

		b, _ := strconv.ParseBool(attr.Value)
		field.SetBool(b)

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		n, err := strconv.ParseInt(attr.Value, 0, 64)
		if err != nil {
			return err
		}

		field.SetInt(n)

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		n, err := strconv.ParseUint(attr.Value, 0, 64)
		if err != nil {
			return err
		}

		field.SetUint(n)

	case reflect.Float64, reflect.Float32:
		n, err := strconv.ParseFloat(attr.Value, 64)
		if err != nil {
			return err
		}

		field.SetFloat(n)

	default:
		log.Warnf("in \033[35m%T\033[00m: field \033[36m%v\033[00m of type \033[34m%T\033[00m can't be mapped",
			compo.Interface(),
			attr.Name,
			field.Interface())
	}

	return nil
}
