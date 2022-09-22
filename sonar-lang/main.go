package main

import (
	"fmt"
	"os"
	"path"

	"github.com/icheka/sonar-lang/sonar-lang/evaluator"
	"github.com/icheka/sonar-lang/sonar-lang/inputs"
	"github.com/icheka/sonar-lang/sonar-lang/lexer"
	"github.com/icheka/sonar-lang/sonar-lang/object"
	"github.com/icheka/sonar-lang/sonar-lang/parser"
	"github.com/icheka/sonar-lang/sonar-lang/repl"
)

func evaluate(source string, options *lexer.LexerOptions) {
	p := parser.New(lexer.New(source, options))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(os.Stderr, p.Errors())
		return
	}
	if evaluated := evaluator.Eval(program, object.NewEnvironment()); evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
	}
}

func main() {
	evaluator.InitStdlib()

	args := os.Args[1:]
	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	if len(args) == 2 {
		switch args[0] {
		case "-f":
			cwd, err := os.Getwd()
			if err != nil {
				panic("Something went wrong!")
			}

			filePath := path.Join(cwd, args[1])
			input := &inputs.FileInput{Path: filePath}
			source := input.Read()

			evaluate(source, &lexer.LexerOptions{Path: filePath})
			return
		case "-text":
			evaluate(args[1], nil)
			return
		}
	}

	fmt.Println("Usage: go run main.go [-f [path]] | [-text input]")
}
