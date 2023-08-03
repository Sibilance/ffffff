package macro

import "gopkg.in/yaml.v3"

const (
	VoidTag   = "!void"
	UnwrapTag = "!unwrap"
)

func IsVoid(node *yaml.Node) bool {
	return node.ShortTag() == VoidTag
}

func IsUnwrap(node *yaml.Node) bool {
	return node.ShortTag() == UnwrapTag
}
