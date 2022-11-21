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

// Test that sequences are decoded correctly.
func TestSequence(t *testing.T) {
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
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "Item 1",
			},
			{
				Name:     "test[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "Item 2",
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that mappings are decoded correctly.
func TestMapping(t *testing.T) {
	actual := ReadTestFile(t)
	testFile := getTestFile(t)

	expected := &Node{
		Name:     "test",
		FileName: testFile,
		Line:     1,
		Column:   1,
		Kind:     MappingNode,
		Tag:      "!!map",
		Mapping: map[string]Node{
			"key 1": {
				Name:     "test.key 1",
				FileName: testFile,
				Line:     1,
				Column:   8,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "value 1",
			},
			"key 2": {
				Name:     "test.key 2",
				FileName: testFile,
				Line:     2,
				Column:   8,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "value 2",
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that different scalar types are decoded correctly.
func TestScalar(t *testing.T) {
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
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Str:      "true",
				Bool:     true,
			},
			{
				Name:     "test[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Str:      "false",
				Bool:     false,
			},
			{
				Name:     "test[2]",
				FileName: testFile,
				Line:     3,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Str:      "False",
				Bool:     false,
			},
			{
				Name:     "test[3]",
				FileName: testFile,
				Line:     4,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Str:      "True",
				Bool:     true,
			},
			{
				Name:     "test[4]",
				FileName: testFile,
				Line:     5,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "0",
				Int:      0,
			},
			{
				Name:     "test[5]",
				FileName: testFile,
				Line:     6,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "1",
				Int:      1,
			},
			{
				Name:     "test[6]",
				FileName: testFile,
				Line:     7,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "+1",
				Int:      1,
			},
			{
				Name:     "test[7]",
				FileName: testFile,
				Comment:  "# largest int64",
				Line:     8,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "9223372036854775807",
				Int:      9223372036854775807,
			},
			{
				Name:     "test[8]",
				FileName: testFile,
				Comment:  "# smallest int64",
				Line:     9,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "-9223372036854775808",
				Int:      -9223372036854775808,
			},
			{
				Name:     "test[9]",
				FileName: testFile,
				Comment:  "# binary",
				Line:     10,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "0b1101",
				Int:      0b1101,
			},
			{
				Name:     "test[10]",
				FileName: testFile,
				Comment:  "# octal",
				Line:     11,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "0o14",
				Int:      0o14,
			},
			{
				Name:     "test[11]",
				FileName: testFile,
				Comment:  "# hex",
				Line:     12,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Str:      "0xFF",
				Int:      0xFF,
			},
			{
				Name:     "test[12]",
				FileName: testFile,
				Comment:  "# Pi",
				Line:     13,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Str:      "3.14159",
				Float:    3.14159,
			},
			{
				Name:     "test[13]",
				FileName: testFile,
				Comment:  "# Avogadro",
				Line:     14,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Str:      "6.022e+23",
				Float:    6.022e+23,
			},
			{
				Name:     "test[14]",
				FileName: testFile,
				Comment:  "# null",
				Line:     15,
				Column:   3,
				Kind:     NullNode,
				Tag:      "!!null",
				Str:      "~",
			},
			{
				Name:     "test[15]",
				FileName: testFile,
				Line:     16,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "A string",
			},
			{
				Name:     "test[16]",
				FileName: testFile,
				Comment:  "# Also a string",
				Line:     17,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Str:      "123",
			},
		},
	}

	compareNodes(t, actual, expected)
}
