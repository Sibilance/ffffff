package ast

type Kind uint8

const (
	UndefinedNode Kind = iota
	SequenceNode
	MappingNode
	ScalarNode
)

func (kind Kind) String() string {
	switch kind {
	case SequenceNode:
		return "SequenceNode"
	case MappingNode:
		return "MappingNode"
	case ScalarNode:
		return "ScalarNode"
	default:
		return "UndefinedNode"
	}
}
