package markup

import (
	"encoding/json"
	"reflect"

	"github.com/murlokswarm/errors"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

// Call invokes the method named by methodName from the component which own
// the node associated to nodeID.
func Call(nodeID uid.ID, methodName string, argJSON string) {
	n, mounted := nodes[nodeID]
	if !mounted {
		log.Error(errors.Newf("node with ID = %v does not belong to a mounted component.", nodeID))
		return
	}

	c := n.Mount
	v := reflect.ValueOf(c)
	m := v.MethodByName(methodName)

	if !m.IsValid() {
		log.Warnf("%T doesn't have a method named %v", c, methodName)
		return
	}

	mt := m.Type()

	switch numIn := mt.NumIn(); numIn {
	case 0:
		m.Call([]reflect.Value{})
		return

	case 1:
		at := mt.In(0)
		arg := createCallArg(at, argJSON)
		m.Call([]reflect.Value{arg})
		return

	default:
		log.Error(errors.Newf("%T %v should have 1 parameter max: %v", c, methodName, mt))
		return
	}
}

func createCallArg(t reflect.Type, argJSON string) reflect.Value {
	arg := reflect.New(t)
	i := arg.Interface()
	json.Unmarshal([]byte(argJSON), i)
	return arg.Elem()
}
