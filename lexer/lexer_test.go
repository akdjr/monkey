package lexer

import (
	"akdjr/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct{
		expectedType token.TokenType
		expectedLiteral string
	}
}