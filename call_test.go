package ml

import (
	"reflect"
	"testing"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

type ComponentWithFunc struct {
	Name string
}

func (c *ComponentWithFunc) OnCallTest() {
	log.Info("OnCallTest")
}

func (c *ComponentWithFunc) OnCallTestWithArg(arg string, number int) {
	log.Info("OnCallTestWithArg")
}

func (c *ComponentWithFunc) OnCallTestWithMultipleArgs(arg int) {
	log.Info("OnCallTestWithMultipleArgs")
}

func (c *ComponentWithFunc) Render() string {
	return `
<h1>{{.Name}}</h1>
    `
}

func init() {
	RegisterComponent("ComponentWithFunc", func() Componer {
		return &ComponentWithFunc{}
	})
}

type FuncArg struct {
	Number int
	String string
}

func TestCall(t *testing.T) {
	c := &ComponentWithFunc{}
	ctx := uid.Context()

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	rootElem := compoElements[c]

	// no arg
	if err := Call(rootElem.ID, "OnCallTest", ""); err != nil {
		t.Error(err)
	}

	// arg
	if err := Call(rootElem.ID, "OnCallTestWithArg", "42"); err != nil {
		t.Error(err)
	}
}

func TestCallError(t *testing.T) {
	c := &ComponentWithFunc{}
	ctx := uid.Context()

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	rootElem := compoElements[c]

	// undefined method
	Call(rootElem.ID, "OnCallFoo", "")

	// multiple args
	Call(rootElem.ID, "OnCallTestWithMultipleArgs", "")

	// unmounted elem
	if err := Call(uid.Elem(), "OnCallTest", ""); err == nil {
		t.Error("should error")
	}
}

func TestCreateCallArg(t *testing.T) {
	var arg FuncArg
	var number float64
	var str string
	var ret reflect.Value
	var converted bool
	var err error

	// struct
	if ret, err = createCallArg(reflect.TypeOf(arg), `{"Number": 42, "String": "Hello"}`); err != nil {
		t.Fatal(err)
	}

	expectedArg := FuncArg{Number: 42, String: "Hello"}

	if arg, converted = ret.Interface().(FuncArg); !converted {
		t.Fatal("type conversion failed")
	}

	if arg != expectedArg {
		t.Errorf("arg should be %v: %v", expectedArg, arg)
	}

	// number
	if ret, err = createCallArg(reflect.TypeOf(number), "42"); err != nil {
		t.Fatal(err)
	}

	if number, converted = ret.Interface().(float64); !converted {
		t.Fatal("type conversion failed")
	}

	if expectedNumber := 42.0; number != expectedNumber {
		t.Errorf("arg should be %v: %v", expectedNumber, number)
	}

	// string
	if ret, err = createCallArg(reflect.TypeOf(str), `"Hello"`); err != nil {
		t.Fatal(err)
	}

	if str, converted = ret.Interface().(string); !converted {
		t.Fatal("type conversion failed")
	}

	if expectedStr := "Hello"; str != expectedStr {
		t.Errorf("arg should be %v: %v", expectedStr, str)
	}
}

func TestCreateCallArgError(t *testing.T) {
	if _, err := createCallArg(reflect.TypeOf("str"), "Hello"); err == nil {
		t.Error("should error")
	}
}
