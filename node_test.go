package markup

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestNodeString(t *testing.T) {
	t.Log(Node{
		ID:  uuid.NewV1(),
		Tag: "div",
	})
}
