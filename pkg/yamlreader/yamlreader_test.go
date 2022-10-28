package yamlreader

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestReadFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(cwd, "test", t.Name()+".yaml")

	nodes, err := ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(nodes) != 1 {
		t.Fatal("expected exactly one document, got:", len(nodes))
	}

	node := nodes[0]

	if node.Kind != yaml.DocumentNode {
		t.Fatal("expected DocumentNode, got:", node.Kind)
	}

	if len(node.Content) != 1 {
		t.Fatal("expected one child, got:", len(node.Content))
	}

	child := node.Content[0]

	if child.Kind != yaml.ScalarNode {
		t.Fatal("expected ScalarNode, got:", child.Kind)
	}

	if child.Value != "Hello World" {
		t.Fatal("expected Hello World, got:", child.Value)
	}
}
