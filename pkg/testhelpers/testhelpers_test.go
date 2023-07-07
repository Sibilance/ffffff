package testhelpers

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGetTestPath(t *testing.T) {
	testPath := getTestPath(t, 0)

	if !strings.HasSuffix(testPath, "/pkg/testhelpers/testhelpers_test/TestGetTestPath") {
		t.Fatalf("unexpected test path: %s", testPath)
	}

	contents, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != "Contents of TestGetTestPath" {
		t.Fatalf("unexpected contents: %s", contents)
	}
}

func TestGetTestFile(t *testing.T) {
	contents, err := os.ReadFile(GetTestFile(t))
	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != "This is TestGetTestFile.yaml." {
		t.Fatalf("unexpected contents: %s", contents)
	}
}

func TestGetTestYaml(t *testing.T) {
	documents := GetTestYaml(t)

	if len(documents) != 2 {
		t.Fatalf("expected two documents, got %d", len(documents))
	}

	err := CompareNodeLists(
		documents,
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "document 1",
					},
				},
			},
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "document 2",
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetYamlTestCases(t *testing.T) {
	inputs, outputs, errors := GetYamlTestCases(t, 3)

	if len(inputs) != 3 {
		t.Fatalf("expected three inputs, got %d", len(inputs))
	}
	if len(outputs) != 3 {
		t.Fatalf("expected three outputs, got %d", len(outputs))
	}
	if len(errors) != 3 {
		t.Fatalf("expected three errors, got %d", len(errors))
	}

	err := CompareNodeLists(
		inputs[0],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "input document 1.1",
					},
				},
			},
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "input document 1.2",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = CompareNodeLists(
		outputs[0],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "output document 1.1",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if errors[0] != "" {
		t.Fatal("expected no error for test case 1")
	}

	err = CompareNodeLists(
		inputs[1],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "input document 2.1",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = CompareNodeLists(
		outputs[1],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "output document 2.1",
					},
				},
			},
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "output document 2.2",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if errors[1] != "" {
		t.Fatal("expected no error for test case 2")
	}

	err = CompareNodeLists(
		inputs[2],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "input document 3.1",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = CompareNodeLists(
		outputs[2],
		[]*yaml.Node{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if errors[2] != "error message 3.1" {
		t.Fatalf("expected error for test case 3, got '%s'", errors[2])
	}
}

func TestGetYamlTestCasesNil(t *testing.T) {
	inputs, outputs, _ := GetYamlTestCases(t, 1)

	err := CompareNodeLists(
		inputs[0],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind: yaml.ScalarNode,
						Tag:  "!!null",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = CompareNodeLists(
		outputs[0],
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind: yaml.ScalarNode,
						Tag:  "!!null",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompareNodes(t *testing.T) {
	expected := yaml.Node{
		Kind:        yaml.DocumentNode,
		Tag:         "!myTag",
		Value:       "myValue",
		Anchor:      "myAnchor",
		HeadComment: "# head comment",
		LineComment: "# line comment",
		FootComment: "# foot comment",
		Content: []*yaml.Node{
			{},
		},
	}

	actual := expected
	err := CompareNodes(&actual, &expected)
	if err != nil {
		t.Fatalf("expected nodes to be identical")
	}

	actual.Kind = yaml.AliasNode
	err = CompareNodes(&actual, &expected)
	actual.Kind = expected.Kind
	if err == nil || err.Error() != "expected Kind DocumentNode, got AliasNode" {
		t.Fatal(err)
	}

	actual.Tag = "!newTag"
	err = CompareNodes(&actual, &expected)
	actual.Tag = expected.Tag
	if err == nil || err.Error() != "expected Tag !myTag, got !newTag" {
		t.Fatal(err)
	}

	actual.Value = "newValue"
	err = CompareNodes(&actual, &expected)
	actual.Value = expected.Value
	if err == nil || err.Error() != "expected Value 'myValue', got 'newValue'" {
		t.Fatal(err)
	}

	actual.Anchor = "newAnchor"
	err = CompareNodes(&actual, &expected)
	actual.Anchor = expected.Anchor
	if err == nil || err.Error() != "expected Anchor myAnchor, got newAnchor" {
		t.Fatal(err)
	}

	actual.HeadComment = "# new head comment"
	err = CompareNodes(&actual, &expected)
	actual.HeadComment = expected.HeadComment
	if err != nil {
		t.Fatalf("unexpected comment comparison, %s", err)
	}

	actual.LineComment = "# new line comment"
	err = CompareNodes(&actual, &expected)
	actual.LineComment = expected.LineComment
	if err != nil {
		t.Fatalf("unexpected comment comparison, %s", err)
	}

	actual.FootComment = "# new foot comment"
	err = CompareNodes(&actual, &expected)
	actual.FootComment = expected.FootComment
	if err != nil {
		t.Fatalf("unexpected comment comparison, %s", err)
	}

	actual.Content = append(actual.Content, &yaml.Node{})
	err = CompareNodes(&actual, &expected)
	actual.Content = expected.Content
	if err == nil || err.Error() != "expected 1 children, got 2" {
		t.Fatal(err)
	}
}

func TestCompareNodesRecursive(t *testing.T) {
	err := CompareNodes(
		&yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{
					Value: "First Value",
				},
				{
					Value: "Actual Value",
				},
			},
		},
		&yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{
					Value: "First Value",
				},
				{
					Value: "Expected Value",
				},
			},
		},
	)
	if err == nil || err.Error() != "1: expected Value 'Expected Value', got 'Actual Value'" {
		t.Fatal(err)
	}
}

func TestCompareNodeLists(t *testing.T) {
	err := CompareNodeLists(
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
			},
			{
				Kind: yaml.AliasNode,
			},
		},
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
			},
			{
				Kind: yaml.AliasNode,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = CompareNodeLists(
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
			},
			{
				Kind: yaml.AliasNode,
			},
		},
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
			},
		},
	)
	if err == nil || err.Error() != "expected 1 nodes, got 2" {
		t.Fatal(err)
	}

	err = CompareNodeLists(
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
			},
			{
				Kind: yaml.AliasNode,
			},
		},
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
			},
			{
				Kind: yaml.DocumentNode,
			},
		},
	)
	if err == nil || err.Error() != "1: expected Kind DocumentNode, got AliasNode" {
		t.Fatal(err)
	}
}
