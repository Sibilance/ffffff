package ast

const ClassTag Tag = "!class"

type Class struct {
	Node

	Fields  map[string]TypeDefinition
	Methods map[string]Function
}

func (c *Class) Parse(n Node) {}
