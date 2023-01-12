package ast

type Class[N Node] struct {
	Node N

	Name    string
	Fields  map[string]TypeDefinition[N]
	Methods map[string]Function[N]
}
