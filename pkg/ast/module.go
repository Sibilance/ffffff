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
func ParseModule(name string, node Node) (*Module, Error) {
	if err := assertNodeTag(node, ModuleTag); err != nil {
		return nil, err
	}
	return nil, nil
}
