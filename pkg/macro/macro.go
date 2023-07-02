package macro

import "gopkg.in/yaml.v3"

type Context struct {
	// Builtins and globals will go here.
}

func ProcessNode(context Context, node *yaml.Node) error {
	return nil
}
