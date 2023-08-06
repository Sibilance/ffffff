package macro

import (
	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

type Macro func(*Context, *yaml.Node) error

func DefaultMacros() (macros map[string]Macro) {
	macros = make(map[string]Macro, 2)
	macros["!void"] = Void
	macros["!unwrap"] = Unwrap
	return
}

func Void(context *Context, node *yaml.Node) error {
	node.Kind = yamlhelpers.VoidKind
	return nil
}

func Unwrap(context *Context, node *yaml.Node) error {
	node.Kind |= yamlhelpers.UnwrapKind
	return nil
}
