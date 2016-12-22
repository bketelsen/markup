package markup

import (
	"testing"

	"github.com/murlokswarm/uid"
)

type CompoSync struct {
	TextChange      bool
	HTMLAttrChange  bool
	HTMLTagChange   bool
	CompoAttrChange bool
	CompoChange     bool
	TypeChange      bool
	AddRemove       bool
}

func (c *CompoSync) Render() string {
	return `
<div>
    <!-- TextChange -->
    <p>{{if .TextChange}}Maxoo{{else}}Jonhy{{end}}</p>

    <!-- HTMLAttrChange -->
    <p class="{{if .HTMLAttrChange}}boo{{end}}">Say something</p>

    <!-- HTMLTagChange -->
    {{if .HTMLTagChange}}
        <h1>Hello</h1>
    {{else}}
        <h2>Hello</h2>
    {{end}}

    <!-- CompoAttrChange -->
    <SubCompoSync Name="{{if .CompoAttrChange}}Max{{else}}Maxence{{end}}" />

    <!-- CompoChange -->
    <div>
        {{if .CompoChange}}
            <SubCompoSyncBis />          
        {{else}}
            <SubCompoSync Name="Jonhzy" />
        {{end}}
    </div>

     <!-- TypeChange -->
    <div>
        {{if .TypeChange}}
            <div>
                <h1>I'm changed</h1>
                <SubCompoSyncBis /> 
            </div>           
        {{else}}
            <SubCompoSync Name="Bravo" />
        {{end}}
    </div>


    <!-- AddRemove -->
    <div>
        {{if .AddRemove}}<h1>Plop!</h1>{{end}}
    </div>
</div>
    `
}

type SubCompoSync struct {
	Name string
}

func (c *SubCompoSync) Render() string {
	return `
<div>
    <h1>{{html .Name}}</h1>
    <p>Whoa</p>
</div>
    `
}

type SubCompoSyncBis struct {
}

func (c *SubCompoSyncBis) Render() string {
	return `<p>I'm sexy</p>`
}

type CompoSyncError struct {
	BadTemplate bool
	BadRoot     bool
	BadMarkup   bool
}

func (c *CompoSyncError) Render() string {
	return `
{{if .BadTemplate}}
    {{.Unknown}}
{{else if .BadRoot}}
    <CompoBadRoot />
{{else if .BadMarkup}}
    <div></p>
{{else}}
    <div>Murloks!!!</div>
{{end}}
    `
}

func init() {
	Register(&CompoSync{})
	Register(&SubCompoSync{})
	Register(&SubCompoSyncBis{})
	Register(&CompoSyncError{})
}

func TestSynchronizeTextChange(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.TextChange = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}
}

func TestSynchronizeHTMLAttrChange(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.HTMLAttrChange = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != AttrSync {
		t.Error("s.Scope should be AttrSync")
	}

	if s.Attributes["class"] != "boo" {
		t.Error(`s.Attributes["class"] should be boo:`, s.Attributes["class"])
	}
}

func TestSynchronizeHTMLTagChange(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.HTMLTagChange = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}
}

func TestSynchronizeCompoAttrChange(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.CompoAttrChange = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}
}

func TestSynchronizeCompoChange(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.CompoChange = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}
}

func TestSynchronizeTypeChange(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.TypeChange = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}
}

func TestAddRemove(t *testing.T) {
	c := &CompoSync{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	// Add.
	c.AddRemove = true
	syncs := Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s := syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}

	// Remove.
	c.AddRemove = false
	syncs = Synchronize(c)

	if l := len(syncs); l != 1 {
		t.Error("l should be 1:", l)
	}

	s = syncs[0]
	t.Log(s.Node.Markup())

	if s.Scope != FullSync {
		t.Error("s.Scope should be FullSync")
	}
}

func TestSynchronizeBadTemplate(t *testing.T) {
	defer func() { recover() }()

	c := &CompoSyncError{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.BadTemplate = true
	Synchronize(c)
	t.Error("should panic")
}

func TestSynchronizeBadRoot(t *testing.T) {
	defer func() { recover() }()

	c := &CompoSyncError{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.BadRoot = true
	Synchronize(c)
	t.Error("should panic")
}

func TestSynchronizeBadMarkup(t *testing.T) {
	defer func() { recover() }()

	c := &CompoSyncError{}
	ctx := uid.Context()

	Mount(c, ctx)
	defer Dismount(c)

	c.BadMarkup = true
	Synchronize(c)
	t.Error("should panic")
}
