package ast

type Module struct {
	AstNode

	Name      string
	Imports   map[string]Import
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}
