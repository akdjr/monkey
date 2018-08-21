package lexer

import (
	"akdjr/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	type expectedTokenType struct {
		expectedType    token.TokenType
		expectedLiteral string
	}

	type inputTest struct {
		input          string
		expectedTokens []expectedTokenType
	}

	tests := []inputTest{
		{
			input: `=+(){},;`,
			expectedTokens: []expectedTokenType{
				{token.ASSIGN, "="},
				{token.PLUS, "+"},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.COMMA, ","},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `let five = 5;
			let ten = 10;
			
			let add = fn(x, y) {
				x + y;
			};
			
			let result = add(five, ten);
			`,
			expectedTokens: []expectedTokenType{
				{token.LET, "let"},
				{token.IDENTIFIER, "five"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},

				{token.LET, "let"},
				{token.IDENTIFIER, "ten"},
				{token.ASSIGN, "="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},

				{token.LET, "let"},
				{token.IDENTIFIER, "add"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENTIFIER, "x"},
				{token.COMMA, ","},
				{token.IDENTIFIER, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},

				{token.IDENTIFIER, "x"},
				{token.PLUS, "+"},
				{token.IDENTIFIER, "y"},
				{token.SEMICOLON, ";"},

				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},

				{token.LET, "let"},
				{token.IDENTIFIER, "result"},
				{token.ASSIGN, "="},
				{token.IDENTIFIER, "add"},
				{token.LPAREN, "("},
				{token.IDENTIFIER, "five"},
				{token.COMMA, ","},
				{token.IDENTIFIER, "ten"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},

				{token.EOF, ""},
			},
		},
		{
			input: `!-/*5;
			5 < 10 > 5;
			=`,
			expectedTokens: []expectedTokenType{
				{token.BANG, "!"},
				{token.MINUS, "-"},
				{token.SLASH, "/"},
				{token.ASTERISK, "*"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},

				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.GT, ">"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.ASSIGN, "="},
				{token.EOF, ""},
			},
		},
		{
			input: `if (5 < 10) {
				return true;
			} else {
				return false;
			}
			`,
			expectedTokens: []expectedTokenType{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},

				{token.RETURN, "return"},
				{token.TRUE, "true"},
				{token.SEMICOLON, ";"},

				{token.RBRACE, "}"},
				{token.ELSE, "else"},
				{token.LBRACE, "{"},

				{token.RETURN, "return"},
				{token.FALSE, "false"},
				{token.SEMICOLON, ";"},

				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `10 == 10;
			10 != 9;
			`,
			expectedTokens: []expectedTokenType{
				{token.INT, "10"},
				{token.EQ, "=="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},

				{token.INT, "10"},
				{token.NOT_EQ, "!="},
				{token.INT, "9"},
				{token.SEMICOLON, ";"},

				{token.EOF, ""},
			},
		},
		{
			input: `?%$#@^\|~'"
			`,
			expectedTokens: []expectedTokenType{
				{token.ILLEGAL, "?"},
				{token.ILLEGAL, "%"},
				{token.ILLEGAL, "$"},
				{token.ILLEGAL, "#"},
				{token.ILLEGAL, "@"},
				{token.ILLEGAL, "^"},
				{token.ILLEGAL, "\\"},
				{token.ILLEGAL, "|"},
				{token.ILLEGAL, "~"},
				{token.ILLEGAL, "'"},
				{token.ILLEGAL, "\""},
				{token.EOF, ""},
			},
		},
	}

	lexer := New("")

	for i, tt := range tests {
		lexer.Load(tt.input)

		for j, tok := range tt.expectedTokens {
			nextToken := lexer.NextToken()

			if nextToken.Type != tok.expectedType {
				t.Fatalf("input[%d] test[%d] - tokentype wrong. expected=%q, got=%q", i, j, tok.expectedType, nextToken.Type)
			}

			if nextToken.Literal != tok.expectedLiteral {
				t.Fatalf("input[%d] test[%d] - literal wrong. expected=%q, got=%q", i, j, tok.expectedLiteral, nextToken.Literal)
			}
		}
	}
}
