package ast

import (
	"fmt"
	"strings"
)

type Error interface {
	Node

	Error() string
	InnerErrors() []Error
	appendError(Error)
	orNil(string) Error
}

func newError(node Node, message string, args ...fmt.Stringer) Error {
	return &error_{
		Node:        node,
		message:     fmt.Sprintf(message, args),
		innerErrors: nil,
	}
}

type error_ struct {
	Node
	message     string
	innerErrors []Error
}

func (err *error_) Error() string {
	return fmt.Sprintf("%s: %s", err.Node.LocatorMessage(), err.message)
}

func (err *error_) InnerErrors() []Error {
	return err.innerErrors
}

func (err *error_) appendError(innerError Error) {
	err.innerErrors = append(err.innerErrors, innerError)
}

func (err *error_) orNil(message string) Error {
	if err.message != "" || len(err.innerErrors) > 0 {
		if err.message == "" {
			err.message = message
		}
		return err
	}
	return nil
}

func assertEqual[T interface {
	fmt.Stringer
	comparable
}](node Node, actual T, expecteds ...T) Error {
	for _, expected := range expecteds {
		if actual == expected {
			return nil
		}
	}
	if len(expecteds) == 1 {
		expected := expecteds[0]
		return newError(
			node,
			fmt.Sprintf(
				"expected %s, got %s",
				expected,
				actual,
			),
		)
	}
	expectedStrs := []string{}
	for _, expected := range expecteds {
		expectedStrs = append(expectedStrs, expected.String())
	}
	return newError(
		node,
		fmt.Sprintf(
			"expected one of [%s], got %s",
			strings.Join(expectedStrs, ", "),
			actual,
		),
	)
}

func assertNodeKind(node Node, kinds ...Kind) Error {
	return assertEqual(node, node.Kind(), kinds...)
}

func assertNodeTag(node Node, tags ...Tag) Error {
	return assertEqual(node, node.Tag(), tags...)
}
