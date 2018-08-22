package ast

import (
	"akdjr/monkey/token"
)

// Node represents a node in the AST.  TokenLiteral() returns the literal value of the token that the node is associated with
type Node interface {
	TokenLiteral() string
}

// Statement represents a statment in the AST.  Statements do not produce values
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the AST.  Expressions produce a value
type Expression interface {
	Node
	expressionNode()
}

// Program is a node that will be the root of the AST.  It is represented as a series of statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// LetStatement represeents a let statement of the format "let <identifier> = <expression>"
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier represents an expression that holds an identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
