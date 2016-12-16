package markup

// import (
// 	"encoding/xml"
// 	"testing"
// )

// func TestMakeAttr(t *testing.T) {
// 	name := "Width"
// 	value := "42"

// 	xmlAttr := xml.Attr{
// 		Name:  xml.Name{Local: name},
// 		Value: value,
// 	}

// 	attr := makeAttr(xmlAttr)

// 	if attr.Name != name {
// 		t.Errorf("attr.Name should be %v: %v", name, attr.Name)
// 	}

// 	if attr.Value != value {
// 		t.Errorf("attr.Value should be %v: %v", value, attr.Value)
// 	}
// }

// func TestAttrAttr(t *testing.T) {
// 	attrInt := Attr{
// 		Name:  "int",
// 		Value: "42",
// 	}

// 	attrs := AttrList{
// 		attrInt,
// 	}

// 	attr, ok := attrs.Attr("int")
// 	if !ok {
// 		t.Error("attr name int should be found")
// 	}

// 	if attr != attrInt {
// 		t.Error("attr should be equal to attrInt")
// 	}

// 	if _, ok = attrs.Attr("foo"); ok {
// 		t.Error("attr should not be found")
// 	}
// }

// func TestAttrIsEvent(t *testing.T) {
// 	attr := Attr{
// 		Name: "_onchange",
// 	}

// 	if !attr.isEvent() {
// 		t.Error("attr should be an event")
// 	}

// 	attr.Name = "onchange"

// 	if attr.isEvent() {
// 		t.Error("attr should  not be an event")
// 	}
// }

// func TestMakeAttrList(t *testing.T) {
// 	xmlAttrs := []xml.Attr{
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Width"},
// 			Value: "42",
// 		},
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Height"},
// 			Value: "21",
// 		},
// 	}

// 	attrs := makeAttrList(xmlAttrs)

// 	if l := len(attrs); l != 2 {
// 		t.Fatal("attrs len should be 2:", l)
// 	}

// 	if name := attrs[0].Name; name != "Width" {
// 		t.Errorf("attrs[0].Name should be %v: %v", "Width", name)
// 	}

// 	if value := attrs[0].Value; value != "42" {
// 		t.Errorf("attrs[0].Value should be %v: %v", "42", value)
// 	}

// 	if name := attrs[1].Name; name != "Height" {
// 		t.Errorf("attrs[1].Name should be %v: %v", "Width", name)
// 	}

// 	if value := attrs[1].Value; value != "21" {
// 		t.Errorf("attrs[1].Value should be %v: %v", "42", value)
// 	}
// }

// func TestAttrListEquals(t *testing.T) {
// 	l1 := makeAttrList([]xml.Attr{
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Width"},
// 			Value: "42",
// 		},
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Height"},
// 			Value: "21",
// 		},
// 	})

// 	l1Bis := makeAttrList([]xml.Attr{
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Width"},
// 			Value: "42",
// 		},
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Height"},
// 			Value: "21",
// 		},
// 	})

// 	l2 := makeAttrList([]xml.Attr{
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Width"},
// 			Value: "42",
// 		},
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Height"},
// 			Value: "42",
// 		},
// 	})

// 	l3 := makeAttrList([]xml.Attr{
// 		xml.Attr{
// 			Name:  xml.Name{Local: "Height"},
// 			Value: "21",
// 		},
// 	})

// 	if !l1.equals(l1Bis) {
// 		t.Error("l1 and l1Bis should be equals")
// 	}

// 	if l1.equals(l2) {
// 		t.Error("l1 and l2 should not be equals")
// 	}

// 	if l1.equals(l3) {
// 		t.Error("l1 and l3 should not be equals")
// 	}

// }
