package parser

import (
	"akdjr/monkey/ast"
	"akdjr/monkey/lexer"
	"akdjr/monkey/token"
)

// Parser represents an instance of a parser.  It takes a lexer and creates an AST, a tree of statemens and expressions that represents the grammar of the language
type Parser struct {
	l *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
}

// New creates a new Parser from a Lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	// Read two tokens to set currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	// set current token to the peeked token
	// read the next token from the lexer
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram generates and returns an AST of the program
func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
