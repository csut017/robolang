package robolang

// Token contains a token from the scanner.
type Token struct {
	Type     TokenType `json:"-"`
	TypeName string    `json:"type"`
	Value    string    `json:"value"`
	LineNum  int       `json:"lineNum"`
	LinePos  int       `json:"linePos"`
}

// String converts the token to a human-readable form.
func (t *Token) String() string {
	return "`" + t.Value + "` [" + t.Type.String() + "]"
}

// TokenType defines what the token can be used for.
type TokenType int

//go:generate stringer -type=TokenType

const (
	// ILLEGAL is an invalid token.
	ILLEGAL TokenType = iota
	// EOF is the end of the file
	EOF
	// WHITESPACE are whitespace characters
	WHITESPACE
)
