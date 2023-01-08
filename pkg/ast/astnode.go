package ast

type ASTNode interface {
	FileName() string
	Line() int
	Column() int

	Kind() Kind
	AsSequence() ([]ASTNode, error)
	AsMapping() (map[string]ASTNode, error)
	AsString() (string, error)

	Tag() string
}

type Kind uint8

const (
	UndefinedNode Kind = iota
	SequenceNode
	MappingNode
	StringNode
)
