package macro

import (
	"fmt"

	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

func ProcessDocuments(context *Context, documents *[]*yaml.Node) error {
	originalDocuments := *documents
	*documents = nil

	for i, document := range originalDocuments {
		localContext := context.New(fmt.Sprintf("[%d]", i))
		if document.Kind != yaml.DocumentNode {
			return localContext.Error(
				document,
				fmt.Sprintf("expected DocumentNode, got %s", yamlhelpers.KindString(document.Kind)),
			)
		}
		if len(document.Content) != 1 {
			return localContext.Error(document, "expected exactly one child node of document")
		}
		child := document.Content[0]
		err := ProcessNode(localContext, child)
		if err != nil {
			return err
		}
		if IsUnwrap(child) {
			switch child.Kind {
			case yaml.SequenceNode:
				for _, newChild := range child.Content {
					*documents = append(*documents, &yaml.Node{
						Kind: yaml.DocumentNode,
						Content: []*yaml.Node{
							newChild,
						},
						Line:   newChild.Line,
						Column: newChild.Column,
					})
				}
			case yaml.MappingNode:
				for i, newValue := range child.Content {
					if i&1 == 0 {
						continue
					}
					newKey := child.Content[i-1]
					*documents = append(*documents, &yaml.Node{
						Kind: yaml.DocumentNode,
						Content: []*yaml.Node{
							{
								Kind: yaml.MappingNode,
								Tag:  "!!map",
								Content: []*yaml.Node{
									newKey,
									newValue,
								},
								Line:   newKey.Line,
								Column: newKey.Column,
							},
						},
						Line:   newKey.Line,
						Column: newKey.Column,
					})
				}
			default:
				return localContext.Error(child, fmt.Sprintf("cannot unwrap %s", yamlhelpers.KindString(child.Kind)))
			}
		} else if !IsVoid(child) {
			*documents = append(*documents, document)
		}
	}

	return nil
}
