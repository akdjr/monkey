package token

// TokenType holdes the type of the token.  String for now
// TODO: could change to int
type TokenType string

// Token represents a single token in the source
// Type is the type of token
// Literal is the literal value of the token (such as the name of the identifier or the value of a literal)
// TODO: Add some more information, such as line number, column number for enhanced error reporting
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// TODO: could change to int
const (
	// Illegal token and End of File token
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifers and literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	// Operators
	ASSIGN   = "ASSIGN"
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// New creates a new Token from a byte
// TODO: Full unicode support instead of simple ASCII
func New(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// LookupIdentifier checks to see if the input identifier is a reserved keyword
func LookupIdentifier(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		// ident is a reserved keyword
		return tokenType
	}

	// ident is NOT a keyword, return IDENTIFIER
	return IDENTIFIER
}
