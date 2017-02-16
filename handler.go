package markup

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/uid"
	"github.com/pkg/errors"
)

type changeArg struct {
	Value string
}

func HandleHTMLEvent(nodeID uid.ID, name string, argJSON string) {
	var err error

	if len(name) == 0 {
		return
	}

	n, mounted := nodes[nodeID]
	if !mounted {
		log.Error(errors.Errorf("node with ID = %v does not belong to a mounted component.", nodeID))
		return
	}

	c := n.Mount
	v := reflect.ValueOf(c)

	if m := v.MethodByName(name); m.IsValid() {
		t := m.Type()
		target := fmt.Sprintf("%v.%v", t.Name(), name)
		callComponentMethod(target, m, argJSON)
		return
	}

	if f := v.FieldByName(name); f.IsValid() {
		if f, err = getFieldValue(f, strings.Split(name, ".")); err != nil {
			err = errors.Wrapf(err, "unable to map %v", name)
			log.Error(err)
			return
		}

		if err = mapComponentField(f, argJSON); err != nil {
			err = errors.Wrapf(err, "unable to map %v", name)
			log.Error(err)
		}
		return
	}

	log.Warnf("%T doesn't have a method or a field named %v", c, name)
}

func callComponentMethod(target string, m reflect.Value, argJSON string) {
	t := m.Type()
	if t.NumIn() == 0 {
		m.Call([]reflect.Value{})
		return
	}

	argt := t.In(0)
	argv := reflect.New(argt)
	argi := argv.Interface()
	arg := argv.Elem()

	if err := json.Unmarshal([]byte(argJSON), argi); err != nil {
		log.Error(errors.Wrapf(err, "callComponentMethod for %v", target))
		return
	}

	m.Call([]reflect.Value{arg})

	if t.NumIn() > 1 {
		log.Warnf("%v should have only 1 argument", target)
	}
}

func getFieldValue(f reflect.Value, pipeline []string) (v reflect.Value, err error) {
	if len(pipeline) == 1 {
		v = f
		return
	}

	switch k := f.Kind(); k {
	case reflect.Ptr:
		if f.IsNil() {
			newv := reflect.New(f.Type().Elem())
			f.Set(newv)
		}
		return getFieldValue(f, pipeline)

	case reflect.Struct:
		f := f.FieldByName(pipeline[1])
		if !f.IsValid() {
			err = errors.Errorf("no field named %v in %v", pipeline[1], pipeline[0])
			return
		}
		return getFieldValue(f, pipeline[1:])

	case reflect.Map:
		mapt := f.Type()
		if keyk := mapt.Key().Kind(); keyk != reflect.String {
			err = errors.Errorf("%v key type must be a %v: %v",
				pipeline[0],
				reflect.String,
				keyk)
			return
		}

		if f.IsNil() {
			newv := reflect.MakeMap(mapt)
			f.Set(newv)
		}

		mapvv := reflect.Zero(mapt.Elem())
		f.SetMapIndex(reflect.ValueOf(pipeline[1]), mapvv)
		return getFieldValue(f, pipeline[1:])

	default:
		err = errors.Errorf("%v must be a %v or a %v: %v",
			pipeline[0],
			reflect.Struct,
			reflect.Map,
			k)
		return
	}
}

func mapComponentField(f reflect.Value, argJSON string) error {
	var arg changeArg
	if err := json.Unmarshal([]byte(argJSON), &arg); err != nil {
		return errors.Wrapf(err, "unable to unmarshal %v", argJSON)
	}

	if f.Kind() == reflect.String {
		f.SetString(arg.Value)
		return nil
	}

	ft := f.Type()
	fpv := reflect.New(ft)
	fpi := fpv.Interface()

	if err := json.Unmarshal([]byte(arg.Value), fpi); err != nil {
		return errors.Wrapf(err, "unable to unmarshal %v", arg.Value)
	}

	f.Set(fpv.Elem())
	return nil
}
