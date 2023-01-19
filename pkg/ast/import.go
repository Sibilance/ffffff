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
	if err := assertNodeTag(node, ModuleTag); err != nil {
		return nil, err
	}
	if err := assertNodeKind(node, ScalarNode, SequenceNode); err != nil {
		return nil, err
	}

	var path []string
	err := newError(node, "")

	switch node.Kind() {
	case ScalarNode:
		path = strings.Split(node.AsScalar(), ".")
	case SequenceNode:
		for _, innerNode := range node.AsSequence() {
			if innerError := assertNodeKind(innerNode, ScalarNode); innerError != nil {
				err.appendError(innerError)
			}
			path = append(path, innerNode.AsScalar())
		}
	}

	return &Import{
		Node: node,
		Path: path,
	}, err.orNil("invalid Import definition")
}
