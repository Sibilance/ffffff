package macro

import (
	"fmt"
	"strings"

	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

const (
	VoidTag   = "!void"
	UnwrapTag = "!unwrap"
)

type Context struct {
	parent *Context
	label  string

	// Builtins and globals will go here.
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
	}
}

func ProcessNode(context *Context, node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		children := node.Content
		node.Content = nil
		for i, child := range children {
			err := ProcessNode(context.New(fmt.Sprintf("[%d]", i)), child)
			if err != nil {
				return err
			}
			if IsUnwrap(child) {
				child.Tag = ""
				*node = *child
				node.Tag = node.ShortTag()
				break
			}
			if !IsVoid(child) {
				node.Content = append(node.Content, child)
			}
		}

	case yaml.MappingNode:
		children := node.Content
		node.Content = nil
		for i, child := range children {
			if i&1 == 0 {
				err := ProcessNode(context.New(fmt.Sprintf("[%d](key)", i/2)), child)
				if err != nil {
					return err
				}
			} else {
				key := children[i-1]
				var contextName string
				if key.Kind == yaml.ScalarNode && key.Tag == "!!str" {
					contextName = key.Value
				} else {
					contextName = fmt.Sprintf("[%d](value)", i/2)
				}
				if !IsVoid(key) {
					err := ProcessNode(context.New(contextName), child)
					if err != nil {
						return err
					}
					if !IsVoid(child) {
						node.Content = append(node.Content, key, child)
					}
				}
			}
		}

	case yaml.ScalarNode:

	case yaml.AliasNode:

	default:
		return context.Error(node, fmt.Sprintf("unexpected node kind, %s", yamlhelpers.KindString(node.Kind)))
	}

	return nil
}

func ProcessDocuments(context *Context, documents *[]*yaml.Node) error {
	originalDocuments := *documents
	*documents = nil

	for i, document := range originalDocuments {
		localContext := context.New(fmt.Sprintf("[%d]", i))
		if document.Kind != yaml.DocumentNode {
			return localContext.Error(
				document,
				fmt.Sprintf("expected DocumentNode, got %s", yamlhelpers.KindString(document.Kind)),
			)
		}
		if len(document.Content) != 1 {
			return localContext.Error(document, "expected exactly one child node of document")
		}
		child := document.Content[0]
		err := ProcessNode(localContext, child)
		if err != nil {
			return err
		}
		if !IsVoid(child) {
			*documents = append(*documents, document)
		}
		if IsUnwrap(child) {
			switch child.Kind {
			case yaml.SequenceNode:
				*documents = nil
				for _, newChild := range child.Content {
					*documents = append(*documents, &yaml.Node{
						Kind: yaml.DocumentNode,
						Content: []*yaml.Node{
							newChild,
						},
						Line:   newChild.Line,
						Column: newChild.Column,
					})
				}
				return nil
			case yaml.MappingNode:
				*documents = nil
				for i, newValue := range child.Content {
					if i&1 == 0 {
						continue
					}
					newKey := child.Content[i-1]
					*documents = append(*documents, &yaml.Node{
						Kind: yaml.DocumentNode,
						Content: []*yaml.Node{
							{
								Kind: yaml.MappingNode,
								Tag:  "!!map",
								Content: []*yaml.Node{
									newKey,
									newValue,
								},
								Line:   newKey.Line,
								Column: newKey.Column,
							},
						},
						Line:   newKey.Line,
						Column: newKey.Column,
					})
					fmt.Printf("%+v\n", (*documents)[len(*documents)-1])
				}
			default:
				return localContext.Error(child, fmt.Sprintf("cannot unwrap %s", yamlhelpers.KindString(child.Kind)))
			}
		}
	}

	return nil
}

func IsVoid(node *yaml.Node) bool {
	return node.ShortTag() == VoidTag
}

func IsUnwrap(node *yaml.Node) bool {
	return node.ShortTag() == UnwrapTag
}
