package markup

import (
	"testing"

	"github.com/murlokswarm/uid"
)

func TestNodeString(t *testing.T) {
	t.Log(Node{
		ID:  uid.Elem(),
		Tag: "div",
	})
}
