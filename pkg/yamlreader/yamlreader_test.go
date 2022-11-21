package yamlreader

import (
	"os"
	"path/filepath"
	"strings"
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

	if actual.Bool != expected.Bool {
		t.Fatalf("%s: expected Bool=%t, got %t", actual.Name, expected.Bool, actual.Bool)
	}

	if actual.Int != expected.Int {
		t.Fatalf("%s: expected Int=%d, got %d", actual.Name, expected.Int, actual.Int)
	}

	if actual.Float != expected.Float {
		t.Fatalf("%s: expected Float=%f, got %f", actual.Name, expected.Float, actual.Float)
	}

	if actual.Str != expected.Str {
		t.Fatalf("%s: expected Str=%s, got %s", actual.Name, expected.Str, actual.Str)
	}
}

func ReadTestFile(t *testing.T) *Node {
	testFile := getTestFile(t)

	actual := &Node{
		Name: "test",
	}
	err := actual.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	return actual
}

// Test that reading a missing file returns an error.
func TestReadFileMissing(t *testing.T) {
	testFile := getTestFile(t)

	actual := &Node{}
	err := actual.ReadFile(testFile)
	if err == nil {
		t.Fatalf("no error returned for missing file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("wrong error returned: %s", err.Error())
	}
}

// Test that reading a file containing a single document returns the only
// node found in that document. Checks that comments, line numbers, Kinds,
// and Tags are preserved properly.
func TestReadFile(t *testing.T) {
	actual := ReadTestFile(t)

	expected := &Node{
		Name:     "test",
		FileName: getTestFile(t),
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
	actual := ReadTestFile(t)

	expected := &Node{
		Name:     "test",
		FileName: getTestFile(t),
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
	actual := ReadTestFile(t)

	expected := &Node{
		Name:     "test",
		FileName: getTestFile(t),
		Kind:     NullNode,
		Tag:      "!!null",
	}

	compareNodes(t, actual, expected)
}

// Test that reading a file with multiple documents returns a sequence node.
func TestReadFileMultipleDocuments(t *testing.T) {
	actual := ReadTestFile(t)
	testFile := getTestFile(t)

	expected := &Node{
		Name:     "test",
		FileName: testFile,
		Kind:     SequenceNode,
		Tag:      "!!seq",
		Sequence: []Node{
			{
				Name:     "test[0]",
				FileName: testFile,
				Line:     2,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "Document one",
			},
			{
				Name:     "test[1]",
				FileName: testFile,
				Line:     4,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "Document two",
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that aliases are followed correctly.
func TestReadFileAlias(t *testing.T) {
	actual := ReadTestFile(t)
	testFile := getTestFile(t)

	expected := &Node{
		Name:     "test",
		FileName: testFile,
		Line:     1,
		Column:   1,
		Kind:     SequenceNode,
		Tag:      "!!seq",
		Sequence: []Node{
			{
				Name:     "test[0]",
				FileName: testFile,
				Line:     1,
				Column:   3,
				Comment:  "# Comment 1",
				Kind:     StringNode,
				Tag:      "!MyTag",
				Str:      "Thing to be copied.",
			},
			{
				Name:     "test[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Comment:  "# Comment 2",
				Kind:     StringNode,
				Tag:      "!MyTag",
				Str:      "Thing to be copied.",
			},
		},
	}

	compareNodes(t, actual, expected)
}

func TestReadFileRecursiveAlias(t *testing.T) {
	testFile := getTestFile(t)

	actual := &Node{}
	err := actual.ReadFile(testFile)
	if err == nil {
		t.Fatalf("no error returned for recursive alias")
	}
	if !strings.Contains(err.Error(), "recursive alias detected") {
		t.Fatalf("wrong error returned: %s", err.Error())
	}
}
