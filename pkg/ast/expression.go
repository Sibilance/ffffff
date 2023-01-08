package ast

type Expression interface{}

type AssignmentExpression struct {
	Node

	Target VariableExpression
	Value  Expression
}

type CallExpression struct {
	Node

	Callable            Expression
	PositionalArguments []Expression
	KeywordArguments    map[string]Expression
}

type LiteralExpression struct {
	Node
	// TODO: What literal types should be supported?
	// How best to support nested literals? Lists, maps, structs, etc.
	// How best to support variable precision types like int?
}

type VariableExpression struct {
	Node

	Name string
}
