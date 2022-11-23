package yamlreader

import (
	"os"
	"path/filepath"
	"strings"
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

	if (actual.Value == nil) != (expected.Value == nil) {
		t.Fatalf("%s: expected (Value == nil)=%t, got (Value == nil)=%t", actual.Name, expected.Value == nil, actual.Value == nil)
	}

	if actual.Value != nil {
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

func ReadTestDirectory(t *testing.T) *Node {
	node := Node{}
	err := node.ReadDirectory(getTestPath(t))
	if err != nil {
		t.Fatal(err)
	}
	return &node
}

func ReadTestFile(t *testing.T) *Node {
	node := Node{}
	err := node.ReadFile(getTestFile(t))
	if err != nil {
		t.Fatal(err)
	}
	return &node
}

// Test reading a nested directory structure.
func TestReadDirectory(t *testing.T) {
	actual := ReadTestDirectory(t)
	testPath := getTestPath(t)

	expected := &Node{
		Name:     t.Name(),
		FileName: testPath,
		Kind:     MappingNode,
		Tag:      "!!map",
		Mapping: map[string]Node{
			"nested": {
				Name:     t.Name() + ".nested",
				FileName: testPath + "/nested",
				Kind:     MappingNode,
				Tag:      "!!map",
				Mapping: map[string]Node{
					"inner": {
						Name:     t.Name() + ".nested.inner",
						FileName: testPath + "/nested/inner",
						Kind:     MappingNode,
						Tag:      "!!map",
						Mapping: map[string]Node{
							"deep": {
								Name:     t.Name() + ".nested.inner.deep",
								FileName: testPath + "/nested/inner/deep.yaml",
								Line:     1,
								Column:   1,
								Kind:     StringNode,
								Tag:      "!!str",
								Raw:      "Deep nested content",
								Value:    StringValue{Value: "Deep nested content"},
							},
						},
					},
				},
			},
			"shallow": {
				Name:     t.Name() + ".shallow",
				FileName: testPath + "/shallow",
				Kind:     MappingNode,
				Tag:      "!!map",
				Mapping: map[string]Node{
					"shallow": {
						Name:     t.Name() + ".shallow.shallow",
						FileName: testPath + "/shallow/shallow.yaml",
						Line:     1,
						Column:   1,
						Kind:     StringNode,
						Tag:      "!!str",
						Raw:      "Shallow content",
						Value:    StringValue{Value: "Shallow content"},
					},
				},
			},
			"library": {
				Name:     t.Name() + ".library",
				FileName: testPath + "/library.yaml",
				Line:     1,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Library content",
				Value:    StringValue{Value: "Library content"},
			},
			"main": {
				Name:     t.Name() + ".main",
				FileName: testPath + "/main.yaml",
				Line:     1,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Main content",
				Value:    StringValue{Value: "Main content"},
			},
		},
	}

	compareNodes(t, actual, expected)
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
		Name:     t.Name(),
		FileName: getTestFile(t),
		Line:     2,
		Column:   1,
		Comment:  "# Head Comment\n\n# Line Comment",
		Kind:     StringNode,
		Tag:      "!!str",
		Raw:      "Hello World",
		Value:    StringValue{Value: "Hello World"},
	}

	compareNodes(t, actual, expected)
}

// Test that reading a file containing 'null' returns a null node.
func TestReadFileNull(t *testing.T) {
	actual := ReadTestFile(t)

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
	actual := ReadTestFile(t)

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
	actual := ReadTestFile(t)
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
				Value:    StringValue{Value: "Document one"},
			},
			{
				Name:     t.Name() + "[1]",
				FileName: testFile,
				Line:     4,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Document two",
				Value:    StringValue{Value: "Document two"},
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
				Value:    StringValue{Value: "Thing to be copied."},
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
				Value:    StringValue{Value: "Thing to be copied."},
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that recursive aliases are detected and return an error.
func TestReadFileRecursiveAlias(t *testing.T) {
	testFile := getTestFile(t)
	t.Fatal("some message")

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
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Item 1",
				Value:    StringValue{Value: "Item 1"},
			},
			{
				Name:     t.Name() + "[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Item 2",
				Value:    StringValue{Value: "Item 1"},
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
		Name:     t.Name(),
		FileName: testFile,
		Line:     1,
		Column:   1,
		Kind:     MappingNode,
		Tag:      "!!map",
		Mapping: map[string]Node{
			"key 1": {
				Name:     t.Name() + ".key 1",
				FileName: testFile,
				Line:     1,
				Column:   8,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "value 1",
				Value:    StringValue{Value: "value 1"},
			},
			"key 2": {
				Name:     t.Name() + ".key 2",
				FileName: testFile,
				Line:     2,
				Column:   8,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "value 2",
				Value:    StringValue{Value: "value 1"},
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
		Name:     "TestScalar",
		FileName: testFile,
		Line:     1,
		Column:   1,
		Kind:     SequenceNode,
		Tag:      "!!seq",
		Sequence: []Node{
			{
				Name:     "TestScalar[0]",
				FileName: testFile,
				Line:     1,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "true",
				Value:    BoolValue{Value: true},
			},
			{
				Name:     "TestScalar[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "false",
				Value:    BoolValue{Value: false},
			},
			{
				Name:     "TestScalar[2]",
				FileName: testFile,
				Line:     3,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "False",
				Value:    BoolValue{Value: false},
			},
			{
				Name:     "TestScalar[3]",
				FileName: testFile,
				Line:     4,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "True",
				Value:    BoolValue{Value: true},
			},
			{
				Name:     "TestScalar[4]",
				FileName: testFile,
				Line:     5,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "0",
				Value:    IntValue{Value: 0},
			},
			{
				Name:     "TestScalar[5]",
				FileName: testFile,
				Line:     6,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "1",
				Value:    IntValue{Value: 1},
			},
			{
				Name:     "TestScalar[6]",
				FileName: testFile,
				Line:     7,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "+1",
				Value:    IntValue{Value: 1},
			},
			{
				Name:     "TestScalar[7]",
				FileName: testFile,
				Comment:  "# largest int64",
				Line:     8,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "9223372036854775807",
				Value:    IntValue{Value: 9223372036854775807},
			},
			{
				Name:     "TestScalar[8]",
				FileName: testFile,
				Comment:  "# smallest int64",
				Line:     9,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "-9223372036854775808",
				Value:    IntValue{Value: -9223372036854775808},
			},
			{
				Name:     "TestScalar[9]",
				FileName: testFile,
				Comment:  "# binary",
				Line:     10,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "0b1101",
				Value:    IntValue{Value: 0b1101},
			},
			{
				Name:     "TestScalar[10]",
				FileName: testFile,
				Comment:  "# octal",
				Line:     11,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "0o14",
				Value:    IntValue{Value: 0o14},
			},
			{
				Name:     "TestScalar[11]",
				FileName: testFile,
				Comment:  "# hex",
				Line:     12,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "0xFF",
				Value:    IntValue{Value: 0xFF},
			},
			{
				Name:     "TestScalar[12]",
				FileName: testFile,
				Comment:  "# Pi",
				Line:     13,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Raw:      "3.14159",
				Value:    FloatValue{Value: 3.14159},
			},
			{
				Name:     "TestScalar[13]",
				FileName: testFile,
				Comment:  "# Avogadro",
				Line:     14,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Raw:      "6.022e+23",
				Value:    FloatValue{Value: 6.022e+23},
			},
			{
				Name:     "TestScalar[14]",
				FileName: testFile,
				Comment:  "# null",
				Line:     15,
				Column:   3,
				Kind:     NullNode,
				Tag:      "!!null",
				Raw:      "~",
			},
			{
				Name:     "TestScalar[15]",
				FileName: testFile,
				Line:     16,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "A string",
				Value:    StringValue{Value: "A string"},
			},
			{
				Name:     "TestScalar[16]",
				FileName: testFile,
				Comment:  "# Also a string",
				Line:     17,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "123",
				Value:    StringValue{Value: "123"},
			},
			{
				Name:     "TestScalar[17]",
				FileName: testFile,
				Comment:  "# overflow -> float",
				Line:     18,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Raw:      "-9223372036854775809",
				Value:    FloatValue{Value: -9223372036854775809},
			},
			{
				Name:     "TestScalar[18]",
				FileName: testFile,
				Comment:  "# overflow -> float",
				Line:     19,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Raw:      "9223372036854775808",
				Value:    FloatValue{Value: 9223372036854775808},
			},
			{
				Name:     "TestScalar[19]",
				FileName: testFile,
				Comment:  "# largest uint64 -> float",
				Line:     20,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Raw:      "18446744073709551615",
				Value:    FloatValue{Value: 18446744073709551615},
			},
			{
				Name:     "TestScalar[20]",
				FileName: testFile,
				Comment:  "# one greater -> float",
				Line:     21,
				Column:   3,
				Kind:     FloatNode,
				Tag:      "!!float",
				Raw:      "18446744073709551616",
				Value:    FloatValue{Value: 18446744073709551616},
			},
		},
	}

	compareNodes(t, actual, expected)
}
