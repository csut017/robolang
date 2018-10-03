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

func TestDoubleParsing(t *testing.T) {
	input := "clear()"
	parser := NewParser(input)
	parser.Parse()
	result := parser.Parse()
	expected := "NodeFunction:clear"
	compareResults(t, input, expected, result)
}

func TestParseFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"clear()", "NodeFunction:clear"},
		{"say(text='hello')", "NodeFunction:say(NodeArgument:text->(NodeConstant:hello))"},
		{"say(text='hello') # test\nsay(text='world')", "NodeFunction:say(NodeArgument:text->(NodeConstant:hello))\nNodeFunction:say(NodeArgument:text->(NodeConstant:world))"},
		{"show(resource=@hello)", "NodeFunction:show(NodeArgument:resource->(NodeResource:hello))"},
		{"set(variable=&count,value=1)", "NodeFunction:set(NodeArgument:variable->(NodeVariable:count),NodeArgument:value->(NodeConstant:1))"},
		{"waitForTime(duration=5m)", "NodeFunction:waitForTime(NodeArgument:duration->(NodeConstant:5m))"},
		{"waitForInput():\n  clear()", "NodeFunction:waitForInput->(NodeFunction:clear)"},
		{"clear()\nsay(text=@hello)", "NodeFunction:clear\nNodeFunction:say(NodeArgument:text->(NodeResource:hello))"},
		{"waitForInput():\n  clear()\n  say(text=@hello)", "NodeFunction:waitForInput->(NodeFunction:clear,NodeFunction:say(NodeArgument:text->(NodeResource:hello)))"},
		{"waitForInput():\n  #clear()\n  say(text=@hello)", "NodeFunction:waitForInput->(NodeFunction:say(NodeArgument:text->(NodeResource:hello)))"},
	}
	for _, test := range tests {
		t.Logf("==== Parsing `%s` ====", test.input)
		parser := NewParser(test.input)
		// parser.Log = t.Logf
		result := parser.Parse()
		compareResults(t, test.input, test.expected, result)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"clear", "Unexpected token '<EOF>', expected TokenOpenBracket"},
		{"clear(", "Unexpected token '<EOF>', expected TokenCloseBracket or TokenIdentifier"},
		{"clear--", "Unexpected token '-', expected TokenOpenBracket"},
	}
	for _, test := range tests {
		t.Logf("==== Parsing `%s` ====", test.input)
		parser := NewParser(test.input)
		// parser.Log = t.Logf
		result := parser.Parse()
		if len(result.Errors) != 1 {
			t.Errorf("Unable to parse `%s`: expected 1 error, found %d errors", test.input, len(result.Errors))
		} else {
			expected, actual := errors.New(test.expected), result.Errors[0]
			var msg string
			if pe, ok := actual.(*ParseError); ok {
				msg = pe.Message
			} else {
				msg = actual.Error()
			}
			if expected.Error() != msg {
				t.Errorf("Unable to parse `%s`: expected [%v], found [%v]", test.input, expected, msg)
			}
		}
	}
}

func compareResults(t *testing.T, input, expected string, result *ParseResult) {
	if len(result.Errors) > 0 {
		t.Errorf("Unexpected errors while parsing `%s`: [%v]", input, result.Errors)
		return
	}

	nodes := make([]string, len(result.Nodes))
	for pos, node := range result.Nodes {
		nodes[pos] = node.String()
	}
	actual := strings.Join(nodes, "\n")
	if expected != actual {
		t.Errorf("Unable to parse `%s`: expected [%s], got [%s]",
			input,
			expected,
			actual)
	}
}
