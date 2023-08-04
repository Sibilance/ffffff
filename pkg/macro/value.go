package macro

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Value struct {
	yaml.Node
}

func (v Value) Bool() bool {
	switch v.Kind {
	case yaml.DocumentNode:
		return Value{*v.Content[0]}.Bool()
	case yaml.SequenceNode:
		return len(v.Content) > 0
	case yaml.MappingNode:
		return len(v.Content) > 0
	case yaml.ScalarNode:
		switch v.Node.ShortTag() {
		case "!!null":
			return false
		case "!!bool":
			var ret bool
			err := v.Node.Decode(&ret)
			if err != nil {
				return true
			}
			return ret
		case "!!int":
			var ret int64
			err := v.Node.Decode(&ret)
			if err != nil {
				return true
			}
			return ret != 0
		case "!!float":
			var ret float64
			err := v.Node.Decode(&ret)
			if err != nil {
				return true
			}
			return ret != 0
		case "!!str":
			return len(v.Value) > 0
		default:
			return true
		}
	case yaml.AliasNode:
		return Value{*v.Alias}.Bool()
	default:
		panic(fmt.Errorf("undefined Kind %d", v.Kind))
	}
}
