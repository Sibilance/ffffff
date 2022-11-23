package yamlreader

import "testing"

// Test reading a nested directory structure.
func TestReadDirectory(t *testing.T) {
	actual := readTestDirectory(t)
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
								Scalar:   StringScalar{Value: "Deep nested content"},
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
						Scalar:   StringScalar{Value: "Shallow content"},
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
				Scalar:   StringScalar{Value: "Library content"},
			},
			"main": {
				Name:     t.Name() + ".main",
				FileName: testPath + "/main.yaml",
				Line:     1,
				Column:   1,
				Kind:     StringNode,
				Tag:      "!!str",
				Raw:      "Main content",
				Scalar:   StringScalar{Value: "Main content"},
			},
		},
	}

	compareNodes(t, actual, expected)
}
