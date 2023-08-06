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

		err := macro.ProcessDocuments(macro.DefaultRootContext(), &input)

		if errorMessage != "" {
			if err == nil {
				t.Fatalf("test %d: no error, expected '%s'", i, errorMessage)
			} else if err.Error() != errorMessage {
				t.Fatalf("test %d: expected error '%s', got '%s'", i, errorMessage, err)
			}
		} else {
			if err != nil {
				t.Fatalf("test %d: %s", i, err)
			}

			err = testhelpers.CompareNodeLists(
				input,
				output,
			)
			if err != nil {
				t.Fatalf("test %d: %s", i, err)
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
	testProcessDocuments(t, 3)
}

func TestUnwrapMapping(t *testing.T) {
	testProcessDocuments(t, 2)
}
