package testhelpers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func getTestPath(t *testing.T, depth int) string {
	_, testPath, _, _ := runtime.Caller(depth)
	testDir := strings.TrimSuffix(testPath, filepath.Ext(testPath))
	return filepath.Join(testDir, t.Name())
}

func GetTestFile(t *testing.T) string {
	return getTestPath(t, 2) + ".yaml"
}

func CompareNodes(actual, expected *yaml.Node) error {
	if actual.Kind != expected.Kind {
		return fmt.Errorf("expected Kind %s, got %s", kindString(expected.Kind), kindString(actual.Kind))
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

func kindString(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.ScalarNode:
		return "ScalarNode"
	case yaml.AliasNode:
		return "AliasNode"
	default:
		return fmt.Sprint(kind)
	}
}
