package macro

import (
	"fmt"
	"strings"

	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

const (
	VoidTag = "!void"
)

type Context struct {
	parent *Context
	label  string

	// Builtins and globals will go here.
}

func (c *Context) path() []string {
	if c.parent == nil {
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
	}
}

func ProcessNode(context *Context, node *yaml.Node) error {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) != 1 {
			return context.Error(node, "expected exactly one child node of document")
		}
		child := node.Content[0]
		err := ProcessNode(context.New(""), child)
		if err != nil {
			return err
		}
		if IsVoid(child) {
			Void(node)
		}

	case yaml.SequenceNode:
		children := node.Content
		node.Content = nil
		for i, child := range children {
			err := ProcessNode(context.New(fmt.Sprint(i)), child)
			if err != nil {
				return err
			}
			if !IsVoid(child) {
				node.Content = append(node.Content, child)
			}
		}

	case yaml.MappingNode:

	case yaml.ScalarNode:

	case yaml.AliasNode:

	default:

	}

	return nil
}

func ProcessDocuments(context *Context, documents *[]*yaml.Node) error {
	originalDocuments := *documents
	*documents = []*yaml.Node{}

	for i, document := range originalDocuments {
		if document.Kind != yaml.DocumentNode {
			return context.Error(
				document,
				fmt.Sprintf("expected DocumentNode, got %s", yamlhelpers.KindString(document.Kind)),
			)
		}
		err := ProcessNode(context.New(fmt.Sprint(i)), document)
		if err != nil {
			return err
		}
		if !IsVoid(document) {
			*documents = append(*documents, document)
		}
	}

	return nil
}

func IsVoid(node *yaml.Node) bool {
	return node.ShortTag() == VoidTag
}

func Void(node *yaml.Node) {
	node.Tag = VoidTag
}
