package ml

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestRender(t *testing.T) {
	ctx := uid.Context()
	hello := &Hello{}

	if err := Mount(hello, ctx); err != nil {
		t.Fatal(err)
	}
	defer Dismount(hello)

	// if err := Render(hello); err != nil {
	// 	t.Error(err)
	// }
}
