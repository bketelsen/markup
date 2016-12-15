package markup

import (
	"reflect"
	"testing"

	"html"

	"github.com/murlokswarm/uid"
)

type Foo struct {
}

func (f *Foo) Render() string {
	return fooXML
}

type Bar struct {
	Value int
}

type PropsTest struct {
	String       string
	Bool         bool
	Int          int
	Uint         uint
	Float        float64
	Bar          Bar
	NotSupported *Foo
}

func (p *PropsTest) Render() string {
	return `
<div>
	<p>String: {{.String}}</p>
	<p>Bool: {{.Bool}}</p>
	<p>Int: {{.Int}}</p>
	<p>Uint: {{.Uint}}</p>
	<p>Float: {{.Float}}</p>
</div>
	`
}

func TestRegisterComponent(t *testing.T) {
	RegisterComponent("Foo", func() Componer {
		return &Foo{}
	})

	RegisterComponent("Foo", func() Componer {
		return &Foo{}
	})
}

func TestRegisterComponentPanic(t *testing.T) {
	defer func() { t.Log(recover()) }()

	RegisterComponent("foo", func() Componer {
		return &Foo{}
	})

	t.Error("should panic")
}

func TestComponentToHTML(t *testing.T) {
	c := &Hello{}
	ctx := uid.Context()

	if _, err := ComponentToHTML(c); err == nil {
		t.Error("should error")
	}

	if _, err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	HTML, err := ComponentToHTML(c)
	if err != nil {
		t.Error(err)
	}

	t.Log(HTML)
}

func TestIsComponentName(t *testing.T) {
	if name := "Foo"; !IsComponentName(name) {
		t.Errorf("%v should be a component name", name)
	}

	if name := "foo"; IsComponentName(name) {
		t.Errorf("%v should not be a component name", name)
	}

	if IsComponentName("") {
		t.Error("empty string should not be a component name")
	}
}

func TestCreateComponent(t *testing.T) {
	if _, err := createComponent("Foo"); err != nil {
		t.Error(err)
	}
}

func TestCreateComponentError(t *testing.T) {
	if _, err := createComponent("HyperFoo"); err == nil {
		t.Error("should error")
	}
}

func TestUpdateComponentFields(t *testing.T) {
	attrs := AttrList{
		Attr{Name: "String", Value: "Hello"},
		Attr{Name: "Int", Value: "42"},
		// Attr{Name: "Struct", Value: convertToJSON(Bar{Value: 21})},
	}

	compo := &PropsTest{}

	if err := updateComponentFields(compo, attrs); err != nil {
		t.Error(err)
	}

	if compo.String != "Hello" {
		t.Errorf("compo.String should be \"Hello\": \"%v\"", compo.String)
	}

	if compo.Int != 42 {
		t.Errorf("compo.String should be 42: %v", compo.Int)
	}

	// if compo.Bar.Value != 21 {
	// 	t.Errorf("compo.Bar.Value should be 21: %v", compo.Bar.Value)
	// }
}

func TestUpdateComponentFieldsError(t *testing.T) {
	attrs := AttrList{
		Attr{Name: "String", Value: "Hello"},
		Attr{Name: "Int", Value: "42.42"},
	}

	compo := &PropsTest{}

	if err := updateComponentFields(compo, attrs); err == nil {
		t.Error("should error")
	}
}

func TestUpdateComponentField(t *testing.T) {
	compo := &PropsTest{}
	compoValue := reflect.ValueOf(compo)
	compoValue = reflect.Indirect(compoValue)

	// String
	if err := updateComponentField(compoValue, Attr{Name: "String", Value: "Hello, World"}); err != nil {
		t.Error(err)
	}

	if hello := "Hello, World"; compo.String != hello {
		t.Errorf("compo.String should be \"%v\": \"%v\"", hello, compo.String)
	}

	// Bool
	if err := updateComponentField(compoValue, Attr{Name: "Bool", Value: "true"}); err != nil {
		t.Error(err)
	}

	if b := true; compo.Bool != b {
		t.Errorf("compo.Bool should be \"%v\": \"%v\"", b, compo.Bool)
	}

	// Int
	if err := updateComponentField(compoValue, Attr{Name: "Int", Value: "-42"}); err != nil {
		t.Error(err)
	}

	if n := -42; compo.Int != n {
		t.Errorf("compo.Int should be \"%v\": \"%v\"", n, compo.Int)
	}

	// Uint
	if err := updateComponentField(compoValue, Attr{Name: "Uint", Value: "42"}); err != nil {
		t.Error(err)
	}

	if n := uint(42); compo.Uint != n {
		t.Errorf("compo.Uint should be \"%v\": \"%v\"", n, compo.Uint)
	}

	// Float
	if err := updateComponentField(compoValue, Attr{Name: "Float", Value: "42.42"}); err != nil {
		t.Error(err)
	}

	if n := 42.42; compo.Float != n {
		t.Errorf("compo.Float should be \"%v\": \"%v\"", n, compo.Float)
	}

	// Struct
	bar := Bar{
		Value: 21,
	}
	j := convertToJSON(bar)
	j = html.UnescapeString(j)

	if err := updateComponentField(compoValue, Attr{Name: "Bar", Value: j}); err != nil {
		t.Error(err)
	}

	if n := 21; compo.Bar.Value != n {
		t.Errorf("compo.Bar.Value should be \"%v\": \"%v\"", n, compo.Bar.Value)
	}

	// NotSupported
	if err := updateComponentField(compoValue, Attr{Name: "NotSupported", Value: "Fooooo"}); err != nil {
		t.Error(err)
	}

}

func TestUpdateComponentFieldError(t *testing.T) {
	compo := &PropsTest{}
	compoValue := reflect.ValueOf(compo)
	compoValue = reflect.Indirect(compoValue)

	// No field
	if err := updateComponentField(compoValue, Attr{Name: "Foo"}); err == nil {
		t.Error("should error")
	}

	// Bool
	if err := updateComponentField(compoValue, Attr{Name: "Bool", Value: "foo"}); err == nil {
		t.Error("should error")
	}

	// Int
	if err := updateComponentField(compoValue, Attr{Name: "Int", Value: "-42.42"}); err == nil {
		t.Error("should error")
	}

	// Uint
	if err := updateComponentField(compoValue, Attr{Name: "Uint", Value: "-42"}); err == nil {
		t.Error("should error")
	}

	// Float
	if err := updateComponentField(compoValue, Attr{Name: "Float", Value: "42*$"}); err == nil {
		t.Error("should error")
	}

	// Struct
	bar := Bar{
		Value: 21,
	}
	j := convertToJSON(bar)

	if err := updateComponentField(compoValue, Attr{Name: "Bar", Value: j}); err == nil {
		t.Error("should error")
	}
}
