package ast

type Function struct {
	Node

	Name string

	PositionalArguments []FunctionArgumentDefinition
	KeywordArguments    map[string]FunctionArgumentDefinition
	ReturnType          TypeDefinition

	Body CodeBlock
}

type FunctionArgumentDefinition struct {
	Node

	Name           string
	TypeDefinition TypeDefinition
}
