package ast

import (
	"strings"
)

const ImportTag = "!import"

type Import struct {
	Node

	Path []string
}

/*
ParseImport expects either a string of "."-delimited path components,
or a sequence of path components (strings).
*/
func ParseImport(node Node) (*Import, Error) {
	if err := assertNodeTagIs(node, ModuleTag); err != nil {
		return nil, err
	}

	var path []string
	var err Error
	var innerErrors []Error

	switch node.Kind() {
	case ScalarNode:
		path = strings.Split(node.AsScalar(), ".")
	case SequenceNode:
		for _, innerNode := range node.AsSequence() {
			if err := assertNodeKindIs(innerNode, ScalarNode); err != nil {
				innerErrors = append(innerErrors, err)
			}
			path = append(path, innerNode.AsScalar())
		}
	default:
		return nil, NewError(node, "expected %s or %s", ScalarNode, SequenceNode)
	}

	if len(innerErrors) != 0 {
		err = NewNestedError(node, innerErrors, "invalid Import definition")
	}

	return &Import{
		Node: node,
		Path: path,
	}, err
}
