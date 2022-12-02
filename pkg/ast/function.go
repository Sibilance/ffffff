package ast

type Function struct {
	AstNode

	Name string

	PositionalArguments map[string]TypeDefinition
	KeywordArguments    map[string]TypeDefinition
	ReturnType          TypeDefinition

	Body CodeBlock
}
