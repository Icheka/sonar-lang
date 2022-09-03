package repl

import (
	"bufio"
	"fmt"
	"io"
	"sonar/v2/lexer"
	"sonar/v2/token"
)

func Start(in io.Reader, out io.Writer) {
	prompt := ">> "
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(prompt)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
