package parser

import (
	"akdjr/monkey/ast"
	"akdjr/monkey/lexer"
	"akdjr/monkey/token"
	"fmt"
	"strconv"
)

// whenever a token type is found, call the appropriate parsing functions and return an AST node representing it
// a token can have up to two parse functions, a prefix and an infix depending on where its found
// prefix token is found before the expression
// infix token is found between expressions and takes an argument that is the "left" side of the token
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// left binding and right binding powers are the same
const (
	_           int = iota
	LOWEST          // lowest precedence
	EQUALS          // == or !=
	LESSGREATER     // > or <
	SUM             // + or -
	PRODUCT         // * or /
	PREFIX          // -X or !X
	CALL            // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// Parser represents an instance of a parser.  It takes a lexer and creates an AST, a tree of statements and expressions that represents the grammar of the language
type Parser struct {
	l      *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
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

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// register parsing functions
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

// Errors returns all parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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
		stmt := p.parseExpressionStatement()

		if stmt != nil {
			return stmt
		}

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

	// currentToken is =, advance forward and parse the expression
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	// if next token is semicolon, advance forward
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// parse a return statement - return <expression>;
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	// current token is return, advance forward to expression
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// if next token is semicolon, advance to it
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currentToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	// if the next token is semicolon, set it to currentS
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// register a parser error when something unimplemented is encountered
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for '%s' found", t)
	p.errors = append(p.errors, msg)
}
func (p *Parser) noInfixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no infix parse function for '%s' found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// parse an expression
	// check for a prefix parsing function and apply it to get the left side of the expression
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	// peek at the next token and examine its precedence.  if the peeked precedence is higher than the current token, this is an infix expression.
	// we keep advancing forward until we encounter a token with a lower precedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			// no infix parse function, just return the current left side
			return leftExp
		}

		// infix expression, advance so that this peeked token is the current token
		// and then evalute the right hand side
		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.currentToken,
	}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	// current token is the prefix operator, advance to the next token to start the next expression at a higher precedence
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parse an infix expression.  given the left side of the operator, parse the right side
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parse a boolean literal.  true/false
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}

// parse a grouped expression - (<expression>)
func (p *Parser) parseGroupedExpression() ast.Expression {
	// currentToken is '('
	// advance the token to get the beginning of the grouped expression and parse as normal with LOWEST precedence as if we are starting from the beginning
	p.nextToken()

	// start from lowest precedence as if this was a fresh expression
	exp := p.parseExpression(LOWEST)

	// after the expression is fully parsed, expect to find ')' as the next token
	// if not, we have an error (missing a closing ')')
	// expectPeek will advance the tokens if we encounter the correct token
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parse an if expression - if (<condition) { <consequence> }
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.currentToken,
	}

	// next token needs to be a '(' to begin condition expression
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// currentToken is '(', advance to get condition expression
	p.nextToken()

	// parse the condition expression
	expression.Condition = p.parseExpression(LOWEST)

	// next token needs to be a ')' to end the condition expression
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// currentToken is ')'
	// next token needs to be a '{' to begin the consequence block
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// check if the next token is an else as we are currently at the '}' of the consequence block statement
	if p.peekTokenIs(token.ELSE) {
		// an else clause is present of the form else { <alternative }
		// advance to the else token
		p.nextToken()

		// check that the next token is '{' for the alternative block
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// parse a BlockStatement of form { <statements> }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// currentToken is an opening brace, parse statements until we hit a '}'
	block := &ast.BlockStatement{
		Token: p.currentToken,
	}
	block.Statements = []ast.Statement{}

	// advance past the '{'
	p.nextToken()

	// iterate through statements until we hit a '}' or end of file
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		// parseStatement advances until we hit a semicolon
		p.nextToken()
	}

	return block
}

// parse a function literal expression of form fn(<identifier params>) { <blockstatement> }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	function := &ast.FunctionLiteral{
		Token: p.currentToken,
	}

	// next token must be a '('
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// parse all prameters till the closing ')'
	function.Parameters = p.parseFunctionParameters()

	// next token must be the '{' to start the function body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	// empty set of parameters
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		// next token is closing paren, empty parameter list
		p.nextToken()
		return identifiers
	}

	// currentToken is now '(' and the next token needs to be an identifier
	p.nextToken()

	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	identifiers = append(identifiers, ident)

	// if the next token is a comma, iterate through all ident/comma token pairs
	for p.peekTokenIs(token.COMMA) {
		// advance to comma, then advance to token after comma
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	// expect next token to be an RPAREN and advance to it
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parse a call expression.  a call expression is essential an infix expression with the '(' as an operator.  ex. add(5) - left = identifier (add or fn for function literal), operator = '(', right = expression arguments
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.currentToken,
		Function: function,
	}
	exp.Arguments = p.parseCallArguments()

	return exp
}

// parse call arguments.
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// currenttoken is '('
	if p.peekTokenIs(token.RPAREN) {
		// next token is ), advance past it and return empty args
		p.nextToken()
		return args
	}

	// advance to first expression token and parse it
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	// as long as the next token is a comma, keep parsing expressions
	for p.peekTokenIs(token.COMMA) {
		// advance to comma
		p.nextToken()
		// advance to expression
		p.nextToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	// at this point, 1 or more args, require closing paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// check if current token is token type
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// check if peek token is token type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// look at the peekToken.  if it is of the right type, advance the tokens and get the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// get the precedence of the next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// get the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}
