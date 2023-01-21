package ast

import (
	"fmt"
	"strings"
)

func assertOneOf[T interface {
	fmt.Stringer
	comparable
}](actual T, expecteds ...T) error {
	for _, expected := range expecteds {
		if actual == expected {
			return nil
		}
	}
	if len(expecteds) == 1 {
		expected := expecteds[0]
		return fmt.Errorf(
			"expected %s, got %s",
			expected,
			actual,
		)
	}
	expectedStrs := []string{}
	for _, expected := range expecteds {
		expectedStrs = append(expectedStrs, expected.String())
	}
	return fmt.Errorf(
		"expected one of [%s], got %s",
		strings.Join(expectedStrs, ", "),
		actual,
	)
}

func assertNodeKind(node Node, kinds ...Kind) (err error) {
	err = assertOneOf(node.Kind(), kinds...)
	node.ReportError(err)
	return
}

func assertNodeTag(node Node, tags ...Tag) (err error) {
	err = assertOneOf(node.Tag(), tags...)
	node.ReportError(err)
	return
}
