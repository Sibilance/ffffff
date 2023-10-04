package value

import (
	"fmt"
	"math/big"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Value interface {
	Bool() bool
	String() string
	MarshalYAML() (any, error)
	UnmarshalYAML(*yaml.Node) error
}

type NullValue struct{}

func (v NullValue) Bool() bool {
	return false
}

func (v NullValue) String() string {
	return "null"
}

func (v NullValue) MarshalYAML() (any, error) {
	return nil, nil
}

func (v NullValue) UnmarshalYAML(node *yaml.Node) error {
	if node.ShortTag() != "!!null" {
		return fmt.Errorf("cannot unmarshal %s into null", node.ShortTag())
	}
	return nil
}

type BoolValue struct {
	value bool
}

func (v BoolValue) Bool() bool {
	return v.value
}

func (v BoolValue) String() string {
	if v.value {
		return "true"
	}
	return "false"
}

func (v BoolValue) MarshalYAML() (any, error) {
	return v.value, nil
}

func (v *BoolValue) UnmarshalYAML(node *yaml.Node) error {
	if node.ShortTag() != "!!bool" {
		return fmt.Errorf("cannot unmarshal %s into bool", node.ShortTag())
	}
	return node.Decode(&v.value)
}

type IntValue struct {
	value big.Int
}

func (v *IntValue) Bool() bool {
	return v.value.Sign() != 0
}

func (v *IntValue) String() string {
	return v.value.String()
}

func (v *IntValue) MarshalYAML() (any, error) {
	n := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: v.String(),
	}
	return n, nil
}

func (v *IntValue) UnmarshalYAML(node *yaml.Node) error {
	// Yaml library can misidentify long integers as floats.
	if node.ShortTag() != "!!int" && node.ShortTag() != "!!float" {
		return fmt.Errorf("cannot unmarshal %s into int", node.ShortTag())
	}
	_, success := v.value.SetString(node.Value, 0)
	if !success {
		return fmt.Errorf("cannot unmarshal value \"%s\" into int", node.Value)
	}
	return nil
}

type FloatValue struct {
	value float64
}

func (v FloatValue) Bool() bool {
	return v.value != 0
}

func (v FloatValue) String() string {
	return strconv.FormatFloat(v.value, 'g', -1, 64)
}

func (v FloatValue) MarshalYAML() (any, error) {
	n := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!float",
		Value: v.String(),
	}
	return n, nil
}

func (v *FloatValue) UnmarshalYAML(node *yaml.Node) error {
	if node.ShortTag() != "!!float" {
		return fmt.Errorf("cannot unmarshal %s into float", node.ShortTag())
	}
	return node.Decode(&v.value)
}

type StringValue struct {
	value string
}

func (v StringValue) Bool() bool {
	return len(v.value) > 0
}

func (v StringValue) String() string {
	return v.value
}

func (v StringValue) MarshalYAML() (any, error) {
	return v.value, nil
}

func (v *StringValue) UnmarshalYAML(node *yaml.Node) error {
	if node.ShortTag() != "!!str" {
		return fmt.Errorf("cannot unmarshal %s into string", node.ShortTag())
	}
	return node.Decode(&v.value)
}
