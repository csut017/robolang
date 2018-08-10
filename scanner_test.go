package robolang

import "testing"

func TestScanTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected Token
	}{
		{"", Token{Type: TokenEOF, Value: string(rune(0))}},
		{" ", Token{Type: TokenWhitespace, Value: " "}},
		{"\t", Token{Type: TokenWhitespace, Value: "\t"}},
		{"\r", Token{Type: TokenWhitespace, Value: "\r"}},
		{"\n", Token{Type: TokenNewLine, Value: "\n"}},
		{":", Token{Type: TokenColon, Value: ":"}},
		{",", Token{Type: TokenComma, Value: ","}},
		{"(", Token{Type: TokenOpenBracket, Value: "("}},
		{")", Token{Type: TokenCloseBracket, Value: ")"}},
		{"=", Token{Type: TokenEquals, Value: "="}},
		{"_test", Token{Type: TokenVariable, Value: "test"}},
		{"@test", Token{Type: TokenReference, Value: "test"}},
		{"test", Token{Type: TokenFunction, Value: "test"}},
		{"'text'", Token{Type: TokenText, Value: "text"}},
		{"1", Token{Type: TokenNumber, Value: "1"}},
		{"1.2", Token{Type: TokenNumber, Value: "1.2"}},
		{"-1", Token{Type: TokenNumber, Value: "-1"}},
		{"-1.2", Token{Type: TokenNumber, Value: "-1.2"}},
		{"1d", Token{Type: TokenDuration, Value: "1d"}},
		{"1h", Token{Type: TokenDuration, Value: "1h"}},
		{"1m", Token{Type: TokenDuration, Value: "1m"}},
		{"1s", Token{Type: TokenDuration, Value: "1s"}},
		{"1d2h3m4s", Token{Type: TokenDuration, Value: "1d2h3m4s"}},
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
