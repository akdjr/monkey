package repl

import (
	"akdjr/monkey/evaluator"
	"akdjr/monkey/lexer"
	"akdjr/monkey/object"
	"akdjr/monkey/parser"
	"bufio"
	"io"
)

// PROMPT is the repl line prompt
const PROMPT = ">>"

// Start reads tokens until it hits EOF
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		io.WriteString(out, PROMPT)
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
			for _, msg := range errors {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		} else {
			// io.WriteString(out, program.String())
			result := evaluator.Eval(program, env)
			if result != nil {
				io.WriteString(out, result.Inspect())
			} else {
				io.WriteString(out, "nil")
			}
			io.WriteString(out, "\n")
		}
	}
}
