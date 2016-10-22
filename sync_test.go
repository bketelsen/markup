package markup

import (
	"testing"

	"github.com/murlokswarm/uid"
)

const (
	testHTML            TestType = 0
	testHTMLAlt                  = 1
	testComponent                = 2
	testComponentAlt             = 3
	testText                     = 4
	testChildMountError          = 5
	testBadMarkup                = 6
	testBadTemplate              = 7
)

type TestType uint8

type SyncComponent struct {
	Input                             string
	Text                              string
	Number                            int
	TestType                          TestType
	TestHTMLToText                    bool
	TestComponentToDifferentComponent bool
	TestComponentToText               bool
}

func (c *SyncComponent) Render() string {
	return `
<div>
	<h1>SyncComponent</h1>
	<div>
		==========  TEST BEGIN ==========
		<!-- Test begin -->
		{{if eq .TestType 0}}
			<input value="{{.Input}}" />
		{{end}}

		{{if eq .TestType 1}}
			<div>
				<p>Booboo</p>
			</div>
		{{end}}

		{{if eq .TestType 2}}
			<SubSyncComponent Number="{{.Number}}" />
		{{end}}

		{{if eq .TestType 3}}
			<Hello />
		{{end}}

		{{if eq .TestType 4}}
			{{.Text}}
		{{end}}

		{{if eq .TestType 5}}
			<div>
				<Hello MountError="true" />
			</div>
		{{end}}

		{{if eq .TestType 6}}
			<div></p>
		{{end}}

		{{if eq .TestType 7}}
			<div>{{.Boo}}</div>
		{{end}}
		<!-- Test end -->
		===========  TEST END ===========
	</div>
</div>
	`
}

type SubSyncComponent struct {
	Number uint
}

func (c *SubSyncComponent) Render() string {
	return `
<div>
	<h1>SubSyncComponent</h1>
	<p>{{.Number}}</p>
</div>
	`
}

func init() {
	RegisterComponent("SyncComponent", func() Componer {
		return &SyncComponent{}
	})

	RegisterComponent("SubSyncComponent", func() Componer {
		return &SubSyncComponent{}
	})
}

func TestSyncNoChange(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testHTML}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 0 {
		t.Error("should not have changed elements")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncHTMLToHTML(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testHTMLAlt}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testHTML

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncHTMLToComponent(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testHTML}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testComponent

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncHTMLToText(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testHTML}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testText
	c.Text = "La vie en rose"

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncComponentToComponent(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testComponent}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.Number = 42

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncComponentToDifferentComponent(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testComponent}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testComponentAlt

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncComponentToHTML(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testComponent}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testHTML

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncComponentToText(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testComponent}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testText
	c.Text = "Boo"

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncTextToText(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testText, Text: "Hello"}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.Text = "World"

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncTextToHTML(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testText, Text: "Hello"}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}

	c.TestType = testHTML

	changed, err := Sync(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 1 {
		t.Error("changed should be equal to 1")
	}

	if HTML, err := ComponentToHTML(c); err == nil {
		t.Log(HTML)
	}
}

func TestSyncChildMountError(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{Input: "Hello"}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	c.TestType = testChildMountError

	if _, err := Sync(c); err == nil {
		t.Error("should error")
	}
}

func TestSyncAttrErrorError(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{TestType: testComponent}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	c.Number = -42

	if _, err := Sync(c); err == nil {
		t.Error("should error")
	}
}

func TestSyncNotMounted(t *testing.T) {
	c := &SyncComponent{TestType: testComponent}

	if _, err := Sync(c); err == nil {
		t.Error("should error")
	}
}

func TestSyncBadMarkup(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	c.TestType = testBadMarkup

	if _, err := Sync(c); err == nil {
		t.Error("should error")
	}
}

func TestSyncBadTemplate(t *testing.T) {
	ctx := uid.Context()
	c := &SyncComponent{}

	if err := Mount(c, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(c)

	c.TestType = testBadTemplate

	if _, err := Sync(c); err == nil {
		t.Error("should error")
	}
}
