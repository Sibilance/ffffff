package yamlreader

import (
	"os"
	"path/filepath"
	"testing"
)

func getTestPath(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(cwd, "test", t.Name())
}

func getTestFile(t *testing.T) string {
	return getTestPath(t) + ".yaml"
}

func compareNodes(t *testing.T, actual, expected *Node) {
	if actual.Name != expected.Name {
		t.Fatalf("%s: expected Name=%s, got %s", actual.Name, expected.Name, actual.Name)
	}

	if actual.FileName != expected.FileName {
		t.Fatalf("%s: expected FileName=%s, got %s", actual.Name, expected.FileName, actual.FileName)
	}

	if actual.Line != expected.Line {
		t.Fatalf("%s: expected Line=%d, got %d", actual.Name, expected.Line, actual.Line)
	}

	if actual.Column != expected.Column {
		t.Fatalf("%s: expected Column=%d, got %d", actual.Name, expected.Column, actual.Column)
	}

	if actual.Comment != expected.Comment {
		t.Fatalf("%s: expected Comment:\n%s\n\ngot:\n%s", actual.Name, expected.Comment, actual.Comment)
	}

	if actual.Kind != expected.Kind {
		t.Fatalf("%s: expected Kind=%s, got %s", actual.Name, expected.Kind, actual.Kind)
	}

	if actual.Tag != expected.Tag {
		t.Fatalf("%s: expected Tag=%s, got %s", actual.Name, expected.Tag, actual.Tag)
	}

	if len(actual.Sequence) != len(expected.Sequence) {
		t.Fatalf("%s: expected len(Sequence)=%d, got %d", actual.Name, len(expected.Sequence), len(actual.Sequence))
	}

	for i, innerExpected := range expected.Sequence {
		innerActual := actual.Sequence[i]
		compareNodes(t, &innerActual, &innerExpected)
	}

	if len(actual.Mapping) != len(expected.Mapping) {
		t.Fatalf("%s: expected len(Mapping)=%d, got %d", actual.Name, len(expected.Mapping), len(actual.Mapping))
	}

	for k, innerExpected := range expected.Mapping {
		innerActual := actual.Mapping[k]
		compareNodes(t, &innerActual, &innerExpected)
	}

	if (actual.Scalar == nil) != (expected.Scalar == nil) {
		t.Fatalf("%s: expected (Value == nil)=%t, got (Value == nil)=%t", actual.Name, expected.Scalar == nil, actual.Scalar == nil)
	}

	if actual.Scalar != nil {
		if actual.Bool() != expected.Bool() {
			t.Fatalf("%s: expected Bool=%t, got %t", actual.Name, expected.Bool(), actual.Bool())
		}

		if actual.Int() != expected.Int() {
			t.Fatalf("%s: expected Int=%d, got %d", actual.Name, expected.Int(), actual.Int())
		}

		if actual.Float() != expected.Float() {
			t.Fatalf("%s: expected Float=%f, got %f", actual.Name, expected.Float(), actual.Float())
		}

		if actual.Str() != expected.Str() {
			t.Fatalf("%s: expected Str=%s, got %s", actual.Name, expected.Str(), actual.Str())
		}
	}
}

func readTestDirectory(t *testing.T) *Node {
	node := Node{}
	err := node.ReadDirectory(getTestPath(t))
	if err != nil {
		t.Fatal(err)
	}
	return &node
}

func readTestFile(t *testing.T) *Node {
	node := Node{}
	err := node.ReadFile(getTestFile(t))
	if err != nil {
		t.Fatal(err)
	}
	return &node
}
