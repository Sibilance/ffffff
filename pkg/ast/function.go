package ast

type Function struct {
	ASTNode

	Name string

	PositionalArguments []FunctionArgumentDefinition
	KeywordArguments    map[string]FunctionArgumentDefinition
	ReturnType          TypeDefinition

	Body CodeBlock
}

type FunctionArgumentDefinition struct {
	ASTNode

	Name           string
	TypeDefinition TypeDefinition
}
