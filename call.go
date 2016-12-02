package markup

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
)

// Call invokes a method from the component which own the element associated to
// elementID.
func Call(elementID uid.ID, componentMethodName string, jsonArg string) error {
	elem, mounted := elements[elementID]
	if !mounted {
		return fmt.Errorf("element with id %v is not mounted", elementID)
	}

	component := elem.Component
	componentValue := reflect.ValueOf(component)
	methodValue := componentValue.MethodByName(componentMethodName)

	if !methodValue.IsValid() {
		log.Warnf("%T doesn't have a method named %v", component, componentMethodName)
		return nil
	}

	methodType := methodValue.Type()

	switch numIn := methodType.NumIn(); numIn {
	case 0:
		methodValue.Call([]reflect.Value{})
		return nil

	case 1:
		argType := methodType.In(0)

		arg, err := createCallArg(argType, jsonArg)
		if err != nil {
			return err
		}

		methodValue.Call([]reflect.Value{arg})
		return nil

	default:
		return fmt.Errorf("%T.%v must have 1 parameter max: %v", component, componentMethodName, methodType)
	}
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
