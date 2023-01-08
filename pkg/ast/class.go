package ast

type Class struct {
	Node

	Name    string
	Fields  map[string]TypeDefinition
	Methods map[string]Function
}
