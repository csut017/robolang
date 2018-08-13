package robolang

import (
	"bufio"
	"bytes"
)

// Scanner is used for splitting an input stream into tokens.
type Scanner struct {
	lineNum int
	linePos int
	r       *bufio.Reader
}

// NewScanner starts a new scanner.
func NewScanner(s string) *Scanner {
	r := bytes.NewBufferString(s)
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan reads the next token from the input stream.
func (s *Scanner) Scan() *Token {
	ch := s.read()

	// Check the multi-character tokens
	if s.isWhitespace(ch) {
		s.unread()
		return s.makeToken(TokenWhitespace, s.scanWhitespace())
	} else if ch == '@' {
		return s.scanIdentifier(TokenResource)
	} else if ch == '&' {
		return s.scanIdentifier(TokenVariable)
	} else if s.isLetter(ch) {
		s.unread()
		return s.scanIdentifier(TokenIdentifier)
	} else if s.isDigit(ch) {
		s.unread()
		return s.scanNumber()
	} else if ch == '\'' {
		return s.scanText()
	} else if ch == '#' {
		return s.scanComment()
	}

	// Check the single character tokens
	t, ok := chars[ch]
	if ok {
		tok := s.makeToken(t, string(ch))
		if t == TokenNewLine {
			s.linePos = 0
			s.lineNum++
		}
		return tok
	}

	return s.makeToken(TokenIllegal, string(ch))
}

var (
	eof   = rune(0)
	chars = map[rune]TokenType{
		rune(0): TokenEOF,
		'\n':    TokenNewLine,
		'(':     TokenOpenBracket,
		')':     TokenCloseBracket,
		',':     TokenComma,
		':':     TokenColon,
		'=':     TokenEquals,
		'+':     TokenOperator,
		'-':     TokenOperator,
		'*':     TokenOperator,
		'/':     TokenOperator,
		'%':     TokenOperator,
		'<':     TokenOperator,
		'>':     TokenOperator,
	}
	whitespace = map[rune]bool{
		' ':  true,
		'\t': true,
		'\r': true,
	}
)

func (s *Scanner) isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch == '.')
}

func (s *Scanner) isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (s *Scanner) isWhitespace(ch rune) bool {
	return whitespace[ch]
}

func (s *Scanner) makeToken(t TokenType, v string) *Token {
	return &Token{
		Type:     t,
		TypeName: t.String(),
		Value:    v,
		LineNum:  s.lineNum,
		LinePos:  s.linePos - len(v),
	}
}

func (s *Scanner) read() rune {
	s.linePos++
	ch, _, err := s.r.ReadRune()

	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) scanComment() *Token {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '\n' {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return s.makeToken(TokenComment, buf.String())
}

func (s *Scanner) scanIdent() string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !s.isLetter(ch) && !s.isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func (s *Scanner) scanIdentifier(typ TokenType) *Token {
	value := s.scanIdent()
	return s.makeToken(typ, value)
}

func (s *Scanner) scanNumber() *Token {
	var buf bytes.Buffer
	elements, numbers, isDuration := map[rune]bool{
		's': true,
		'm': true,
		'h': true,
		'd': true,
	}, map[rune]bool{
		'.': true,
	}, false
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !s.isDigit(ch) && !numbers[ch] {
			if elements[ch] {
				isDuration = true
				elements[ch] = false
				buf.WriteRune(ch)
			} else {
				s.unread()
				break
			}
		} else {
			if numbers[ch] {
				numbers[ch] = false
			}
			buf.WriteRune(ch)
		}
	}
	if isDuration {
		return s.makeToken(TokenDuration, buf.String())
	}
	return s.makeToken(TokenNumber, buf.String())
}

func (s *Scanner) scanText() *Token {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '\'' {
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return s.makeToken(TokenText, buf.String())
}

func (s *Scanner) scanWhitespace() string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !s.isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func (s *Scanner) unread() {
	s.linePos--
	_ = s.r.UnreadRune()
}
