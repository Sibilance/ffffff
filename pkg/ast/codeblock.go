package ast

type CodeBlock[N Node] struct {
	Node N

	Statements []Statement[N]
}
