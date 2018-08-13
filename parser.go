package robolang

import (
	"fmt"
)

// Parser converts a stream of input into an Abstract Syntax Tree (AST).
type Parser struct {
	Log func(string, ...interface{})

	functionArgMap map[TokenType]func() (*Node, error)
	result         *ParseResult
	s              *Scanner
	buf            struct {
		token *Token
		n     int
	}
}

var (
	skipTokens = map[TokenType]bool{
		TokenWhitespace: true,
		TokenComment:    true,
	}
)

// NewParser builds a new parser instance.
func NewParser(s string) *Parser {
	p := &Parser{
		s:   NewScanner(s),
		Log: func(string, ...interface{}) {},
	}
	p.functionArgMap = map[TokenType]func() (*Node, error){
		TokenDuration: p.parseConstant,
		TokenNumber:   p.parseConstant,
		TokenResource: p.parseResource,
		TokenText:     p.parseConstant,
		TokenVariable: p.parseVariable,
	}
	return p
}

// Parse will parse the entire stream into an AST
func (p *Parser) Parse() *ParseResult {
	if p.result != nil {
		// The parser can only parse once - any subsequent parses should always return the same result
		return p.result
	}

	p.result = &ParseResult{}
	tok := p.scanNextToken()
	if tok.Type == TokenEOF {
		return p.result.addErrorf("Nothing to parse")
	}

	for ; tok.Type != TokenEOF; tok = p.scanNextToken() {
		if tok.Type != TokenNewLine {
			p.unscan()
			node, err := p.parseItem()
			if node != nil {
				p.result.addNode(node)
			}
			if err != nil {
				return p.result.addError(err)
			}
		}
	}

	return p.result
}

func (p *Parser) clearToNewLine() {
	p.Log("Clearing to newline")
	for tok := p.scanNextToken(); tok.Type != TokenEOF && tok.Type != TokenNewLine; tok = p.scanNextToken() {
	}
	p.unscan()
}

func (p *Parser) makeNode(tok *Token, tokenType NodeType) *Node {
	return &Node{
		Token:    tok,
		Type:     tokenType,
		TypeText: tokenType.String(),
	}
}

func (p *Parser) makeUnexpectedError(tok *Token, expected string) error {
	if expected != "" {
		return fmt.Errorf("Unexpected token '%s', expected %s at line %d, pos %d", tok.Value, expected, tok.LineNum, tok.LinePos)
	}
	return fmt.Errorf("Unexpected token '%s' at line %d, pos %d", tok.Value, tok.LineNum, tok.LinePos)
}

func (p *Parser) parseConstant() (*Node, error) {
	tok := p.scanNextToken()
	p.Log("parsing constant %s", tok.Value)
	return p.makeNode(tok, NodeConstant), nil
}

func (p *Parser) parseFunction() (*Node, error) {
	tok := p.scanNextToken()
	if tok.Type != TokenIdentifier {
		p.clearToNewLine()
		return p.makeNode(tok, NodeInvalid), p.makeUnexpectedError(tok, "")
	}

	node := p.makeNode(tok, NodeFunction)
	p.Log("parsing function %s", tok.Value)
	if err := p.validateNextToken(TokenOpenBracket); err != nil {
		return node, err
	}

	for tok = p.scanNextToken(); tok.Type != TokenCloseBracket; {
		p.unscan()
		arg, err := p.parseFunctionArg()
		if err != nil {
			return node, err
		}

		node.AddArgument(arg)
		tok = p.scanNextToken()
		if tok.Type == TokenComma {
			tok = p.scanNextToken()
		}
	}

	tok = p.scanNextToken()
	if tok.Type != TokenColon {
		p.unscan()
		return node, nil
	}

	if err := p.validateNextToken(TokenNewLine); err != nil {
		return node, err
	}

	cont := true
	tok = p.scan()
	if tok.Type != TokenWhitespace {
		err := p.makeUnexpectedError(tok, "")
		return node, err
	}
	whitespace := tok.Value

	for cont {
		tok = p.scanNextToken()
		if tok.Type != TokenNewLine {
			p.unscan()
			child, err := p.parseItem()
			if err != nil {
				return node, err
			}
			node.AddChild(child)
		}

		for tok = p.scan(); tok.Type == TokenNewLine; tok = p.scan() {
		}
		cont = (tok.Type == TokenWhitespace) && (tok.Value == whitespace)
	}

	p.unscan()
	return node, nil
}

func (p *Parser) parseFunctionArg() (*Node, error) {
	tok := p.scanNextToken()
	if tok.Type != TokenIdentifier {
		p.clearToNewLine()
		return p.makeNode(tok, NodeInvalid), p.makeUnexpectedError(tok, "")
	}

	p.Log("parsing function argument %s", tok.Value)
	node := p.makeNode(tok, NodeArgument)
	if err := p.validateNextToken(TokenEquals); err != nil {
		return node, err
	}

	tok = p.scanNextToken()
	parseFunc, ok := p.functionArgMap[tok.Type]
	if ok {
		p.unscan()
		child, err := parseFunc()
		if err != nil {
			return node, err
		}
		node.AddChild(child)
	} else {
		err := fmt.Errorf("Unable to parse function arg, found '%s' at line %d, pos %d", tok.String(), tok.LineNum, tok.LinePos)
		return node, err
	}
	return node, nil
}

func (p *Parser) parseItem() (*Node, error) {
	tok := p.scanNextToken()
	switch tok.Type {
	case TokenIdentifier:
		p.unscan()
		return p.parseFunction()
	}

	return p.makeNode(tok, NodeInvalid), p.makeUnexpectedError(tok, "")
}

func (p *Parser) parseResource() (*Node, error) {
	tok := p.scanNextToken()
	p.Log("parsing resource %s", tok.Value)
	return p.makeNode(tok, NodeResource), nil
}

func (p *Parser) parseVariable() (*Node, error) {
	tok := p.scanNextToken()
	p.Log("parsing variable %s", tok.Value)
	return p.makeNode(tok, NodeVariable), nil
}

func (p *Parser) scan() *Token {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.token
	}

	token := p.s.Scan()
	p.buf.token = token
	p.result.addToken(token)
	return token
}

func (p *Parser) scanNextToken() *Token {
	token := p.scan()
	for {
		if !skipTokens[token.Type] {
			break
		}
		token = p.scan()
	}

	p.Log("Scanned %s", token.String())
	return token
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) validateNextToken(expected TokenType) error {
	tok := p.scanNextToken()
	return p.validateToken(tok, expected)
}

func (p *Parser) validateToken(tok *Token, expected TokenType) error {
	if tok.Type != expected {
		text := expected.String()
		return p.makeUnexpectedError(tok, text)
	}
	return nil
}

// ParseResult is generated from the parser after.
type ParseResult struct {
	Errors []error
	Nodes  []*Node
	Tokens []*Token
}

// Script converts the parse result to an executable script
func (result *ParseResult) Script() *Script {
	return &Script{
		Nodes: result.Nodes,
	}
}

func (result *ParseResult) addError(err error) *ParseResult {
	result.Errors = append(result.Errors, err)
	return result
}

func (result *ParseResult) addErrorf(format string, a ...interface{}) *ParseResult {
	return result.addError(fmt.Errorf(format, a...))
}

func (result *ParseResult) addNode(node *Node) *ParseResult {
	result.Nodes = append(result.Nodes, node)
	return result
}

func (result *ParseResult) addToken(token *Token) *ParseResult {
	result.Tokens = append(result.Tokens, token)
	return result
}
