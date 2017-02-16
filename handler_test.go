package markup

type CompoWithFunc struct {
	Name             string
	calledWithNoArg  bool
	calledWithOneArg bool
}

func (c *CompoWithFunc) OnCallTest() {
	c.calledWithNoArg = true
}

func (c *CompoWithFunc) OnCallTestWithArg(arg string) {
	c.calledWithOneArg = true
}

func (c *CompoWithFunc) OnCallTestWithMultipleArgs(arg int, number int) {
}

func (c *CompoWithFunc) Render() string {
	return `
<h1>{{.Name}}</h1>
    `
}

type FuncArg struct {
	Number int
	String string
}

func init() {
	Register(&CompoWithFunc{})
}

// func TestCall(t *testing.T) {
// 	c := &CompoWithFunc{}
// 	ctx := uid.Context()

// 	Mount(c, ctx)
// 	defer Dismount(c)

// 	root := Root(c)
// 	Call(root.ID, "OnCallTest", "")
// 	Call(root.ID, "OnCallTestWithArg", `"42"`)
// 	Call(root.ID, "Nonexistent", "")

// 	if !c.calledWithNoArg {
// 		t.Error("OnCallTest should have been called")
// 	}

// 	if !c.calledWithOneArg {
// 		t.Error("OnCallTestWithArg should have been called")
// 	}
// }

// func TestCallNotMountedNode(t *testing.T) {
// 	Call("elem--42", "", "")
// }

// func TestCallMethodWithMultipleArgs(t *testing.T) {
// 	c := &CompoWithFunc{}
// 	ctx := uid.Context()

// 	Mount(c, ctx)
// 	defer Dismount(c)

// 	root := Root(c)
// 	Call(root.ID, "OnCallTestWithMultipleArgs", "")
// }
