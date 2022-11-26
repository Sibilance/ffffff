package languagemodels

type Module struct {
	Name string

	Imports   map[string]*Module
	Classes   map[string]Class
	Functions map[string]Function
	Constants map[string]Expression
}
