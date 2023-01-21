package ast

import (
	"errors"
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
func ParseImport(node Node) *Import {
	if assertNodeTag(node, ModuleTag) != nil {
		return nil
	}
	if assertNodeKind(node, ScalarNode, SequenceNode) != nil {
		return nil
	}

	var path []string
	var failed bool

	switch node.Kind() {
	case ScalarNode:
		path = strings.Split(node.AsScalar(), ".")
	case SequenceNode:
		for _, innerNode := range node.AsSequence() {
			if assertNodeKind(innerNode, ScalarNode) != nil {
				failed = true
				continue
			}
			path = append(path, innerNode.AsScalar())
		}
	}

	if failed {
		node.ReportError(errors.New("invalid Import definition"))
	}

	return &Import{
		Node: node,
		Path: path,
	}
}
