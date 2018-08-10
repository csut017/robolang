package robolang

import "testing"

func TestString(t *testing.T) {
	tok := Token{
		Type:  WHITESPACE,
		Value: " ",
	}
	actual, expected := tok.String(), "` ` [WHITESPACE]"
	if actual != expected {
		t.Errorf("Token.String() output does not match: expected [%s], got [%s]", expected, actual)
	}
}
