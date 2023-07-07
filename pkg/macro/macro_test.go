package macro_test

import (
	"testing"

	"github.com/sibilance/ffffff/pkg/macro"
	"github.com/sibilance/ffffff/pkg/testhelpers"
)

func testProcessDocuments(t *testing.T, count int) {
	inputs, outputs, errors := testhelpers.GetYamlTestCases(t, count)

	for i, input := range inputs {
		output := outputs[i]
		errorMessage := errors[i]

		err := macro.ProcessDocuments(&macro.Context{}, &input)

		if errorMessage != "" {
			if err == nil {
				t.Fatalf("no error, expected '%s'", errorMessage)
			} else if err.Error() != errorMessage {
				t.Fatalf("expected error '%s', got '%s'", errorMessage, err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}

			err = testhelpers.CompareNodeLists(
				input,
				output,
			)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestVoidDocument(t *testing.T) {
	testProcessDocuments(t, 1)
}

func TestVoidSequence(t *testing.T) {
	testProcessDocuments(t, 1)
}

func TestVoidMapping(t *testing.T) {
	testProcessDocuments(t, 1)
}

func TestVoidNested(t *testing.T) {
	testProcessDocuments(t, 3)
}

func TestUnwrapDocument(t *testing.T) {
	testProcessDocuments(t, 3)
}

func TestUnwrapSequence(t *testing.T) {
	testProcessDocuments(t, 4)
}
