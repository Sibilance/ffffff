package ast

type Module struct {
	Node

	Name      string
	Imports   map[string]Import
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}

func ParseModule(astNode Node) (mod Module, err Error) {
	return
}
