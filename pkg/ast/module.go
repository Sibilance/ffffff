package ast

const ModuleTag Tag = "!module"

type Module struct {
	Node

	Name      string
	Imports   map[string]Import
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}

/*
ParseModule expects a map from global variable names to imports, classes,
functions, and constants.
*/
func ParseModule(name string, node Node) (mod Module, err Error) {
	if node.Tag() != ModuleTag {
		err = NewError(node, "expected %s tag", ModuleTag)
	}
	return
}
