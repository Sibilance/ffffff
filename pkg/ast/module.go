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
func ParseModule(name string, node Node) *Module {
	if assertNodeTag(node, ModuleTag) != nil {
		return nil
	}
	return nil
}
