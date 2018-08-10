package robolang

import "testing"

func TestScanTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected Token
	}{
		{"", Token{Type: EOF, Value: string(rune(0))}},
		{" ", Token{Type: WHITESPACE, Value: " "}},
		{"\t", Token{Type: WHITESPACE, Value: "\t"}},
		{"\r", Token{Type: WHITESPACE, Value: "\r"}},
		{"\n", Token{Type: NEWLINE, Value: "\n"}},
		{":", Token{Type: COLON, Value: ":"}},
		{",", Token{Type: COMMA, Value: ","}},
		{"(", Token{Type: OPENBRACKET, Value: "("}},
		{")", Token{Type: CLOSEBRACKET, Value: ")"}},
		{"=", Token{Type: EQUALS, Value: "="}},
		{"_test", Token{Type: VARIABLE, Value: "test"}},
		{"@test", Token{Type: REFERENCE, Value: "test"}},
		{"test", Token{Type: FUNCTION, Value: "test"}},
		{"'text'", Token{Type: TEXT, Value: "text"}},
		{"1", Token{Type: NUMBER, Value: "1"}},
		{"1.2", Token{Type: NUMBER, Value: "1.2"}},
		{"-1", Token{Type: NUMBER, Value: "-1"}},
		{"-1.2", Token{Type: NUMBER, Value: "-1.2"}},
		{"1d", Token{Type: DURATION, Value: "1d"}},
		{"1h", Token{Type: DURATION, Value: "1h"}},
		{"1m", Token{Type: DURATION, Value: "1m"}},
		{"1s", Token{Type: DURATION, Value: "1s"}},
		{"1d2h3m4s", Token{Type: DURATION, Value: "1d2h3m4s"}},
	}

	for _, test := range tests {
		scanner := NewScanner(test.input)
		actual := scanner.Scan()
		if actual.String() != test.expected.String() {
			t.Errorf("Unable to scan [%s]: expected [%s], got [%s]",
				test.input,
				test.expected.String(),
				actual.String())
		}
	}
}
