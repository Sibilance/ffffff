package testhelpers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sibilance/ffffff/pkg/readyaml"
	"github.com/sibilance/ffffff/pkg/yamlhelpers"
	"gopkg.in/yaml.v3"
)

func getTestPath(t *testing.T, depth int) string {
	_, testPath, _, _ := runtime.Caller(depth + 1)
	testDir := strings.TrimSuffix(testPath, filepath.Ext(testPath))
	return filepath.Join(testDir, t.Name())
}

func getTestFile(t *testing.T, depth int) string {
	return getTestPath(t, depth+1) + ".yaml"
}

func GetTestFile(t *testing.T) string {
	return getTestFile(t, 1)
}

func getTestYaml(t *testing.T, depth int) []*yaml.Node {
	documents, err := readyaml.ReadFile(getTestFile(t, depth+1))
	if err != nil {
		t.Fatal(err)
	}

	return documents
}

func GetTestYaml(t *testing.T) []*yaml.Node {
	return getTestYaml(t, 1)
}

func GetYamlTestCases(t *testing.T, count int) (inputs, outputs [][]*yaml.Node, errors []string) {
	type Mode int
	const (
		inputPrefix  = "# Input"
		outputPrefix = "# Output"
		errorPrefix  = "# Error"
	)
	const (
		inputMode Mode = iota
		outputMode
		errorMode
	)
	documents := getTestYaml(t, 1)

	if len(documents) < 1 {
		t.Fatal("expected at least one document")
	}
	if len(documents[0].Content) != 1 {
		t.Fatal("expected exactly one child of document")
	}

	mode := inputMode
	var inputDocuments, outputDocuments []*yaml.Node
	var errorMessage string
	for i, document := range documents {
		if mode == inputMode {
			inputDocuments = append(inputDocuments, document)
		} else if mode == outputMode {
			outputDocuments = append(outputDocuments, document)
		} else if mode == errorMode {
			if len(document.Content) != 1 {
				t.Fatal("expected exactly one child of document")
			}
			child := document.Content[0]
			if child.Kind != yaml.ScalarNode {
				t.Fatal("expected string error message")
			}
			errorMessage = child.Value
		}

		footComment := document.FootComment
		if i == 0 && strings.HasPrefix(footComment, inputPrefix) {
			// A bug in the yaml library sometimes attaches comments above the first document
			// as foot comments to that document. If this happens, try to remove it by searching
			// for `outputPrefix` and `errorPrefix`.
			index := strings.Index(footComment, "\n"+outputPrefix)
			if index < 0 {
				index = strings.Index(footComment, "\n"+errorPrefix)
			}
			if index > 0 {
				footComment = footComment[index+1:]
			}
		}
		if strings.HasPrefix(footComment, inputPrefix) {
			if mode == outputMode || mode == errorMode {
				mode = inputMode
				inputs = append(inputs, inputDocuments)
				outputs = append(outputs, outputDocuments)
				errors = append(errors, errorMessage)
				inputDocuments = nil
				outputDocuments = nil
				errorMessage = ""
			}
		}
		if strings.HasPrefix(footComment, outputPrefix) {
			mode = outputMode
		}
		if strings.HasPrefix(footComment, errorPrefix) {
			mode = errorMode
		}
	}

	inputs = append(inputs, inputDocuments)
	outputs = append(outputs, outputDocuments)
	errors = append(errors, errorMessage)

	if len(inputs) != count {
		t.Fatalf("expected %d inputs, got %d", count, len(inputs))
	}
	if len(outputs) != count {
		t.Fatalf("expected %d outputs, got %d", count, len(outputs))
	}
	return
}

func CompareNodes(actual, expected *yaml.Node) error {
	if actual.Kind != expected.Kind {
		return fmt.Errorf("expected Kind %s, got %s", yamlhelpers.KindString(expected.Kind), yamlhelpers.KindString(actual.Kind))
	}

	if actual.Tag != expected.Tag {
		return fmt.Errorf("expected Tag %s, got %s", expected.Tag, actual.Tag)
	}

	if actual.Value != expected.Value {
		return fmt.Errorf("expected Value '%s', got '%s'", expected.Value, actual.Value)
	}

	if actual.Anchor != expected.Anchor {
		return fmt.Errorf("expected Anchor %s, got %s", expected.Anchor, actual.Anchor)
	}

	for i, expectedChild := range expected.Content {
		if i >= len(actual.Content) {
			break
		}
		actualChild := actual.Content[i]
		err := CompareNodes(actualChild, expectedChild)
		if err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
	}

	if len(actual.Content) != len(expected.Content) {
		return fmt.Errorf("expected %d children, got %d", len(expected.Content), len(actual.Content))
	}

	return nil
}

func CompareNodeLists(actual, expected []*yaml.Node) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d nodes, got %d", len(expected), len(actual))
	}

	for i, expectedNode := range expected {
		if i >= len(actual) {
			break
		}
		err := CompareNodes(actual[i], expectedNode)
		if err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
	}

	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d nodes, got %d", len(expected), len(actual))
	}

	return nil
}
