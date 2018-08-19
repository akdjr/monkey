package token

// token type - change to int?
type TokenType string

// A token
// Type is the type of token
// Literal is the literal value of the token (such as the name of the identifier or the value of a literal)
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Illegal token and End of File token
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifers and literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	// Operators
	ASSIGN = "ASSIGN"
	PLUS   = "+"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
