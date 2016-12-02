package markup

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
}

func (b *BadTemplate) Render() string {
	return `
<div wtf="{{.Wtf}}"></div>
`
}

type BadMarkup struct {
}

func (b *BadMarkup) Render() string {
	return `<div></span>`
}

type NonexistentChild struct {
}

func (n *NonexistentChild) Render() string {
	return `
<div>
    <Unknown />
</div>
`
}

type EmptyComponent struct{}

func (c *EmptyComponent) Render() string {
	return `<div></div>`
}

type CompoRoot struct {
}

func (c *CompoRoot) Render() string {
	return `<Hello />`
}

type NoPointer struct {
}

func (n NoPointer) Render() string {
	return `<div></div>`
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
	compoElemsLen := len(components)

	root, err := Mount(hello, ctx)
	if err != nil {
		t.Error(err)
	}

	defer Dismount(hello)

	if l := len(elements); l != elementsLen+2 {
		t.Errorf("l should be %v: %v", elementsLen+2, l)
	}

	if l := len(components); l != compoElemsLen+2 {
		t.Errorf("l should be %v: %v", compoElemsLen+2, l)
	}

	t.Log(root.HTML())
}

func TestMountEmpty(t *testing.T) {
	empty1 := &EmptyComponent{}
	empty2 := &EmptyComponent{}
	ctx := uid.Context()

	elemNum := len(elements)
	compoNum := len(components)

	if _, err := Mount(empty1, ctx); err != nil {
		t.Error(err)
	}

	if l := len(elements); l != elemNum+1 {
		t.Errorf("l should be %v: %v", elemNum+1, l)
	}

	if l := len(components); l != compoNum+1 {
		t.Errorf("l should be %v: %v", compoNum+1, l)
	}

	_, err := Mount(empty2, ctx)
	if err != nil {
		t.Error(err)
	}

	if l := len(elements); l != elemNum+1 {
		t.Errorf("l should be %v: %v", elemNum+1, l)
	}

	if l := len(components); l != compoNum+1 {
		t.Errorf("l should be %v: %v", compoNum+1, l)
	}

	Dismount(empty1)

	if l := len(elements); l != elemNum+1 {
		t.Errorf("l should be %v: %v", elemNum+1, l)
	}

	if l := len(components); l != compoNum+1 {
		t.Errorf("l should be %v: %v", compoNum+1, l)
	}

	Dismount(empty2)

	if l := len(elements); l != elemNum {
		t.Errorf("l should be %v: %v", elemNum, l)
	}

	if l := len(components); l != compoNum {
		t.Errorf("l should be %v: %v", compoNum, l)
	}
}

func TestMountError(t *testing.T) {
	ctx := uid.Context()

	// bad template
	badTpl := &BadTemplate{}

	if _, err := Mount(badTpl, ctx); err == nil {
		t.Error("should error")
	}

	// bad markup
	badMarkup := &BadMarkup{}

	if _, err := Mount(badMarkup, ctx); err == nil {
		t.Error("should error")
	}

	// nonexistent child
	nonexistentChild := &NonexistentChild{}

	if _, err := Mount(nonexistentChild, ctx); err == nil {
		t.Error("should error")
	}

	// bad attibute
	hello := &Hello{
		Number: 42.99,
	}

	if _, err := Mount(hello, ctx); err == nil {
		t.Error("should error")
	}

	// mount error
	hello = &Hello{
		MountError: true,
	}

	if _, err := Mount(hello, ctx); err == nil {
		t.Error("should error")
	}

	// component root
	compoRoot := &CompoRoot{}
	if _, err := Mount(compoRoot, ctx); err == nil {
		t.Error("should error")
	}

	// no pointer
	noPointer := NoPointer{}
	if _, err := Mount(noPointer, ctx); err == nil {
		t.Error("should error")
	}
}

func TestDismount(t *testing.T) {
	hello := &Hello{}
	ctx := uid.Context()

	elementsLen := len(elements)
	compoElemsLen := len(components)

	if _, err := Mount(hello, ctx); err != nil {
		t.Fatal(err)
	}

	Dismount(hello)

	if l := len(elements); l != elementsLen {
		t.Errorf("l should be %v: %v", elementsLen, l)
	}

	if l := len(components); l != compoElemsLen {
		t.Errorf("l should be %v: %v", compoElemsLen, l)
	}
}

func TestDismountError(t *testing.T) {
	// not mounted
	hello := &Hello{}
	Dismount(hello)
}
