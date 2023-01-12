package ast

type TypeDefinition[N Node] struct {
	Node N

	Type Expression[N]
}
