package value

import (
	"fmt"
	"math/big"

	"gopkg.in/yaml.v3"
)

type Value interface {
	Bool() bool
	String() string
	MarshalYAML() (interface{}, error)
	UnmarshalYAML(*yaml.Node) error
}

type NullValue struct{}

func (v NullValue) Bool() bool {
	return false
}

func (v NullValue) String() string {
	return "null"
}

func (v NullValue) MarshalYAML() (interface{}, error) {
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

func (v BoolValue) MarshalYAML() (interface{}, error) {
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

type StringValue struct {
	value string
}

func (v StringValue) Bool() bool {
	return len(v.value) > 0
}

func (v StringValue) String() string {
	return v.value
}

func (v StringValue) MarshalYAML() (interface{}, error) {
	return v.value, nil
}

func (v *StringValue) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&v.value)
}
