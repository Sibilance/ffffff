package yamlreader

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Kind int8

const (
	UndefinedNode Kind = iota
	SequenceNode
	MappingNode
	ScalarNode
)

func (kind Kind) String() string {
	switch kind {
	case UndefinedNode:
		return "UndefinedNode"
	case SequenceNode:
		return "SequenceNode"
	case MappingNode:
		return "MappingNode"
	case ScalarNode:
		return "ScalarNode"
	}
	return "unknown"
}

type Node struct {
	Name     string
	FileName string
	Line     int
	Column   int
	Comment  string
	Kind     Kind
	Tag      string
	Sequence []Node
	Mapping  map[string]Node
	Scalar   string
}

func (n *Node) String() string {
	return fmt.Sprintf("%s (%s:%d:%d)", n.Name, n.FileName, n.Line, n.Column)
}

func (n *Node) UnmarshalYAML(yamlNode *yaml.Node) error {
	recursionDetection := make(map[*yaml.Node]struct{})
	return n.unmarshalYAML(yamlNode, recursionDetection)
}

func (n *Node) unmarshalYAML(yamlNode *yaml.Node, recursionDetection map[*yaml.Node]struct{}) error {
	n.Tag = yamlNode.ShortTag()
	n.Line = yamlNode.Line
	n.Column = yamlNode.Column
	n.Comment = strings.Trim(yamlNode.HeadComment+"\n\n"+yamlNode.LineComment, "\n")

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
		return n.unmarshalYAML(yamlNode.Alias, recursionDetection)

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
		mapping := map[string]*yaml.Node{}
		if err := yamlNode.Decode(mapping); err != nil {
			return fmt.Errorf("%s: %s", n, err)
		}
		for k, innerYamlNode := range mapping {
			innerNode := Node{
				Name:     fmt.Sprintf("%s.%s", n.Name, k),
				FileName: n.FileName,
			}
			if err := innerNode.unmarshalYAML(innerYamlNode, recursionDetection); err != nil {
				return err
			}
			n.Mapping[k] = innerNode
		}

	case yaml.ScalarNode:
		n.Kind = ScalarNode
		n.Scalar = yamlNode.Value

	default:
		return fmt.Errorf("%s: unknown node kind", n)
	}

	return nil
}
