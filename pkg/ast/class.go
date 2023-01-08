package ast

type Class struct {
	ASTNode

	Name    string
	Fields  map[string]TypeDefinition
	Methods map[string]Function
}
