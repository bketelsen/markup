package ml

import (
	"testing"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

type Hello struct {
	Number     float64
	MountError bool
}

func (h *Hello) Render() string {
	return `
<div>
    Hello,
    <World Greeting="World" Number="{{.Number}}" MountError="{{.MountError}}" />
</div>
    `
}

type World struct {
	Greeting   string
	Number     uint
	MountError bool
}

func (w *World) Render() string {
	return `
<input value="{{.Greeting}} {{.Number}}" />
{{if .MountError}}
	<BadMarkup />
{{end}}
    `
}

func (w *World) OnMount() {
	log.Infof("%T is mounted: %+v", w, w)
}

func (w *World) OnDismount() {
	log.Infof("%T is dismounted: %+v", w, w)
}

type BadTemplate struct {
	Placebo string
}

func (b *BadTemplate) Render() string {
	return `
<div wtf="{{.Wtf}}"></div>
`
}

type BadMarkup struct {
	Placebo string
}

func (b *BadMarkup) Render() string {
	return `
<div></span>
`
}

type NonexistentChild struct {
	Placebo string
}

func (n *NonexistentChild) Render() string {
	return `
<div>
    <Unknown />
</div>
`
}

type NoField struct {
}

func (n *NoField) Render() string {
	return `
<div></div>
`
}

type ComponentRoot struct {
	Placebo bool
}

func (c *ComponentRoot) Render() string {
	return `
<Hello />
`
}

func init() {
	RegisterComponent("Hello", func() Componer {
		return &Hello{}
	})

	RegisterComponent("World", func() Componer {
		return &World{}
	})

	RegisterComponent("NonexistentChild", func() Componer {
		return &NonexistentChild{}
	})
}

func TestMount(t *testing.T) {
	hello := &Hello{}
	ctx := uid.Context()

	elementsLen := len(elements)
	compoElemsLen := len(compoElements)

	if err := Mount(hello, ctx); err != nil {
		t.Error(err)
	}

	if l := len(elements); l != elementsLen+2 {
		t.Errorf("l should be %v: %v", elementsLen+2, l)
	}

	if l := len(compoElements); l != compoElemsLen+2 {
		t.Errorf("l should be %v: %v", compoElemsLen+2, l)
	}

	elem := compoElements[hello]
	t.Log(elem.HTML())
}

func TestMountError(t *testing.T) {
	ctx := uid.Context()

	// no field
	noField := &NoField{}

	if err := Mount(noField, ctx); err == nil {
		t.Error("should error")
	}

	// already mounted
	world := &World{}

	if err := Mount(world, ctx); err != nil {
		t.Fatal(err)
	}

	if err := Mount(world, ctx); err == nil {
		t.Error("should error")
	}

	// bad template
	badTpl := &BadTemplate{}

	if err := Mount(badTpl, ctx); err == nil {
		t.Error("should error")
	}

	// bad markup
	badMarkup := &BadMarkup{}

	if err := Mount(badMarkup, ctx); err == nil {
		t.Error("should error")
	}

	// nonexistent child
	nonexistentChild := &NonexistentChild{}

	if err := Mount(nonexistentChild, ctx); err == nil {
		t.Error("should error")
	}

	// bad attibute
	hello := &Hello{
		Number: 42.99,
	}

	if err := Mount(hello, ctx); err == nil {
		t.Error("should error")
	}

	// mount error
	hello = &Hello{
		MountError: true,
	}

	if err := Mount(hello, ctx); err == nil {
		t.Error("should error")
	}

	// component root
	compoRoot := &ComponentRoot{}
	if err := Mount(compoRoot, ctx); err == nil {
		t.Error("should error")
	}
}

func TestDismount(t *testing.T) {
	hello := &Hello{}
	ctx := uid.Context()

	elementsLen := len(elements)
	compoElemsLen := len(compoElements)

	if err := Mount(hello, ctx); err != nil {
		t.Fatal(err)
	}

	if err := Dismount(hello); err != nil {
		t.Error(err)
	}

	if l := len(elements); l != elementsLen {
		t.Errorf("l should be %v: %v", elementsLen, l)
	}

	if l := len(compoElements); l != compoElemsLen {
		t.Errorf("l should be %v: %v", compoElemsLen, l)
	}
}

func TestDismountError(t *testing.T) {
	// not mounted
	hello := &Hello{}
	Dismount(hello)
}
