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
		{"#test", Token{Type: TokenVariable, Value: "test"}},
		{"@test", Token{Type: TokenReference, Value: "test"}},
		{"test", Token{Type: TokenFunction, Value: "test"}},
		{"'text'", Token{Type: TokenText, Value: "text"}},
		{"1", Token{Type: TokenNumber, Value: "1"}},
		{"1.2", Token{Type: TokenNumber, Value: "1.2"}},
		{"1d", Token{Type: TokenDuration, Value: "1d"}},
		{"1h", Token{Type: TokenDuration, Value: "1h"}},
		{"1m", Token{Type: TokenDuration, Value: "1m"}},
		{"1s", Token{Type: TokenDuration, Value: "1s"}},
		{"1d2h3m4s", Token{Type: TokenDuration, Value: "1d2h3m4s"}},
		{"+", Token{Type: TokenOperator, Value: "+"}},
		{"-", Token{Type: TokenOperator, Value: "-"}},
		{"*", Token{Type: TokenOperator, Value: "*"}},
		{"/", Token{Type: TokenOperator, Value: "/"}},
		{"%", Token{Type: TokenOperator, Value: "%"}},
		{"<", Token{Type: TokenOperator, Value: "<"}},
		{">", Token{Type: TokenOperator, Value: ">"}},
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

func TestScanLocation(t *testing.T) {
	expectedResults := []Token{
		Token{Type: TokenFunction, Value: "test", LinePos: 0, LineNum: 0},
		Token{Type: TokenOpenBracket, Value: "(", LinePos: 4, LineNum: 0},
		Token{Type: TokenCloseBracket, Value: ")", LinePos: 5, LineNum: 0},
		Token{Type: TokenColon, Value: ":", LinePos: 6, LineNum: 0},
		Token{Type: TokenNewLine, Value: "\n", LinePos: 7, LineNum: 0},
		Token{Type: TokenWhitespace, Value: "  ", LinePos: 0, LineNum: 1},
		Token{Type: TokenFunction, Value: "stop", LinePos: 2, LineNum: 1},
		Token{Type: TokenOpenBracket, Value: "(", LinePos: 6, LineNum: 1},
		Token{Type: TokenCloseBracket, Value: ")", LinePos: 7, LineNum: 1},
	}
	input := "test():\n  stop()"
	scanner := NewScanner(input)
	for tok, pos := scanner.Scan(), 0; tok.Type != TokenEOF; pos++ {
		expected := expectedResults[pos]
		if tok.String() != expected.String() {
			t.Errorf("Unable to scan token #%d: expected [%s], got [%s]",
				pos,
				expected.String(),
				tok.String())
			break
		}
		if (tok.LineNum != expected.LineNum) || (tok.LinePos != expected.LinePos) {
			t.Errorf("Unexpected location for token #%d: expected [%d,%d], got [%d,%d]",
				pos,
				expected.LineNum, expected.LinePos,
				tok.LineNum, tok.LinePos)
			break
		}
		tok = scanner.Scan()
	}
}
