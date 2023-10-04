package macro

import (
	"fmt"
	"strings"

	"github.com/sibilance/ffffff/pkg/value"
	"gopkg.in/yaml.v3"
)

type Context struct {
	parent *Context
	label  string

	macros map[string]Macro
	locals map[string]value.Value
}

func DefaultRootContext() *Context {
	return &Context{
		macros: DefaultMacros(),
	}
}

func (c *Context) path() []string {
	if c.parent == nil {
		if c.label == "" {
			return nil
		}
		return []string{c.label}
	}
	if c.label != "" {
		return append(c.parent.path(), c.label)
	}
	return c.parent.path()
}

func (c *Context) Path() string {
	return strings.Join(c.path(), ".")
}

func (c *Context) Error(node *yaml.Node, msg string) error {
	return fmt.Errorf("%s:%d:%d: %s", c.Path(), node.Line, node.Column, msg)
}

func (c *Context) New(label string) *Context {
	return &Context{
		parent: c,
		label:  label,
		macros: c.macros,
		locals: c.locals,
	}
}
