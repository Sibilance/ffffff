package ast

import "fmt"

type Error interface {
	Node

	Error() string
	InnerErrors() []Error
}

func NewError(node Node, message string, args ...fmt.Stringer) Error {
	return simpleError{
		Node:        node,
		message:     fmt.Sprintf(message, args),
		innerErrors: nil,
	}
}

func NewNestedError(node Node, innerErrors []Error, message string, args ...fmt.Stringer) Error {
	return simpleError{
		Node:        node,
		message:     fmt.Sprintf(message, args),
		innerErrors: innerErrors,
	}
}

type simpleError struct {
	Node
	message     string
	innerErrors []Error
}

func (err simpleError) Error() string {
	return err.message
}

func (err simpleError) InnerErrors() []Error {
	return err.innerErrors
}

func assertNodeKindIs(node Node, kind Kind) Error {
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

func assertNodeTagIs(node Node, tag Tag) Error {
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
