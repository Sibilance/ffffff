package ast

import "errors"

type Node interface {
	FileName() string
	Line() int
	Column() int

	Tag() string
	Kind() Kind
	AsSequence() ([]Node, Error)
	AsMapping() (map[string]Node, Error)
	AsScalar() (string, Error)
}

type NodeTemplate struct {
	FileName string
	Line     int
	Column   int

	Tag      string
	Sequence []NodeTemplate
	Mapping  map[string]NodeTemplate
	Scalar   string
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
		fileName: template.FileName,
		line:     template.Line,
		column:   template.Column,
		tag:      template.Tag,
		kind:     kind,
		scalar:   template.Scalar,
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

type Kind uint8

const (
	UndefinedNode Kind = iota
	SequenceNode
	MappingNode
	ScalarNode
)

type simpleNode struct {
	fileName string
	line     int
	column   int

	tag      string
	kind     Kind
	sequence []Node
	mapping  map[string]Node
	scalar   string
}

func (node simpleNode) FileName() string {
	return node.fileName
}

func (node simpleNode) Line() int {
	return node.line
}

func (node simpleNode) Column() int {
	return node.column
}

func (node simpleNode) Tag() string {
	return node.tag
}

func (node simpleNode) Kind() Kind {
	return node.kind
}

func (node simpleNode) AsSequence() ([]Node, Error) {
	if node.kind == SequenceNode {
		return node.sequence, nil
	}
	return nil, NewError(node, "not a sequence", nil)
}

func (node simpleNode) AsMapping() (map[string]Node, Error) {
	if node.kind == MappingNode {
		return node.mapping, nil
	}
	return nil, NewError(node, "not a mapping", nil)
}

func (node simpleNode) AsScalar() (string, Error) {
	if node.kind == ScalarNode {
		return node.scalar, nil
	}
	return "", NewError(node, "not a scalar", nil)
}
