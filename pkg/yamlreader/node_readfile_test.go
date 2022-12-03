package yamlreader

import (
	"strings"
	"testing"
)

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
	actual := readTestFile(t)

	expected := &Node{
		Name:     t.Name(),
		FileName: getTestFile(t),
		Line:     2,
		Column:   1,
		Comment:  "# Head Comment\n# Line Comment",
		Kind:     StringNode,
		Tag:      "!!str",
		Raw:      "Hello World",
		Scalar:   StringScalar{Value: "Hello World"},
	}

	compareNodes(t, actual, expected)
}

// Test that reading a file containing 'null' returns a null node.
func TestReadFileNull(t *testing.T) {
	actual := readTestFile(t)

	expected := &Node{
		Name:     t.Name(),
		FileName: getTestFile(t),
		Line:     1,
		Column:   1,
		Kind:     NullNode,
		Tag:      "!!null",
		Raw:      "null",
	}

	compareNodes(t, actual, expected)
}

// Test that reading an empty file returns a null node.
func TestReadFileEmpty(t *testing.T) {
	actual := readTestFile(t)

	expected := &Node{
		Name:     t.Name(),
		FileName: getTestFile(t),
		Kind:     NullNode,
		Tag:      "!!null",
	}

	compareNodes(t, actual, expected)
}

// Test that reading a file with multiple documents returns a sequence node.
func TestReadFileMultipleDocuments(t *testing.T) {
	actual := readTestFile(t)
	testFile := getTestFile(t)

	expected := &Node{
		Name:     t.Name(),
		FileName: testFile,
		Kind:     SequenceNode,
		Tag:      "!!seq",
		Sequence: []Node{
			{
				Name:     t.Name() + "[0]",
				FileName: testFile,
				Line:     2,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Document one",
				Scalar:   StringScalar{Value: "Document one"},
			},
			{
				Name:     t.Name() + "[1]",
				FileName: testFile,
				Line:     4,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Document two",
				Scalar:   StringScalar{Value: "Document two"},
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that aliases are followed correctly.
func TestReadFileAlias(t *testing.T) {
	actual := readTestFile(t)
	testFile := getTestFile(t)

	expected := &Node{
		Name:     t.Name(),
		FileName: testFile,
		Line:     1,
		Column:   1,
		Kind:     SequenceNode,
		Tag:      "!!seq",
		Sequence: []Node{
			{
				Name:     t.Name() + "[0]",
				FileName: testFile,
				Line:     1,
				Column:   3,
				Comment:  "# Comment 1",
				Kind:     StringNode,
				Tag:      "!MyTag",
				Raw:      "Thing to be copied.",
				Scalar:   StringScalar{Value: "Thing to be copied."},
			},
			{
				Name:     t.Name() + "[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Comment:  "# Comment 2",
				Kind:     StringNode,
				Tag:      "!MyTag",
				Raw:      "Thing to be copied.",
				Scalar:   StringScalar{Value: "Thing to be copied."},
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that recursive aliases are detected and return an error.
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
