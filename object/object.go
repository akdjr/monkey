package object

import (
	"akdjr/monkey/ast"
	"bytes"
	"fmt"
	"strings"
)

// ObjectType represents the type of object
// TODO: could replace type with something else?
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
)

// Object is an internal representation of a value.  Every value will be wrapped in a struct that fulfills this interface
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer literal value that is the result of evaluating an interger literal in source
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Boolean represents a boolean literal
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Null represents no value.  fun with nulls!
type Null struct {
}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

// ReturnValue represents a return value.  This wraps an object that is the result of a return <expression>
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Error represents an internal error.  These are any internal or user errors that spawn as a result of invalid operators, unsupported operations, or anything else.
// TODO: add stack trace, line/column numbers (LEXER support required), etc.
type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Function represents a function expression.  We keep track of the parameters, the body, and the environment that it is invoked with
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
