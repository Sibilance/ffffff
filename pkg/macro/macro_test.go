package macro

import (
	"testing"

	"github.com/sibilance/ffffff/pkg/testhelpers"
	"gopkg.in/yaml.v3"
)

func TestVoidDocument(t *testing.T) {
	inputs, _ := testhelpers.GetYamlTestCases(t, 1)

	err := ProcessDocuments(&Context{}, &inputs[0])
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.CompareNodeLists(
		inputs[0],
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
