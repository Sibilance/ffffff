package ast

type Function[N Node] struct {
	Node N

	Name string

	PositionalArguments []FunctionArgumentDefinition[N]
	KeywordArguments    map[string]FunctionArgumentDefinition[N]
	ReturnType          TypeDefinition[N]

	Body CodeBlock[N]
}

type FunctionArgumentDefinition[N Node] struct {
	Node

	Name           string
	TypeDefinition TypeDefinition[N]
}
