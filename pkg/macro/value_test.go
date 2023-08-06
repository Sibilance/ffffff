package macro

import (
	"testing"

	"github.com/sibilance/ffffff/pkg/testhelpers"
)

func TestValueBool(t *testing.T) {
	inputs, outputs, _ := testhelpers.GetYamlTestCases(t, 16)

	for i, input := range inputs {
		if len(input) != 1 {
			t.Fatalf("expected exactly one input")
		}
		if len(outputs[i]) != 1 {
			t.Fatalf("expected exactly one output")
		}
		var output bool
		err := outputs[i][0].Decode(&output)
		if err != nil {
			t.Fatal(err)
		}

		actual := Value{*input[0]}.Bool()

		if actual != output {
			t.Fatalf("test %d: expected %v, got %v", i, output, actual)
		}
	}
}

func TestValueBoolAlias(t *testing.T) {
	inputs, outputs, _ := testhelpers.GetYamlTestCases(t, 2)

	for i, input := range inputs {
		if len(input) != 1 {
			t.Fatal("expected exactly one input")
		}
		if len(outputs[i]) != 1 {
			t.Fatal("expected exactly one output")
		}
		var output bool
		err := outputs[i][0].Decode(&output)
		if err != nil {
			t.Fatal(err)
		}

		if len(input[0].Content[0].Content) != 2 {
			t.Fatal("expected exactly two items")
		}

		actual := Value{*input[0].Content[0].Content[1]}.Bool()

		if actual != output {
			t.Fatalf("test %d: expected %v, got %v", i, output, actual)
		}
	}
}

func TestValueString(t *testing.T) {
	inputs, outputs, _ := testhelpers.GetYamlTestCases(t, 1)

	for i, input := range inputs {
		if len(input) != 1 {
			t.Fatal("expected exactly one input")
		}
		if len(outputs[i]) != 1 {
			t.Fatal("expected exactly one output")
		}
		var output string
		err := outputs[i][0].Decode(&output)
		if err != nil {
			t.Fatal(err)
		}

		actual := Value{*input[0]}.String()

		if actual != output {
			t.Fatalf("test %d: expected:\n%v\nbut got:\n%v", i, output, actual)
		}
	}
}
