package markup

import (
	"fmt"
	"reflect"

	"github.com/murlokswarm/log"
)

var (
	compoBuilders = map[string]func() Componer{}
	components    = map[Componer]*component{}
)

// Componer is the interface that describes a component.
type Componer interface {
	// Render should returns a markup.
	// The markup can be a template string following the text/template standard
	// package rules.
	Render() string
}

type component struct {
	Count int
	Root  *Node
}

// Register registers a component. Allows the component to be dynamically
// created when a tag with its struct name is found into a markup.
func Register(c Componer) {
	v := reflect.ValueOf(c)

	if k := v.Kind(); k != reflect.Ptr {
		log.Panicf("register accepts only components of kind %v: %v", reflect.Ptr, k)
	}

	t := v.Type().Elem()
	tag := t.Name()

	if !isComponentTag(tag) {
		log.Panicf("non exported components cannot be registered: %v", t)
	}

	compoBuilders[tag] = func() Componer {
		v := reflect.New(t)
		return v.Interface().(Componer)
	}
	log.Info("%v has been registered under the tag %v", t, tag)
}

// Root returns the root node of c.
func Root(c Componer) *Node {
	compo, mounted := components[c]
	if !mounted {
		log.Panicf("%#v is not mounted", c)
	}
	return compo.Root
}

// Markup returns the markup of c.
func Markup(c Componer) string {
	return Root(c).Markup()
}

func componentFromTag(tag string) (c Componer, err error) {
	b, registered := compoBuilders[tag]
	if !registered {
		err = fmt.Errorf("component tagged %v is not registered", tag)
		return
	}

	c = b()
	return
}

// func updateComponentFields(c Componer, attrs AttrList) error {
// 	compo := reflect.Indirect(reflect.ValueOf(c))

// 	for _, attr := range attrs {
// 		if err := updateComponentField(compo, attr); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func updateComponentField(compo reflect.Value, attr Attr) error {
// 	field := compo.FieldByName(attr.Name)
// 	if !field.IsValid() {
// 		return fmt.Errorf("no field %v in %T", attr.Name, compo.Interface())
// 	}

// 	switch field.Kind() {
// 	case reflect.String:
// 		field.SetString(attr.Value)

// 	case reflect.Bool:
// 		if attr.Value != "true" && attr.Value != "false" {
// 			return fmt.Errorf("boolean attributes in a component must be set to true or false: %v", attr.Value)
// 		}

// 		b, _ := strconv.ParseBool(attr.Value)
// 		field.SetBool(b)

// 	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
// 		n, err := strconv.ParseInt(attr.Value, 0, 64)
// 		if err != nil {
// 			return err
// 		}

// 		field.SetInt(n)

// 	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
// 		n, err := strconv.ParseUint(attr.Value, 0, 64)
// 		if err != nil {
// 			return err
// 		}

// 		field.SetUint(n)

// 	case reflect.Float64, reflect.Float32:
// 		n, err := strconv.ParseFloat(attr.Value, 64)
// 		if err != nil {
// 			return err
// 		}

// 		field.SetFloat(n)

// 	case reflect.Struct:
// 		s := reflect.New(field.Type())

// 		if err := json.Unmarshal([]byte(attr.Value), s.Interface()); err != nil {
// 			return err
// 		}

// 		field.Set(s.Elem())

// 	default:
// 		log.Warnf("in \033[35m%T\033[00m: field \033[36m%v\033[00m of type \033[34m%T\033[00m can't be mapped",
// 			compo.Interface(),
// 			attr.Name,
// 			field.Interface())
// 	}
// 	return nil
// }
