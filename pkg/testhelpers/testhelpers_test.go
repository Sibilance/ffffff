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
		t.Fatalf("unexpected error, %s", err)
	}
}

func TestCompareNodes(t *testing.T) {
	var err error
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
	err = CompareNodes(&actual, &expected)
	if err != nil {
		t.Fatalf("expected nodes to be identical")
	}

	actual.Kind = yaml.AliasNode
	err = CompareNodes(&actual, &expected)
	actual.Kind = expected.Kind
	if err == nil || err.Error() != "expected Kind DocumentNode, got AliasNode" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.Tag = "!newTag"
	err = CompareNodes(&actual, &expected)
	actual.Tag = expected.Tag
	if err == nil || err.Error() != "expected Tag !myTag, got !newTag" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.Value = "newValue"
	err = CompareNodes(&actual, &expected)
	actual.Value = expected.Value
	if err == nil || err.Error() != "expected Value 'myValue', got 'newValue'" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.Anchor = "newAnchor"
	err = CompareNodes(&actual, &expected)
	actual.Anchor = expected.Anchor
	if err == nil || err.Error() != "expected Anchor myAnchor, got newAnchor" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.HeadComment = "# new head comment"
	err = CompareNodes(&actual, &expected)
	actual.HeadComment = expected.HeadComment
	if err == nil || err.Error() != "expected HeadComment '# head comment', got '# new head comment'" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.LineComment = "# new line comment"
	err = CompareNodes(&actual, &expected)
	actual.LineComment = expected.LineComment
	if err == nil || err.Error() != "expected LineComment '# line comment', got '# new line comment'" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.FootComment = "# new foot comment"
	err = CompareNodes(&actual, &expected)
	actual.FootComment = expected.FootComment
	if err == nil || err.Error() != "expected FootComment '# foot comment', got '# new foot comment'" {
		t.Fatalf("unexpected error, %s", err)
	}

	actual.Content = append(actual.Content, &yaml.Node{})
	err = CompareNodes(&actual, &expected)
	actual.Content = expected.Content
	if err == nil || err.Error() != "expected 1 children, got 2" {
		t.Fatalf("unexpected error, %s", err)
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
		t.Fatalf("unexpected error, %s", err)
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
		t.Fatalf("unexpected error, %s", err)
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
		t.Fatalf("unexpected error, %s", err)
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
		t.Fatalf("unexpected error, %s", err)
	}
}
