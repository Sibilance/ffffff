package languagemodels

type Class struct {
	Name string

	Methods map[string]Function
	Fields  map[string]TypeDefinition
}
