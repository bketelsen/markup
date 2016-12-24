package markup

import (
	"testing"

	"github.com/murlokswarm/uid"
)

type CompoMount struct {
	mounted             bool
	EmbedsNonRegistered bool
}

func (c *CompoMount) OnMount() {
	c.mounted = true
}

func (c *CompoMount) OnDismount() {
	c.mounted = false
}

func (c *CompoMount) Render() string {
	return `
<div>
    CompoMount is mounted
    <CompoEmpty />

    {{if .EmbedsNonRegistered}}
        <CompoNotRegistered />
    {{end}}
</div>
    `
}

type CompoEmpty struct{}

func (c *CompoEmpty) Render() string {
	return `<p>CompoEmpty is mounted</p>`
}

type CompoNotRegistered struct{}

func (c *CompoNotRegistered) Render() string {
	return `<p>CompoNotRegistered</p>`
}

type CompoBadRenderTemplate struct{}

func (c *CompoBadRenderTemplate) Render() string {
	return `<p>CompoBadRender {{.Foo}}</p>`
}

type CompoBadMarkup struct{}

func (c *CompoBadMarkup) Render() string {
	return `<p>CompoBadMarkup</span>`
}

type CompoBadRoot struct{}

func (c *CompoBadRoot) Render() string {
	return `<CompoEmpty />`
}

type compoNotExported struct{}

func (c *compoNotExported) Render() string {
	return `<p>CompoNotExported</p>`
}

type CompoNoPtr struct{}

func (c CompoNoPtr) Render() string {
	return `<p>CompoNoPtr</p>`
}

func init() {
	Register(&CompoMount{})
	Register(&CompoEmpty{})
	Register(&CompoBadRenderTemplate{})
	Register(&CompoBadMarkup{})
	Register(&CompoBadRoot{})
}

func TestRegisterNotExported(t *testing.T) {
	defer func() { recover() }()
	Register(&compoNotExported{})
	t.Error("should panic")
}

func TestRegisterNoPtr(t *testing.T) {
	defer func() { recover() }()
	Register(CompoNoPtr{})
	t.Error("should panic")

}

func TestRootNotMounted(t *testing.T) {
	defer func() { recover() }()
	Root(&CompoEmpty{})
	t.Error("should panic")
}

func TestMount(t *testing.T) {
	ctx := uid.Context()
	c := &CompoMount{}
	root := Mount(c, ctx)
	t.Log(root)

	if l := len(components); l != 2 {
		t.Error("components len should be 2", l)
	}

	if l := len(nodes); l != 2 {
		t.Error("node len should be 2", l)
	}

	if !c.mounted {
		t.Error("c.mounted should be true:", c.mounted)
	}

	t.Log(Markup(c))

	Dismount(c)

	if l := len(components); l != 0 {
		t.Error("components len should be 0", l)
	}

	if l := len(nodes); l != 0 {
		t.Error("node len should be 0", l)
	}

	if c.mounted {
		t.Error("c.mounted should be false:", c.mounted)
	}

	Dismount(c)
}

func TestMountNotRegistered(t *testing.T) {
	defer func() { recover() }()

	ctx := uid.Context()
	c := &CompoNotRegistered{}
	Mount(c, ctx)
	t.Error("should panic")
}

func TestMountEmbedsNotRegistered(t *testing.T) {
	defer func() { recover() }()

	ctx := uid.Context()
	c := &CompoMount{EmbedsNonRegistered: true}
	Mount(c, ctx)
	t.Error("should panic")
}

func TestMountAlreadyMounted(t *testing.T) {
	defer func() { recover() }()

	ctx := uid.Context()
	c := &CompoMount{}
	Mount(c, ctx)
	Mount(c, ctx)
	t.Error("should panic")
}

func TestMountBadRenderTemplate(t *testing.T) {
	ctx := uid.Context()
	c := &CompoBadRenderTemplate{}
	Mount(c, ctx)
}

func TestMountBadMarkup(t *testing.T) {
	ctx := uid.Context()
	c := &CompoBadMarkup{}
	Mount(c, ctx)
}

func TestMountBadRoot(t *testing.T) {
	ctx := uid.Context()
	c := &CompoBadRoot{}
	Mount(c, ctx)
}
