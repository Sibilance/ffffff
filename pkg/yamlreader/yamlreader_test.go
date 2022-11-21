package yamlreader

import (
	"os"
	"path/filepath"
	"testing"
)

func getTestFile(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(cwd, "test", t.Name()+".yaml")
}

func compareNodes(t *testing.T, actual, expected *Node) {
	if actual.Name != expected.Name {
		t.Fatalf("expected Name=%s, got %s", expected.Name, actual.Name)
	}

	if actual.FileName != expected.FileName {
		t.Fatalf("expected FileName=%s, got %s", expected.FileName, actual.FileName)
	}

	if actual.Line != expected.Line {
		t.Fatalf("expected Line=%d, got %d", expected.Line, actual.Line)
	}

	if actual.Column != expected.Column {
		t.Fatalf("expected Column=%d, got %d", expected.Column, actual.Column)
	}

	if actual.Comment != expected.Comment {
		t.Fatalf("expected Comment:\n%s\n\ngot:\n%s", expected.Comment, actual.Comment)
	}

	if actual.Kind != expected.Kind {
		t.Fatalf("expected Kind=%s, got %s", expected.Kind, actual.Kind)
	}

	if actual.Tag != expected.Tag {
		t.Fatalf("expected Tag=%s, got %s", expected.Tag, actual.Tag)
	}

	if len(actual.Sequence) != len(expected.Sequence) {
		t.Fatalf("expected len(Sequence)=%d, got %d", len(expected.Sequence), len(actual.Sequence))
	}

	for i, innerExpected := range expected.Sequence {
		innerActual := actual.Sequence[i]
		compareNodes(t, &innerActual, &innerExpected)
	}

	if len(actual.Mapping) != len(expected.Mapping) {
		t.Fatalf("expected len(Mapping)=%d, got %d", len(expected.Mapping), len(actual.Mapping))
	}

	for k, innerExpected := range expected.Mapping {
		innerActual := actual.Mapping[k]
		compareNodes(t, &innerActual, &innerExpected)
	}

	if actual.Bool != expected.Bool {
		t.Fatalf("expected Bool=%t, got %t", expected.Bool, actual.Bool)
	}

	if actual.Int != expected.Int {
		t.Fatalf("expected Int=%d, got %d", expected.Int, actual.Int)
	}

	if actual.Float != expected.Float {
		t.Fatalf("expected Float=%f, got %f", expected.Float, actual.Float)
	}

	if actual.Str != expected.Str {
		t.Fatalf("expected Str=%s, got %s", expected.Str, actual.Str)
	}
}

// Test that reading a file containing a single document returns the only
// node found in that document. Checks that comments, line numbers, Kinds,
// and Tags are preserved properly.
func TestReadFile(t *testing.T) {
	testFile := getTestFile(t)

	actual := &Node{
		Name:     "test",
		FileName: testFile,
	}
	err := actual.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Node{
		Name:     "test",
		FileName: testFile,
		Line:     2,
		Column:   1,
		Comment:  "# Head Comment\n\n# Line Comment",
		Kind:     StringNode,
		Tag:      "!!str",
		Str:      "Hello World",
	}

	compareNodes(t, actual, expected)
}

// Test that reading a file containing 'null' returns a null node.
func TestReadFileNull(t *testing.T) {
	testFile := getTestFile(t)

	actual := &Node{
		Name:     "test",
		FileName: testFile,
	}
	err := actual.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Node{
		Name:     "test",
		FileName: testFile,
		Line:     1,
		Column:   1,
		Kind:     NullNode,
		Tag:      "!!null",
		Str:      "null",
	}

	compareNodes(t, actual, expected)
}

// Test that reading an empty file returns a null node.
func TestReadFileEmpty(t *testing.T) {
	testFile := getTestFile(t)

	actual := &Node{
		Name:     "test",
		FileName: testFile,
	}
	err := actual.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Node{
		Name:     "test",
		FileName: testFile,
		Line:     0,
		Column:   0,
		Kind:     NullNode,
		Tag:      "!!null",
	}

	compareNodes(t, actual, expected)
}
