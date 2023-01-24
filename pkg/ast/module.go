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
Module.Parse expects a map from global variable names to imports, classes,
functions, and constants.
*/
func (m *Module) Parse(name string) *Module {
	assertTag(m, ModuleTag)
	return m
}
