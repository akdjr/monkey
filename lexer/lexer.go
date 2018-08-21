package lexer

import (
	"akdjr/monkey/token"
)

// Lexer represents an instance of the lexer.  It holds the input source as a string and maintains the current position in the input which points to the current char and the position of the next character to read
type Lexer struct {
	input        string
	position     int  // points to where we are currently reading
	readPosition int  // points to where we will read next
	currentChar  byte // character at the position where we are currently reading
}

// New creates a new instance of the lexer from an input string.  input contains monkey source
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	l.readChar()

	return l
}

// Load loads a new set of input source into the lexer and resets it
func (l *Lexer) Load(input string) {
	l.input = input
	l.position = 0
	l.readPosition = 0
	l.currentChar = 0

	l.readChar()
}

// NextToken converts the current char into a token
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	// currentChar is ASCII byte
	// TODO: Full unicode support instead of simple ASCII
	// TODO: Add additional two character operators and refactor
	switch l.currentChar {
	case '=':
		// peek ahead for equal operator
		if l.peekChar() == '=' {
			// save the current =, read the next = and create the EQ token
			ch := l.currentChar
			l.readChar()
			literal := string(ch) + string(l.currentChar)
			t = token.Token{
				Type:    token.EQ,
				Literal: literal,
			}
		} else {
			t = token.New(token.ASSIGN, l.currentChar)
		}
	case '+':
		t = token.New(token.PLUS, l.currentChar)
	case '-':
		t = token.New(token.MINUS, l.currentChar)
	case '!':
		// peek ahead for not equal operator
		if l.peekChar() == '=' {
			// save the current =, read the next = and create the EQ token
			ch := l.currentChar
			l.readChar()
			literal := string(ch) + string(l.currentChar)
			t = token.Token{
				Type:    token.NOT_EQ,
				Literal: literal,
			}
		} else {
			t = token.New(token.BANG, l.currentChar)
		}
	case '/':
		t = token.New(token.SLASH, l.currentChar)
	case '*':
		t = token.New(token.ASTERISK, l.currentChar)
	case '<':
		t = token.New(token.LT, l.currentChar)
	case '>':
		t = token.New(token.GT, l.currentChar)
	case ';':
		t = token.New(token.SEMICOLON, l.currentChar)
	case '(':
		t = token.New(token.LPAREN, l.currentChar)
	case ')':
		t = token.New(token.RPAREN, l.currentChar)
	case ',':
		t = token.New(token.COMMA, l.currentChar)
	case '{':
		t = token.New(token.LBRACE, l.currentChar)
	case '}':
		t = token.New(token.RBRACE, l.currentChar)
	case 0:
		t = token.Token{
			Type:    token.EOF,
			Literal: "",
		}

	default:
		// if the current character is a letter, we need to read characters until we hit a non-letter
		if isLetter(l.currentChar) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t
		} else if isDigit(l.currentChar) {
			t.Literal = l.readNumber()
			t.Type = token.INT
			return t
		} else {
			// current character is not a letter or a valid single character token
			t = token.New(token.ILLEGAL, l.currentChar)
		}
	}

	// read the next character into the lexer
	l.readChar()
	return t
}

// skipWhitespace will advance the position past all whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

// read the next character in input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		//  readPosition has reached the end of the input string
		// set the current char to ASCII 0 (NUL), this will be returned as token.EOF
		l.currentChar = 0
	} else {
		// read the character from the input and assign it to currentChar for the lexer to examine
		l.currentChar = l.input[l.readPosition]
	}

	// set position to be the location of this current character
	l.position = l.readPosition

	// advance the read position by 1 character
	l.readPosition++
}

// peek at the next char in the input and return it, but do not advance the read position
// readPosition always points to the next character after the current one
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// readIdentifier reads characters until it hits a non-letter and then returns the string
func (l *Lexer) readIdentifier() string {
	// start at the current character position
	position := l.position

	// keep advancing forward until we hit a non-letter
	for isLetter(l.currentChar) {
		l.readChar()
	}

	// l.currentCharacter and l.position represent the character after the identifier
	// return the slice for this identifier
	return l.input[position:l.position]
}

// readNumber reads characters until it hits a non-digit
// TODO: allow floating point numbers as this only supports integers
// TODO: almost identical to readIdentifier - REFACTOR
func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.currentChar) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// isLetter checks if ch is a valid identifier character
// TODO: Full unicode support instead of simple ASCII
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || (ch == '_')
}

// isDigit checks if ch is a valid digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
