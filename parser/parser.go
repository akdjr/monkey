package parser

import (
	"akdjr/monkey/ast"
	"akdjr/monkey/lexer"
	"akdjr/monkey/token"
	"fmt"
)

// Parser represents an instance of a parser.  It takes a lexer and creates an AST, a tree of statemens and expressions that represents the grammar of the language
type Parser struct {
	l      *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token
}

// New creates a new Parser from a Lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens to set currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns all parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	// set current token to the peeked token
	// read the next token from the lexer
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram generates and returns an AST of the program
func (p *Parser) ParseProgram() *ast.Program {
	// initialize the root node of the AST (a *ast.Program)
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// iterate through all tokens until we hit EOF
	for p.currentToken.Type != token.EOF {
		// parse a statement
		// parseStatement will advance tokens as needed
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// parse statements based on the type of the current token
	switch p.currentToken.Type {
	case token.LET:
		stmt := p.parseLetStatement()

		if stmt != nil {
			return stmt
		}

		return nil
	case token.RETURN:
		stmt := p.parseReturnStatement()

		if stmt != nil {
			return stmt
		}

		return nil
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// parse a let statement - let <identifier> = <expression>;
	stmt := &ast.LetStatement{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skip expressions until we encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// parse a return statement - return <expression>;
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	// TODO: skip expressions until we encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// check if current token is token type
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// check if peek token is token type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// look at the peekToken.  if it is of the right type, get the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
