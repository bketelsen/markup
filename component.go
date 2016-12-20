package markup

import (
	"reflect"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

var (
	compoBuilders = map[string]func() Componer{}
	components    = map[Componer]*component{}
	nodes         = map[uid.ID]*Node{}
)

// Componer is the interface that describes a component.
type Componer interface {
	// Render should returns a markup.
	// The markup can be a template string following the text/template standard
	// package rules.
	Render() string
}

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
	log.Infof("%v has been registered under the tag %v", t, tag)
}

// Registered returns true if c is registered, otherwise false.
func Registered(c Componer) bool {
	v := reflect.Indirect(reflect.ValueOf(c))
	t := v.Type()
	_, registered := compoBuilders[t.Name()]
	return registered
}

// Root returns the root node of c.
func Root(c Componer) *Node {
	compo, mounted := components[c]
	if !mounted {
		log.Panicf("%T is not mounted", c)
	}
	return compo.Root
}

// Markup returns the markup of c.
func Markup(c Componer) string {
	return Root(c).Markup()
}

// Mount retains a component and its underlying nodes.
func Mount(c Componer, ctx uid.ID) (root *Node) {
	if !Registered(c) {
		log.Panicf("%T is not registered", c)
	}

	if compo, mounted := components[c]; mounted {
		// Go uses the same reference for different instances of a same empty struct.
		// This prevents from mounting a same empty struct.
		if t := reflect.TypeOf(c).Elem(); t.NumField() == 0 {
			compo.Count++
			return compo.Root
		}

		log.Panicf("%T is already mounted", c)
	}

	r, err := render(c)
	if err != nil {
		log.Panic(err)
	}

	if root, err = stringToNode(r); err != nil {
		log.Panicf("%T markup returned by Render() has a %v", c, err)
	}

	if root.Type != HTMLNode {
		log.Panicf("%T markup returned by Render() has a syntax error: root node is not a HTMLNode", c)
	}

	mountNode(root, c, ctx)
	components[c] = &component{
		Count: 1,
		Root:  root,
	}

	if mounter, isMounter := c.(Mounter); isMounter {
		mounter.OnMount()
	}
	return
}

func mountNode(n *Node, mount Componer, ctx uid.ID) {
	switch n.Type {
	case HTMLNode:
		mountHTMLNode(n, mount, ctx)

	case ComponentNode:
		mountComponentNode(n, mount, ctx)
	}
}

func mountHTMLNode(n *Node, mount Componer, ctx uid.ID) {
	n.ID = uid.Elem()
	n.ContextID = ctx
	n.Mount = mount
	nodes[n.ID] = n

	for _, c := range n.Children {
		mountNode(c, mount, ctx)
	}
}

func mountComponentNode(n *Node, mount Componer, ctx uid.ID) {
	n.ContextID = ctx
	n.Mount = mount

	b, registed := compoBuilders[n.Tag]
	if !registed {
		log.Panicf("%v is not registered", n.Tag)
	}

	c := b()
	decodeAttributeMap(n.Attributes, c)
	Mount(c, ctx)
	n.Component = c
}

// Dismount dismounts a component.
func Dismount(c Componer) {
	compo, mounted := components[c]
	if !mounted {
		return
	}

	// Go uses the same reference for different instances of a same empty struct.
	// This prevents from dismounting an empty struct that still remains in another context.
	if compo.Count--; compo.Count == 0 {
		dismountNode(compo.Root)
		delete(components, c)

		if dismounter, isDismounter := c.(Dismounter); isDismounter {
			dismounter.OnDismount()
		}
	}
	return
}

func dismountNode(n *Node) {
	switch n.Type {
	case HTMLNode:
		dismountHTMLNode(n)

	case ComponentNode:
		Dismount(n.Component)
	}
}

func dismountHTMLNode(n *Node) {
	for _, c := range n.Children {
		dismountNode(c)
	}

	delete(nodes, n.ID)
}
