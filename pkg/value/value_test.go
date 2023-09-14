package value

import (
	"errors"
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
