package macro

import (
	"testing"
)

func TestValueBool(t *testing.T) {
	if NullValue().Bool() != false {
		t.Fatal("NullValue expected to be false")
	}
	if BoolValue(false).Bool() != false {
		t.Fatal("BoolValue{false} expected to be false")
	}
	if BoolValue(true).Bool() != true {
		t.Fatal("BoolValue{true} expected to be true")
	}
	if IntValue(0).Bool() != false {
		t.Fatal("IntValue{0} expected to be false")
	}
	if IntValue(1).Bool() != true {
		t.Fatal("IntValue{1} expected to be true")
	}
	if IntValue(-1).Bool() != true {
		t.Fatal("IntValue{-1} expected to be true")
	}
	if FloatValue(0).Bool() != false {
		t.Fatal("FloatValue{0} expected to be false")
	}
	if FloatValue(3.14).Bool() != true {
		t.Fatal("FloatValue{3.14} expected to be true")
	}
	if FloatValue(-1.0/12).Bool() != true {
		t.Fatal("FloatValue{-1.0/12} expected to be true")
	}
	if StringValue("").Bool() != false {
		t.Fatal("StringValue{\"\"} expected to be false")
	}
	if StringValue("hello").Bool() != true {
		t.Fatal("StringValue{\"hello\"} expected to be true")
	}
	if BytesValue([]byte("")).Bool() != false {
		t.Fatal("BytesValue{\"\"} expected to be false")
	}
	if BytesValue([]byte("hello")).Bool() != true {
		t.Fatal("BytesValue{\"hello\"} expected to be true")
	}
	if ListValue([]*Value{}).Bool() != false {
		t.Fatal("ListValue{} expected to be false")
	}
	if ListValue([]*Value{NullValue()}).Bool() != true {
		t.Fatal("ListValue{NullValue} expected to be true")
	}
	if StringMapValue(map[string]*Value{}).Bool() != false {
		t.Fatal("StringMapValue{} expected to be false")
	}
	// StringMapValue
	// ValueMapValue
}
