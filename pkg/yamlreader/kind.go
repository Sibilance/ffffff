package yamlreader

type Kind uint8

const (
	UndefinedNode Kind = iota
	SequenceNode
	MappingNode
	BooleanNode
	IntegerNode
	FloatNode
	StringNode
	NullNode
)

func (kind Kind) String() string {
	switch kind {
	case UndefinedNode:
		return "UndefinedNode"
	case SequenceNode:
		return "SequenceNode"
	case MappingNode:
		return "MappingNode"
	case BooleanNode:
		return "BooleanNode"
	case IntegerNode:
		return "IntegerNode"
	case FloatNode:
		return "FloatNode"
	case StringNode:
		return "StringNode"
	case NullNode:
		return "NullNode"
	}
	return "unknown"
}
