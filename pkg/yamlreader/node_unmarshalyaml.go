package yamlreader

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func (n *Node) UnmarshalYAML(yamlNode *yaml.Node) error {
	recursionDetection := make(map[*yaml.Node]struct{})
	return n.unmarshalYAML(yamlNode, recursionDetection)
}

func (n *Node) unmarshalYAML(yamlNode *yaml.Node, recursionDetection map[*yaml.Node]struct{}) error {
	n.Tag = yamlNode.ShortTag()
	n.Line = yamlNode.Line
	n.Column = yamlNode.Column
	n.Comment = strings.Trim(yamlNode.HeadComment+"\n"+yamlNode.LineComment, "\n")

	switch yamlNode.Kind {
	case yaml.DocumentNode:
		if len(yamlNode.Content) != 1 {
			return fmt.Errorf("%s: expected one node in document, got %d", n, len(yamlNode.Content))
		}
		return n.unmarshalYAML(yamlNode.Content[0], recursionDetection)

	case yaml.AliasNode:
		// If we've encountered this particular AliasNode before, it must be a recursive alias.
		// Note that we track the AliasNode itself, not the node it is aliased to, as multiple
		// pointers to the same node are valid and don't indicate recursion.
		if _, contained := recursionDetection[yamlNode]; contained {
			return fmt.Errorf("%s: recursive alias detected", n)
		}
		recursionDetection[yamlNode] = struct{}{}
		// Follow the alias, but keep these attributes from the alias node.
		line, column, comment := n.Line, n.Column, n.Comment
		err := n.unmarshalYAML(yamlNode.Alias, recursionDetection)
		n.Line, n.Column, n.Comment = line, column, comment
		return err

	case yaml.SequenceNode:
		n.Kind = SequenceNode
		for i, innerYamlNode := range yamlNode.Content {
			innerNode := Node{
				Name:     fmt.Sprintf("%s[%d]", n.Name, i),
				FileName: n.FileName,
			}
			if err := innerNode.unmarshalYAML(innerYamlNode, recursionDetection); err != nil {
				return err
			}
			n.Sequence = append(n.Sequence, innerNode)
		}

	case yaml.MappingNode:
		n.Kind = MappingNode
		n.Mapping = make(map[string]Node)
		mapping := map[string]yaml.Node{}
		if err := yamlNode.Decode(mapping); err != nil {
			return fmt.Errorf("%s: %s", n, err)
		}
		for k, innerYamlNode := range mapping {
			innerNode := Node{
				Name:     fmt.Sprintf("%s.%s", n.Name, k),
				FileName: n.FileName,
			}
			if err := innerNode.unmarshalYAML(&innerYamlNode, recursionDetection); err != nil {
				return err
			}
			n.Mapping[k] = innerNode
		}

	case yaml.ScalarNode:
		n.Raw = yamlNode.Value
		switch n.Tag {
		case BoolTag:
			n.Kind = BooleanNode
			value := BoolScalar{}
			yamlNode.Decode(&value.Value)
			n.Scalar = value
		case IntTag:
			n.Kind = IntegerNode
			// We have to check for overflows of positive integers because the yaml
			// library detects anything that fits in int64 or uint64 as !!int.
			var overflowCheck uint64
			yamlNode.Decode(&overflowCheck)
			if overflowCheck > ^uint64(0)>>1 {
				n.Kind = FloatNode
				n.Tag = FloatTag
				value := FloatScalar{}
				yamlNode.Decode(&value.Value)
				n.Scalar = value
			} else {
				value := IntScalar{}
				yamlNode.Decode(&value.Value)
				n.Scalar = value
			}
		case FloatTag:
			n.Kind = FloatNode
			value := FloatScalar{}
			yamlNode.Decode(&value.Value)
			n.Scalar = value
		case NullTag:
			n.Kind = NullNode
		default:
			n.Kind = StringNode
			n.Scalar = StringScalar{Value: n.Raw}
		}

	default:
		return fmt.Errorf("%s: unknown node kind", n)
	}

	return nil
}
