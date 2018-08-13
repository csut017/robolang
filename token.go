package robolang

// Token contains a token from the scanner.
type Token struct {
	LineNum  int       `json:"lineNum"`
	LinePos  int       `json:"linePos"`
	Type     TokenType `json:"-"`
	TypeName string    `json:"type"`
	Value    string    `json:"value"`
}

// String converts the token to a human-readable form.
func (t *Token) String() string {
	return "`" + t.Value + "` [" + t.Type.String() + "]"
}

// TokenType defines what the token can be used for.
type TokenType int

//go:generate stringer -type=TokenType

const (
	// TokenIllegal is an invalid token
	TokenIllegal TokenType = iota

	// TokenEOF is the end of the file
	TokenEOF

	// TokenNewLine is the end of a line (\n)
	TokenNewLine

	// TokenWhitespace are whitespace characters
	TokenWhitespace

	// TokenOpenBracket is an opening bracket sign (()
	TokenOpenBracket

	// TokenCloseBracket is a closing bracket sign ())
	TokenCloseBracket

	// TokenEquals is an equals sign (=)
	TokenEquals

	// TokenComma is a comma sign (,)
	TokenComma

	// TokenColon is a colon sign (:)
	TokenColon

	// TokenIdentifier defines an identifier (name)
	TokenIdentifier

	// TokenVariable defines a variable name (_name)
	TokenVariable

	// TokenResource defines a resource name (@name)
	TokenResource

	// TokenText is a string constant ('text')
	TokenText

	// TokenNumber is a numeric constant (1)
	TokenNumber

	// TokenDuration is a timespan (1d2h3m4s)
	TokenDuration

	// TokenOperator is a maths operator (+-*/%<>)
	TokenOperator
)
