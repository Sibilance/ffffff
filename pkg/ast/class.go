package ast

type Class struct {
	AstNode

	Name    string
	Fields  map[string]TypeDefinition
	Methods map[string]Function
}
