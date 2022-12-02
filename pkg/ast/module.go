package ast

type Module struct {
	Name string

	Imports   map[string]string
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}
