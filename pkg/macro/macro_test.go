package macro

import (
	"testing"

	"github.com/sibilance/ffffff/pkg/readyaml"
	"github.com/sibilance/ffffff/pkg/testhelpers"
	"gopkg.in/yaml.v3"
)

func TestVoidDocument(t *testing.T) {
	documents, err := readyaml.ReadFile(testhelpers.GetTestFile(t))
	if err != nil {
		t.Fatal(err)
	}

	err = ProcessDocuments(&Context{}, &documents)
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.CompareNodeLists(
		documents,
		[]*yaml.Node{
			{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "first non-void document",
					},
				},
			},
		},
	)
}
