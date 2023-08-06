package yamlhelpers

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	VoidKind           yaml.Kind = 0
	UnwrapKind         yaml.Kind = 1 << 5
	UnwrapDocumentNode yaml.Kind = UnwrapKind | yaml.DocumentNode
	UnwrapSequenceNode yaml.Kind = UnwrapKind | yaml.SequenceNode
	UnwrapMappingNode  yaml.Kind = UnwrapKind | yaml.MappingNode
	UnwrapScalarNode   yaml.Kind = UnwrapKind | yaml.ScalarNode
	UnwrapAliasNode    yaml.Kind = UnwrapKind | yaml.AliasNode
)

func IsVoid(n *yaml.Node) bool {
	return n.Kind == VoidKind
}

func IsUnwrap(n *yaml.Node) bool {
	return n.Kind&UnwrapKind > 0
}

func KindString(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.ScalarNode:
		return "ScalarNode"
	case yaml.AliasNode:
		return "AliasNode"
	case UnwrapDocumentNode:
		return "UnwrapDocumentNode"
	case UnwrapSequenceNode:
		return "UnwrapSequenceNode"
	case UnwrapMappingNode:
		return "UnwrapMappingNode"
	case UnwrapScalarNode:
		return "UnwrapScalarNode"
	case UnwrapAliasNode:
		return "UnwrapAliasNode"
	case VoidKind:
		return "VoidKind"
	case UnwrapKind:
		return "UnwrapKind"
	default:
		return fmt.Sprint(kind)
	}
}
