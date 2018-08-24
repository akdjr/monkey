package repl

import (
	"akdjr/monkey/lexer"
	"akdjr/monkey/parser"
	"akdjr/monkey/token"
	"bufio"
	"fmt"
	"io"
)

// PROMPT is the repl line prompt
const PROMPT = ">>"

// Start reads tokens until it hits EOF
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		errors := p.Errors()

		if len(errors) > 0 {
			for _, error := range errors {
				fmt.Println(error)
			}
		} else {
			for _, stmt := range program.Statements {
				fmt.Println(stmt.String())
			}
		}

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
