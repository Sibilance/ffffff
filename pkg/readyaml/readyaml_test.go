package readyaml_test // Separate package to avoid dependency cycle with testhelpers.

import (
	"fmt"
	"testing"

	"github.com/sibilance/ffffff/pkg/readyaml"
	"github.com/sibilance/ffffff/pkg/testhelpers"
	"gopkg.in/yaml.v3"
)

func TestReadFile(t *testing.T) {
	documents, err := readyaml.ReadFile(testhelpers.GetTestFile(t))
	if err != nil {
		t.Fatal(err)
	}

	if len(documents) != 2 {
		t.Fatalf("unexpected number of documents: %d", len(documents))
	}

	expected := []*yaml.Node{
		{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: "Document 1 contents",
				},
			},
		},
		{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: "Document 2 contents",
				},
			},
		},
	}

	for i := range expected {
		if err := testhelpers.CompareNodes(documents[i], expected[i]); err != nil {
			t.Fatalf("unexpected difference in document %d: %s", i, err.Error())
		}
	}
}

func TestReadFileMissing(t *testing.T) {
	fileName := testhelpers.GetTestFile(t)
	_, err := readyaml.ReadFile(fileName)

	if err == nil || err.Error() != fmt.Sprintf("open %s: no such file or directory", fileName) {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestReadFileMap(t *testing.T) {
	documents, err := readyaml.ReadFile(testhelpers.GetTestFile(t))
	if err != nil {
		t.Fatal(err)
	}

	if len(documents) != 1 {
		t.Fatalf("wrong number of documents: %d", len(documents))
	}

	document := documents[0]

	err = testhelpers.CompareNodes(document, &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "Key",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "Value",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "Foo",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "Bar",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
