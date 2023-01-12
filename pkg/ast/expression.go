package ast

type Expression[N Node] interface {
	Node
}

type AssignmentExpression[N Node] struct {
	Node N

	Target VariableExpression[N]
	Value  Expression[N]
}

type CallExpression[N Node] struct {
	Node N

	Callable            Expression[N]
	PositionalArguments []Expression[N]
	KeywordArguments    map[string]Expression[N]
}

type LiteralExpression[N Node] struct {
	Node N
	// TODO: What literal types should be supported?
	// How best to support nested literals? Lists, maps, structs, etc.
	// How best to support variable precision types like int?
}

type VariableExpression[N Node] struct {
	Node N

	Name string
}
