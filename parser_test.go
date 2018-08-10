package robolang

import (
	"errors"
	"strings"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	parser := NewParser("")
	result := parser.Parse()
	if len(result.Errors) != 1 {
		t.Errorf("Unable to parse empty file: expected 1 error, found %d errors", len(result.Errors))
	} else {
		expected, actual := errors.New("Nothing to parse"), result.Errors[0]
		if expected.Error() != actual.Error() {
			t.Errorf("Unable to parse empty file: expected [%v], found [%v]", expected, actual)
		}
	}
}

func TestParseFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"clear()", "[NodeFunction]clear()"},
		{"say(text='hello')", "[NodeFunction]say([NodeArgument]text()->([NodeConstant]hello()))"},
		{"show(resource=@hello)", "[NodeFunction]show([NodeArgument]resource()->([NodeResource]hello()))"},
		{"set(variable=#count,value=1)", "[NodeFunction]set([NodeArgument]variable()->([NodeVariable]count()),[NodeArgument]value()->([NodeConstant]1()))"},
	}
	for _, test := range tests {
		t.Logf("==== Parsing `%s` ====", test.input)
		parser := NewParser(test.input)
		parser.Log = t.Logf
		result := parser.Parse()
		if len(result.Errors) > 0 {
			t.Errorf("Unexpected errors while parsing `%s`: [%v]", test.input, result.Errors)
		} else {
			nodes := make([]string, len(result.Nodes))
			for pos, node := range result.Nodes {
				nodes[pos] = node.String()
			}
			actual := strings.Join(nodes, "\n")
			if test.expected != actual {
				t.Errorf("Unable to parse `%s`: expected [%s], got [%s]",
					test.input,
					test.expected,
					actual)
			}
		}
	}
}
