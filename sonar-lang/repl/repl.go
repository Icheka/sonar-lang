package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/errors"
	"github.com/icheka/sonar-lang/sonar-lang/evaluator"
	"github.com/icheka/sonar-lang/sonar-lang/keys"
	"github.com/icheka/sonar-lang/sonar-lang/lexer"
	"github.com/icheka/sonar-lang/sonar-lang/object"
	"github.com/icheka/sonar-lang/sonar-lang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line, nil)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			return
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func PrintParserErrors(out io.Writer, errors []errors.Error) {
	if len(errors) == 0 {
		return

	}
	for i, e := range errors[0:1] {
		if len(e.File) != 0 {
			io.WriteString(out, fmt.Sprintf("File %s, ", normalisePath(e.File)))
		}

		// add [LINE:COLUMN]
		io.WriteString(out, fmt.Sprintf("line %d:%d\n\n", e.Line, e.Column))

		// and error
		io.WriteString(out, "\t")
		if len(e.LineText) != 0 {
			drawErrorTracer(out, &e)
		}
		io.WriteString(out, fmt.Sprintf("%s\n", e.String()))

		if len(e.Hint) != 0 && keys.Keys.MODE == "DEV" {
			io.WriteString(out, fmt.Sprintf("[Hint] %s\n", e.Hint))
		}
		if i < len(errors)-1 {
			io.WriteString(out, "\n")
		}
	}
}

func normalisePath(p string) string {
	pwd, _ := os.Getwd()
	return strings.Replace(p, pwd, ".", 1)
}

func drawErrorTracer(out io.Writer, err *errors.Error) {
	if len(err.LineText) == 0 {
		return
	}

	indent := []string{"\t"}
	for i := 0; i < err.LineTextTokenPosition; i++ {
		indent = append(indent, " ")
	}
	indent = append(indent, "^")

	io.WriteString(out, fmt.Sprintf("%s\n%s\n\n", err.LineText, strings.Join(indent, "")))
}
