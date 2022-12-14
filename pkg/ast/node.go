package ast

import (
	"errors"
)

type Node interface {
	Tag() Tag
	Kind() Kind
	AsSequence() []Node
	AsMapping() map[string]Node
	AsScalar() string

	LocatorMessage() string
}

type NodeTemplate struct {
	Tag      Tag
	Sequence []NodeTemplate
	Mapping  map[string]NodeTemplate
	Scalar   string

	LocatorMessage string
}

func NewNode(template NodeTemplate) (Node, error) {
	kind := ScalarNode
	if template.Sequence != nil {
		kind = SequenceNode
	}
	if template.Mapping != nil {
		if kind != UndefinedNode {
			return nil, errors.New("cannot define both Sequence and Mapping in node template")
		}
		kind = MappingNode
	}
	if template.Scalar != "" && kind != ScalarNode {
		return nil, errors.New("cannot define more than one of Sequence/Mapping/Scalar in node template")
	}

	node := simpleNode{
		tag:            template.Tag,
		kind:           kind,
		scalar:         template.Scalar,
		locatorMessage: template.LocatorMessage,
	}

	if kind == SequenceNode {
		for _, subTemplate := range template.Sequence {
			subNode, err := NewNode(subTemplate)
			if err != nil {
				return nil, err
			}
			node.sequence = append(node.sequence, subNode)
		}
	}

	if kind == MappingNode {
		node.mapping = map[string]Node{}
		for key, subTemplate := range template.Mapping {
			subNode, err := NewNode(subTemplate)
			if err != nil {
				return nil, err
			}
			node.mapping[key] = subNode
		}
	}

	return node, nil
}

type simpleNode struct {
	tag      Tag
	kind     Kind
	sequence []Node
	mapping  map[string]Node
	scalar   string

	locatorMessage string
}

func (node simpleNode) Tag() Tag {
	return node.tag
}

func (node simpleNode) Kind() Kind {
	return node.kind
}

func (node simpleNode) AsSequence() []Node {
	return node.sequence
}

func (node simpleNode) AsMapping() map[string]Node {
	return node.mapping
}

func (node simpleNode) AsScalar() string {
	return node.scalar
}

func (node simpleNode) LocatorMessage() string {
	return node.locatorMessage
}
