package ast

const ModuleTag Tag = "!module"

type Module struct {
	Node

	Imports   map[string]Import
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}

/*
Module.Parse expects a map from global variable names to imports, classes,
functions, and constants.
*/
func (m *Module) Parse(n Node) {
	m.Node = n
	assertTag(m, ModuleTag)

	for _, innerNode := range m.AsMapping() {
		assertTag(innerNode, ImportTag, ClassTag, FunctionTag)
	}

	parseByTag(m, m.Imports, ImportTag)
	parseByTag(m, m.Classes, ClassTag)
	parseByTag(m, m.Functions, FunctionTag)
	parseByTag(m, m.Constants, ConstantTag)

}
