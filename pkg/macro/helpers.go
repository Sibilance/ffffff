package macro

import "gopkg.in/yaml.v3"

func IsVoid(node *yaml.Node) bool {
	return node.ShortTag() == VoidTag
}

func IsUnwrap(node *yaml.Node) bool {
	return node.ShortTag() == UnwrapTag
}
