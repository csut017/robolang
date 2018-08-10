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
	// ILLEGAL is an invalid token
	ILLEGAL TokenType = iota

	// EOF is the end of the file
	EOF

	// NEWLINE is the end of a line (\n)
	NEWLINE

	// WHITESPACE are whitespace characters
	WHITESPACE

	// OPENBRACKET is an opening bracket sign (()
	OPENBRACKET

	// CLOSEBRACKET is a closing bracket sign ())
	CLOSEBRACKET

	// EQUALS is an equals sign (=)
	EQUALS

	// COMMA is a comma sign (,)
	COMMA

	// COLON is a colon sign (:)
	COLON

	// FUNCTION defines a function name (name)
	FUNCTION

	// VARIABLE defines a variable name (_name)
	VARIABLE

	// REFERENCE defines a reference name (@name)
	REFERENCE

	// TEXT is a string constant ('text')
	TEXT

	// NUMBER is a numeric constant (1)
	NUMBER

	// DURATION is a timespan (1d2h3m4s)
	DURATION
)
