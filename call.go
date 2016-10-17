package ml

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

// Call invokes a method from the component which own the element associated to
// elementID.
func Call(elementID uid.ID, componentMethodName string, jsonArg string) (err error) {
	var mounted bool
	var element *Element

	if element, mounted = elements[elementID]; !mounted {
		err = fmt.Errorf("element with id %v is not mounted", elementID)
		return
	}

	component := element.Component
	componentValue := reflect.ValueOf(component)
	methodValue := componentValue.MethodByName(componentMethodName)

	if !methodValue.IsValid() {
		log.Warnf("%T doesn't have a method named %v", componentMethodName)
		return
	}

	methodType := methodValue.Type()

	switch numIn := methodType.NumIn(); numIn {
	case 0:
		methodValue.Call([]reflect.Value{})

	case 1:
		argType := methodType.In(0)
		arg := reflect.Zero(argType)
		methodValue.Call([]reflect.Value{arg})

	default:
		log.Errorf("%T.%v must have 1 parameter max: %v", component, componentMethodName, methodType)
	}

	return
}

func createCallArg(t reflect.Type, jsonArg string) (arg reflect.Value, err error) {
	arg = reflect.New(t)
	i := arg.Interface()

	if err = json.Unmarshal([]byte(jsonArg), i); err != nil {
		return
	}

	arg = arg.Elem()
	return
}
