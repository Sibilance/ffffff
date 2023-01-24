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
Import.Parse expects either a string of "."-delimited path components,
or a sequence of path components (strings).
*/
func (i *Import) Parse() *Import {
	assertTag(i, ModuleTag)
	assertKind(i, ScalarNode, SequenceNode)

	i.Path = nil
	var err error

	switch i.Kind() {
	case ScalarNode:
		i.Path = strings.Split(i.AsScalar(), ".")
	case SequenceNode:
		for _, innerNode := range i.AsSequence() {
			if assertKind(innerNode, ScalarNode) != nil && err == nil {
				err = errors.New("expected sequence of string path components")
			}
			i.Path = append(i.Path, innerNode.AsScalar())
		}
	}

	if err != nil {
		i.ReportError(err)
	}

	return i
}
