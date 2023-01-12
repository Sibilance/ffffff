package ast

const ModuleTag Tag = "!module"

type Module[N Node] struct {
	Node N

	Name      string
	Imports   map[string]Import[N]
	Classes   map[string]Class[N]
	Functions map[string]Function[N]
	Constants map[string]Expression[N]
}

/*
ParseModule expects a map from global variable names to imports, classes,
functions, and constants.
*/
func ParseModule[N Node](name string, node N) (mod Module[N], err Error[N]) {
	if node.Tag() != ModuleTag {
		err = NewError(node, "expected %s tag", ModuleTag)
	}
	return
}
