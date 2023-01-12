package ast

import "fmt"

type Error[N Node] *simpleError[N]

func NewError[N Node](node N, message string, args ...fmt.Stringer) Error[N] {
	return &simpleError[N]{
		Node:        node,
		message:     fmt.Sprintf(message, args),
		innerErrors: nil,
	}
}

func NewNestedError[N Node](node N, innerErrors []Error[N], message string, args ...fmt.Stringer) Error[N] {
	return &simpleError[N]{
		Node:        node,
		message:     fmt.Sprintf(message, args),
		innerErrors: innerErrors,
	}
}

type simpleError[N Node] struct {
	Node        N
	message     string
	innerErrors []Error[N]
}

func (err simpleError[N]) Error() string {
	return err.message
}

func (err simpleError[N]) InnerErrors() []Error[N] {
	return err.innerErrors
}

func assertNodeKindIs[N Node](node N, kind Kind) Error[N] {
	if node.Kind() != kind {
		return NewError(
			node,
			fmt.Sprintf(
				"expected %s, got %s",
				kind,
				node.Kind(),
			),
		)
	}
	return nil
}

func assertNodeTagIs[N Node](node N, tag Tag) Error[N] {
	if node.Tag() != tag {
		return NewError(
			node,
			fmt.Sprintf(
				"expected %s, got %s",
				tag,
				node.Tag(),
			),
		)
	}
	return nil
}
