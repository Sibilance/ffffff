package ast

import (
	"fmt"
	"strings"
)

func assertOneOf[T interface {
	fmt.Stringer
	comparable
}](msg string, actual T, expecteds []T) error {
	for _, expected := range expecteds {
		if actual == expected {
			return nil
		}
	}
	if len(expecteds) == 1 {
		expected := expecteds[0]
		return fmt.Errorf(
			"%s: expected %s, got %s",
			msg,
			expected,
			actual,
		)
	}
	expectedStrs := []string{}
	for _, expected := range expecteds {
		expectedStrs = append(expectedStrs, expected.String())
	}
	return fmt.Errorf(
		"%s: expected one of [%s], got %s",
		msg,
		strings.Join(expectedStrs, ", "),
		actual,
	)
}

func assertKind(node Node, kinds ...Kind) (err error) {
	err = assertOneOf("unexpected kind", node.Kind(), kinds)
	node.ReportError(err)
	return
}

func assertTag(node Node, tags ...Tag) (err error) {
	err = assertOneOf("unexpected tag", node.Tag(), tags)
	node.ReportError(err)
	return
}
