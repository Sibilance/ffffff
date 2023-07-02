package testhelpers

import (
	"os"
	"strings"
	"testing"
)

func TestGetTestPath(t *testing.T) {
	testPath := getTestPath(t, 1)

	if !strings.HasSuffix(testPath, "/pkg/testhelpers/helpers_test/TestGetTestPath") {
		t.Fatalf("unexpected test path: %s", testPath)
	}

	contents, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != "Contents of TestGetTestPath" {
		t.Fatalf("unexpected contents: %s", contents)
	}
}

func TestGetTestFile(t *testing.T) {
	contents, err := os.ReadFile(GetTestFile(t))
	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != "This is TestGetTestFile.yaml." {
		t.Fatalf("unexpected contents: %s", contents)
	}
}
