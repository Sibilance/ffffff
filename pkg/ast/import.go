package ast

type Import struct {
	AstNode

	Name string
	Path Path
}
