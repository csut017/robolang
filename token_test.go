package robolang

import "testing"

func TestString(t *testing.T) {
	tok := Token{
		Type:  TokenWhitespace,
		Value: " ",
	}
	actual, expected := tok.String(), "` ` [TokenWhitespace]"
	if actual != expected {
		t.Errorf("Token.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}
