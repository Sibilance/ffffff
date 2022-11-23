package yamlreader

import (
	"testing"
)

// Test that sequences are decoded correctly.
func TestSequence(t *testing.T) {
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
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Item 1",
				Scalar:   StringScalar{Value: "Item 1"},
			},
			{
				Name:     t.Name() + "[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Item 2",
				Scalar:   StringScalar{Value: "Item 1"},
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that mappings are decoded correctly.
func TestMapping(t *testing.T) {
	actual := readTestFile(t)
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
				Scalar:   StringScalar{Value: "value 1"},
			},
			"key 2": {
				Name:     t.Name() + ".key 2",
				FileName: testFile,
				Line:     2,
				Column:   8,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "value 2",
				Scalar:   StringScalar{Value: "value 1"},
			},
		},
	}

	compareNodes(t, actual, expected)
}

// Test that different scalar types are decoded correctly.
func TestScalar(t *testing.T) {
	actual := readTestFile(t)
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
				Scalar:   BoolScalar{Value: true},
			},
			{
				Name:     "TestScalar[1]",
				FileName: testFile,
				Line:     2,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "false",
				Scalar:   BoolScalar{Value: false},
			},
			{
				Name:     "TestScalar[2]",
				FileName: testFile,
				Line:     3,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "False",
				Scalar:   BoolScalar{Value: false},
			},
			{
				Name:     "TestScalar[3]",
				FileName: testFile,
				Line:     4,
				Column:   3,
				Kind:     BooleanNode,
				Tag:      "!!bool",
				Raw:      "True",
				Scalar:   BoolScalar{Value: true},
			},
			{
				Name:     "TestScalar[4]",
				FileName: testFile,
				Line:     5,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "0",
				Scalar:   IntScalar{Value: 0},
			},
			{
				Name:     "TestScalar[5]",
				FileName: testFile,
				Line:     6,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "1",
				Scalar:   IntScalar{Value: 1},
			},
			{
				Name:     "TestScalar[6]",
				FileName: testFile,
				Line:     7,
				Column:   3,
				Kind:     IntegerNode,
				Tag:      "!!int",
				Raw:      "+1",
				Scalar:   IntScalar{Value: 1},
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
				Scalar:   IntScalar{Value: 9223372036854775807},
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
				Scalar:   IntScalar{Value: -9223372036854775808},
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
				Scalar:   IntScalar{Value: 0b1101},
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
				Scalar:   IntScalar{Value: 0o14},
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
				Scalar:   IntScalar{Value: 0xFF},
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
				Scalar:   FloatScalar{Value: 3.14159},
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
				Scalar:   FloatScalar{Value: 6.022e+23},
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
				Scalar:   StringScalar{Value: "A string"},
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
				Scalar:   StringScalar{Value: "123"},
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
				Scalar:   FloatScalar{Value: -9223372036854775809},
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
				Scalar:   FloatScalar{Value: 9223372036854775808},
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
				Scalar:   FloatScalar{Value: 18446744073709551615},
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
				Scalar:   FloatScalar{Value: 18446744073709551616},
			},
		},
	}

	compareNodes(t, actual, expected)
}
