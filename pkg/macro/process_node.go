package macro

import (
	"fmt"

	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

func processSequenceNode(context *Context, node *yaml.Node) error {
	children := node.Content
	node.Content = nil
	for i, child := range children {
		localContext := context.New(fmt.Sprintf("[%d]", i))
		err := ProcessNode(localContext, child)
		if err != nil {
			return err
		}
		if IsUnwrap(child) {
			switch child.Kind {
			case yaml.SequenceNode:
				node.Content = append(node.Content, child.Content...)
			case yaml.MappingNode:
				for i, newValue := range child.Content {
					if i&1 == 0 {
						continue
					}
					newKey := child.Content[i-1]
					node.Content = append(node.Content, &yaml.Node{
						Kind: yaml.MappingNode,
						Tag:  "!!map",
						Content: []*yaml.Node{
							newKey,
							newValue,
						},
						Line:   newKey.Line,
						Column: newKey.Column,
					})
				}
			default:
				return localContext.Error(child, fmt.Sprintf("cannot unwrap %s", yamlhelpers.KindString(child.Kind)))
			}
		} else if !IsVoid(child) {
			node.Content = append(node.Content, child)
		}
	}

	return nil
}

func processMappingNode(context *Context, node *yaml.Node) error {
	children := node.Content
	node.Content = nil
	for i, child := range children {
		keyContext := context.New(fmt.Sprintf("[%d](key)", i/2))
		if i&1 == 0 {
			err := ProcessNode(keyContext, child)
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
			if IsUnwrap(key) {
				return keyContext.Error(key, "cannot unwrap a mapping key")
			} else if !IsVoid(key) {
				valueContext := context.New(contextName)
				err := ProcessNode(valueContext, child)
				if err != nil {
					return err
				}
				if IsUnwrap(child) {
					return valueContext.Error(child, "cannot unwrap a mapping value")
				} else if !IsVoid(child) {
					node.Content = append(node.Content, key, child)
				}
			}
		}
	}

	return nil
}

func ProcessNode(context *Context, node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		return processSequenceNode(context, node)

	case yaml.MappingNode:
		return processMappingNode(context, node)

	case yaml.ScalarNode:

	case yaml.AliasNode:

	default:
		return context.Error(node, fmt.Sprintf("unexpected node kind, %s", yamlhelpers.KindString(node.Kind)))
	}

	return nil
}
