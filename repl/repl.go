package repl

import (
	"akdjr/monkey/lexer"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
