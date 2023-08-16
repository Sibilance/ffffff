package macro

import "fmt"

type ValueKind uint8

const (
	NullValueKind ValueKind = iota
	BoolValueKind
	IntValueKind
	FloatValueKind
	StringValueKind
	BytesValueKind
	ListValueKind
	StringMapValueKind
	ValueMapValueKind
)

type Value struct {
	Kind           ValueKind
	BoolValue      bool
	IntValue       int64
	FloatValue     float64
	StringValue    string
	BytesValue     []byte
	ListValue      []*Value
	StringMapValue map[string]*Value
	ValueMapValue  map[*Value]*Value
}

func NullValue() *Value {
	return &Value{}
}

func BoolValue(v bool) *Value {
	return &Value{Kind: BoolValueKind, BoolValue: v}
}

func IntValue(v int64) *Value {
	return &Value{Kind: IntValueKind, IntValue: v}
}

func FloatValue(v float64) *Value {
	return &Value{Kind: FloatValueKind, FloatValue: v}
}

func StringValue(v string) *Value {
	return &Value{Kind: StringValueKind, StringValue: v}
}

func BytesValue(v []byte) *Value {
	return &Value{Kind: BytesValueKind, BytesValue: v}
}

func ListValue(v []*Value) *Value {
	return &Value{Kind: ListValueKind, ListValue: v}
}

func StringMapValue(v map[string]*Value) *Value {
	return &Value{Kind: StringMapValueKind, StringMapValue: v}
}

func ValueMapValue(v map[*Value]*Value) *Value {
	return &Value{Kind: ValueMapValueKind, ValueMapValue: v}
}

func (v *Value) Bool() bool {
	switch v.Kind {
	case NullValueKind:
		return false
	case BoolValueKind:
		return v.BoolValue
	case IntValueKind:
		return v.IntValue != 0
	case FloatValueKind:
		return v.FloatValue != 0
	case StringValueKind:
		return len(v.StringValue) > 0
	case BytesValueKind:
		return len(v.BytesValue) > 0
	case ListValueKind:
		return len(v.ListValue) > 0
	case StringMapValueKind:
		return len(v.StringMapValue) > 0
	case ValueMapValueKind:
		return len(v.ValueMapValue) > 0
	default:
		panic(fmt.Sprintf("unexpected ValueKind %d", v.Kind))
	}
}
