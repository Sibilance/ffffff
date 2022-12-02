package ast

type Function struct {
	Name string

	PositionalArguments []TypeDefinition
	KeywordArguments    map[string]TypeDefinition
	ReturnType          TypeDefinition

	Body CodeBlock
}
