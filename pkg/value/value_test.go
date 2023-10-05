package value

import (
	"errors"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func assertBool(t *testing.T, v Value, b bool) {
	if v.Bool() != b {
		t.Fatalf("expected %v, got %v", b, v.Bool())
	}
}

func assertString(t *testing.T, v Value, s string) {
	if v.String() != s {
		t.Fatalf("expected %v, got %v", s, v.String())
	}
}

func assertMarshalYAML(t *testing.T, v Value, s string) {
	s = strings.TrimSpace(s)
	actual, err := yaml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	actualString := strings.TrimSpace(string(actual))
	if actualString != s {
		t.Fatalf("expected:\n%v\ngot:\n%v", s, actualString)
	}
}

func assertUnmarshalYAML(t *testing.T, v Value, s string, expected Value, expectedErr error) {
	err := yaml.Unmarshal([]byte(s), v)
	if err != nil {
		if expectedErr != nil && err.Error() != expectedErr.Error() {
			t.Fatalf("expected error `%v`, got error `%v`", expectedErr, err)
		} else if expectedErr != nil {
			return
		}
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, expected) {
		t.Fatalf("expected:\n%v\ngot:\n%v", expected, v)
	}
}

func TestNullValue(t *testing.T) {
	assertBool(t, NullValue{}, false)
	assertString(t, NullValue{}, "null")
	assertMarshalYAML(t, NullValue{}, `null`)
	assertUnmarshalYAML(t, NullValue{}, `null`, NullValue{}, nil)
	assertUnmarshalYAML(t, &NullValue{}, `0`, nil,
		errors.New("cannot unmarshal !!int into null"))
}

func TestBoolValue(t *testing.T) {
	assertBool(t, &BoolValue{false}, false)
	assertBool(t, &BoolValue{true}, true)
	assertString(t, &BoolValue{false}, "false")
	assertString(t, &BoolValue{true}, "true")
	assertMarshalYAML(t, &BoolValue{false}, `false`)
	assertMarshalYAML(t, &BoolValue{true}, `true`)
	assertUnmarshalYAML(t, &BoolValue{}, `false`, &BoolValue{false}, nil)
	assertUnmarshalYAML(t, &BoolValue{}, `true`, &BoolValue{true}, nil)
	assertUnmarshalYAML(t, &BoolValue{}, `0`, nil,
		errors.New("cannot unmarshal !!int into bool"))
}

func TestIntValue(t *testing.T) {
	assertBool(t, &IntValue{*big.NewInt(0)}, false)
	assertBool(t, &IntValue{*big.NewInt(1)}, true)
	assertString(t, &IntValue{*big.NewInt(123)}, "123")
	intValue, _ := (&big.Int{}).SetString("-12345678901234567890", 0)
	assertString(t, &IntValue{*intValue}, "-12345678901234567890")
	assertMarshalYAML(t, &IntValue{*big.NewInt(123)}, `123`)
	assertMarshalYAML(t, &IntValue{*intValue}, `!!int -12345678901234567890`)
	assertUnmarshalYAML(t, &IntValue{}, `-12345678901234567890`, &IntValue{*intValue}, nil)
	assertUnmarshalYAML(t, &IntValue{}, `123`, &IntValue{*big.NewInt(123)}, nil)
	assertUnmarshalYAML(t, &IntValue{}, `0xff`, &IntValue{*big.NewInt(255)}, nil)
	assertUnmarshalYAML(t, &IntValue{}, `0o77`, &IntValue{*big.NewInt(63)}, nil)
	assertUnmarshalYAML(t, &IntValue{}, `abc`, nil,
		errors.New("cannot unmarshal !!str into int"))
	assertUnmarshalYAML(t, &IntValue{}, `!!int abc`, nil,
		errors.New("cannot unmarshal value \"abc\" into int"))
}

func TestFloatValue(t *testing.T) {
	assertBool(t, &FloatValue{0.0}, false)
	assertBool(t, &FloatValue{0.1}, true)
	assertString(t, &FloatValue{3.14159}, "3.14159")
	assertString(t, &FloatValue{6.022e23}, "6.022e+23")
	assertString(t, &FloatValue{-1.6e-19}, "-1.6e-19")
	assertMarshalYAML(t, &FloatValue{0.0}, `!!float 0`)
	assertMarshalYAML(t, &FloatValue{3.14159}, `3.14159`)
	assertMarshalYAML(t, &FloatValue{6.022e23}, `6.022e+23`)
	assertMarshalYAML(t, &FloatValue{-1.6e-19}, `-1.6e-19`)
	assertUnmarshalYAML(t, &FloatValue{}, `!!float 0`, &FloatValue{0.0}, nil)
	assertUnmarshalYAML(t, &FloatValue{}, `0.0`, &FloatValue{0.0}, nil)
	assertUnmarshalYAML(t, &FloatValue{}, `3.14159`, &FloatValue{3.14159}, nil)
	assertUnmarshalYAML(t, &FloatValue{}, `6.022e23`, &FloatValue{6.022e23}, nil)
	assertUnmarshalYAML(t, &FloatValue{}, `-1.6e-19`, &FloatValue{-1.6e-19}, nil)
	assertUnmarshalYAML(t, &FloatValue{}, `0`, nil,
		errors.New("cannot unmarshal !!int into float"))
	assertUnmarshalYAML(t, &FloatValue{}, `!!str 3.14159`, nil,
		errors.New("cannot unmarshal !!str into float"))
}

func TestStringValue(t *testing.T) {
	assertBool(t, &StringValue{""}, false)
	assertBool(t, &StringValue{"hello world"}, true)
	assertString(t, &StringValue{""}, "")
	assertString(t, &StringValue{"hello world"}, "hello world")
	assertMarshalYAML(t, &StringValue{""}, `""`)
	assertMarshalYAML(t, &StringValue{"hello world"}, `hello world`)
	assertUnmarshalYAML(t, &StringValue{}, `hello world`, &StringValue{"hello world"}, nil)
	assertUnmarshalYAML(t, &StringValue{}, `!!str 123`, &StringValue{"123"}, nil)
	assertUnmarshalYAML(t, &StringValue{}, `123`, nil,
		errors.New("cannot unmarshal !!int into string"))
}

func TestListValue(t *testing.T) {
	assertBool(t, &ListValue{[]Value{}}, false)
	assertBool(t, &ListValue{[]Value{&StringValue{""}}}, true)
	assertString(t, &ListValue{[]Value{}}, "[]\n")
	assertString(t, &ListValue{[]Value{&StringValue{""}}}, "- \"\"\n")
	assertMarshalYAML(t, &ListValue{[]Value{}}, "[]\n")
	assertMarshalYAML(t, &ListValue{[]Value{&StringValue{""}}}, "- \"\"\n")
	assertUnmarshalYAML(t, &ListValue{}, `[]`, &ListValue{}, nil)
	assertUnmarshalYAML(t, &ListValue{}, `[""]`, &ListValue{[]Value{&StringValue{""}}}, nil)
	assertUnmarshalYAML(t, &ListValue{}, `[1]`, &ListValue{[]Value{&IntValue{*big.NewInt(1)}}}, nil)
	intValue, _ := (&big.Int{}).SetString("-12345678901234567890", 0)
	assertUnmarshalYAML(t, &ListValue{}, `[-12345678901234567890]`, &ListValue{[]Value{&IntValue{*intValue}}}, nil)
	assertUnmarshalYAML(t, &ListValue{}, `[3.14159]`, &ListValue{[]Value{&FloatValue{3.14159}}}, nil)
	assertUnmarshalYAML(t, &ListValue{}, `[6.022e23]`, &ListValue{[]Value{&FloatValue{6.022e23}}}, nil)
	assertUnmarshalYAML(t, &ListValue{}, `[-1.6e-19]`, &ListValue{[]Value{&FloatValue{-1.6e-19}}}, nil)
	assertUnmarshalYAML(t, &ListValue{}, `[[]]`, &ListValue{[]Value{&ListValue{}}}, nil)
	// TODO: test nested map
}
