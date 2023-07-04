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

func GetTestYaml(t *testing.T) []*yaml.Node {
	documents, err := readyaml.ReadFile(getTestFile(t, 1))
	if err != nil {
		t.Fatal(err)
	}

	return documents
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

	if actual.HeadComment != expected.HeadComment {
		return fmt.Errorf("expected HeadComment '%s', got '%s'", expected.HeadComment, actual.HeadComment)
	}

	if actual.LineComment != expected.LineComment {
		return fmt.Errorf("expected LineComment '%s', got '%s'", expected.LineComment, actual.LineComment)
	}

	if actual.FootComment != expected.FootComment {
		return fmt.Errorf("expected FootComment '%s', got '%s'", expected.FootComment, actual.FootComment)
	}

	if len(actual.Content) != len(expected.Content) {
		return fmt.Errorf("expected %d children, got %d", len(expected.Content), len(actual.Content))
	}

	for i, expectedChild := range expected.Content {
		actualChild := actual.Content[i]
		err := CompareNodes(actualChild, expectedChild)
		if err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
	}

	return nil
}

func CompareNodeLists(actual, expected []*yaml.Node) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d nodes, got %d", len(expected), len(actual))
	}

	for i, expectedNode := range expected {
		err := CompareNodes(actual[i], expectedNode)
		if err != nil {
			return fmt.Errorf("%d: %w", i, err)
		}
	}

	return nil
}
