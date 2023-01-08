package ast

type Module struct {
	ASTNode

	Name      string
	Imports   map[string]Import
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}

func ParseModule(astNode ASTNode) (mod Module, err ASTError) {
	return
}
