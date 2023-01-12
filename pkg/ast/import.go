package ast

import (
	"strings"
)

const ImportTag = "!import"

type Import[N Node] struct {
	Node N

	Path []string
}

/*
ParseImport expects either a string of "."-delimited path components,
or a sequence of path components (strings).
*/
func ParseImport[N Node](node N) (*Import[N], Error[N]) {
	if err := assertNodeTagIs(node, ModuleTag); err != nil {
		return nil, err
	}

	var path []string
	var err Error[N]
	var innerErrors []Error[N]

	switch node.Kind() {
	case ScalarNode:
		path = strings.Split(node.AsScalar(), ".")
	case SequenceNode:
		for _, innerNode := range node.AsSequence() {
			innerNode := innerNode.(N)
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

	return &Import[N]{
		Node: node,
		Path: path,
	}, err
}
