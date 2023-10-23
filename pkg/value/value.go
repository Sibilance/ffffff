package value

import (
	"encoding/binary"
	"fmt"
	"hash/maphash"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

type Value interface {
	Bool() bool
	String() string
	Hash(maphash.Seed) (uint64, error)
	Cmp(Value) int
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

func (v NullValue) Hash(seed maphash.Seed) (uint64, error) {
	return maphash.String(seed, ""), nil
}

func (v NullValue) Cmp(other Value) int {
	switch other.(type) {
	case NullValue:
		return 0
	}
	return cmpType(v, other)
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

func (v BoolValue) Hash(seed maphash.Seed) (uint64, error) {
	if v.value {
		return maphash.String(seed, "T"), nil
	} else {
		return maphash.String(seed, "F"), nil
	}
}

func (v BoolValue) Cmp(other Value) int {
	switch other := other.(type) {
	case *BoolValue:
		if !v.value && other.value {
			return -1
		}
		if v.value && !other.value {
			return 1
		}
		return 0
	}
	return cmpType(&v, other)
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

func (v *IntValue) Hash(seed maphash.Seed) (uint64, error) {
	return maphash.Bytes(seed, v.value.Bytes()), nil
}

func (v *IntValue) Cmp(other Value) int {
	switch other := other.(type) {
	case *IntValue:
		return v.value.Cmp(&other.value)
	}
	return cmpType(v, other)
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

func (v FloatValue) Hash(seed maphash.Seed) (uint64, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(v.value))
	return maphash.Bytes(seed, buf[:]), nil
}

func (v FloatValue) Cmp(other Value) int {
	switch other := other.(type) {
	case *FloatValue:
		if v.value < other.value {
			return -1
		}
		if v.value > other.value {
			return 1
		}
		return 0
	}
	return cmpType(&v, other)
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

func (v StringValue) Hash(seed maphash.Seed) (uint64, error) {
	return maphash.String(seed, v.value), nil
}

func (v StringValue) Cmp(other Value) int {
	switch other := other.(type) {
	case *StringValue:
		if v.value < other.value {
			return -1
		}
		if v.value > other.value {
			return 1
		}
		return 0
	}
	return cmpType(&v, other)
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

type ListValue struct {
	value []Value
}

func (v ListValue) Bool() bool {
	return len(v.value) > 0
}

func (v ListValue) String() string {
	s, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func (v ListValue) Hash(seed maphash.Seed) (uint64, error) {
	return 0, fmt.Errorf("list is mutable")
}

func (v ListValue) Cmp(other Value) int {
	switch other := other.(type) {
	case *ListValue:
		for i, item := range v.value {
			if i > len(other.value)-1 {
				return 1
			}
			cmp := item.Cmp(other.value[i])
			if cmp != 0 {
				return cmp
			}
		}
		if len(v.value) < len(other.value) {
			return -1
		}
		return 0
	}
	return cmpType(&v, other)
}

func (v ListValue) MarshalYAML() (any, error) {
	return v.value, nil
}

func (v *ListValue) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("cannot unmarshal %s into list", yamlhelpers.KindString(node.Kind))
	}
	if node.ShortTag() != "!!seq" {
		return fmt.Errorf("cannot unmarshal %s into list", node.ShortTag())
	}
	for i, child := range node.Content {
		value, err := UnmarshalYAML(child)
		if err != nil {
			return fmt.Errorf("error unmarshalling item %d: %w", i, err)
		}
		v.value = append(v.value, value)
	}
	return nil
}

type mapPair struct {
	key   Value
	value Value
}

type MapValue struct {
	seed  maphash.Seed
	value map[uint64][]mapPair
}

func (v MapValue) Bool() bool {
	return len(v.value) > 0
}

func (v MapValue) String() string {
	s, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func (v MapValue) Hash(seed maphash.Seed) (uint64, error) {
	return 0, fmt.Errorf("map is mutable")
}

func (v MapValue) Cmp(other Value) int {
	switch other.(type) {
	case *MapValue:
		// TODO: implement map equality check
		return 0
	}
	return cmpType(&v, other)
}

func (v *MapValue) SetItem(key, value Value) error {
	if v.value == nil {
		v.seed = maphash.MakeSeed()
		v.value = make(map[uint64][]mapPair)
	}
	keyHash, err := key.Hash(v.seed)
	if err != nil {
		return err
	}
	for _, existingPair := range v.value[keyHash] {
		if existingPair.key.Cmp(key) == 0 {
			existingPair.value = value
			return nil
		}
	}
	v.value[keyHash] = append(v.value[keyHash], mapPair{key, value})
	return nil
}

func (v MapValue) MarshalYAML() (any, error) {
	mapNode := &yaml.Node{Kind: yaml.MappingNode}
	for _, pairs := range v.value {
		for _, pair := range pairs {
			var key, value yaml.Node
			err := key.Encode(pair.key)
			if err != nil {
				return nil, fmt.Errorf("error marshalling key: %w", err)
			}
			err = value.Encode(pair.value)
			if err != nil {
				return nil, fmt.Errorf("error marshalling value: %w", err)
			}
			mapNode.Content = append(mapNode.Content, &key, &value)
		}
	}
	return mapNode, nil
}

func (v *MapValue) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("cannot unmarshal %s into map", yamlhelpers.KindString(node.Kind))
	}
	if node.ShortTag() != "!!map" {
		return fmt.Errorf("cannot unmarshal %s into map", node.ShortTag())
	}
	for i, value := range node.Content {
		if i&1 == 0 {
			continue // skip keys
		}
		key := node.Content[i-1]
		keyValue, err := UnmarshalYAML(key)
		if err != nil {
			return fmt.Errorf("error unmarshalling key: %w", err)
		}
		valueValue, err := UnmarshalYAML(value)
		if err != nil {
			return fmt.Errorf("error unmarshalling value: %w", err)
		}
		err = v.SetItem(keyValue, valueValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalYAML(node *yaml.Node) (ret Value, err error) {
	switch node.Kind {
	case yaml.ScalarNode:
		switch node.ShortTag() {
		case "!!null":
			ret = NullValue{}
		case "!!bool":
			ret = &BoolValue{}
		case "!!int":
			ret = &IntValue{}
		case "!!float":
			// It's possible this was a misidentified integer due to size.
			testInt := big.Int{}
			_, success := testInt.SetString(node.Value, 0)
			if success {
				ret = &IntValue{}
			} else {
				ret = &FloatValue{}
			}
		case "!!str":
			ret = &StringValue{}
		default:
			return nil, fmt.Errorf("cannot unmarshal %s", node.ShortTag())
		}
	case yaml.SequenceNode:
		ret = &ListValue{}
	case yaml.MappingNode:
		panic("not implemented yet")
	default:
		return nil, fmt.Errorf("cannot unmarshal %s", yamlhelpers.KindString(node.Kind))
	}
	err = ret.UnmarshalYAML(node)
	return
}

func cmpType(v, other Value) int {
	vType := cmpTypeOrdering(v)
	otherType := cmpTypeOrdering(other)
	if vType < otherType {
		return -1
	}
	if vType > otherType {
		return 1
	}
	return 0
}

// cmpTypeOrdering determines type ordering for sorting different types
func cmpTypeOrdering(v Value) string {
	switch v.(type) {
	case NullValue:
		return "000"
	case *BoolValue:
		return "001"
	case *IntValue:
		return "002"
	case *FloatValue:
		return "003"
	case *StringValue:
		return "004"
	case *ListValue:
		return "005"
	case *MapValue:
		return "006"
	}
	return reflect.TypeOf(v).Name()
}
