package yamlreader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(cwd, "test", t.Name()+".yaml")

	node, err := ReadFile("test", testFile)
	if err != nil {
		t.Fatal(err)
	}

	if node.Name != "test" {
		t.Fatalf("expected Name=test, got %s", node.Name)
	}

	if node.FileName != testFile {
		t.Fatalf("expected FileName=%s, got %s", testFile, node.FileName)
	}

	if node.Line != 2 {
		t.Fatalf("expected Line=2, got %d", node.Line)
	}

	if node.Column != 1 {
		t.Fatalf("expected Column=1, got %d", node.Column)
	}

	if node.Comment != "# Head Comment\n\n# Line Comment" {
		t.Fatalf("unexpected comment:\n%s", node.Comment)
	}

	if node.Kind != ScalarNode {
		t.Fatalf("expected Kind=ScalarNode, got %s", node.Kind)
	}

	if node.Tag != "!!str" {
		t.Fatalf("expected Tag=str, got %s", node.Tag)
	}

	if node.Scalar != "Hello World" {
		t.Fatalf("expected Scalar=Hello World, got %s", node.Scalar)
	}
}
