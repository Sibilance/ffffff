package ast

const FunctionTag Tag = "!function"

type Function struct {
	Node

	PositionalArguments []FunctionArgumentDefinition
	KeywordArguments    map[string]FunctionArgumentDefinition
	ReturnType          TypeDefinition

	Body CodeBlock
}

func (f *Function) Parse(n Node) {}

type FunctionArgumentDefinition struct {
	Node

	Name           string
	TypeDefinition TypeDefinition
}
