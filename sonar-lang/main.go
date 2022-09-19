package main

import (
	"fmt"
	"os"
	"path"

	"github.com/icheka/sonar-lang/evaluator"
	"github.com/icheka/sonar-lang/inputs"
	"github.com/icheka/sonar-lang/lexer"
	"github.com/icheka/sonar-lang/object"
	"github.com/icheka/sonar-lang/parser"
	"github.com/icheka/sonar-lang/repl"
)

func main() {
	evaluator.InitStdlib()

	args := os.Args[1:]
	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	if len(args) == 2 {
		if args[0] == "-f" {
			cwd, err := os.Getwd()
			if err != nil {
				panic("Something went wrong!")
			}
			filePath := path.Join(cwd, args[1])

			input := &inputs.FileInput{Path: filePath}
			source := input.Read()

			p := parser.New(lexer.New(source))
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				repl.PrintParserErrors(os.Stderr, p.Errors())
				return
			}
			if evaluated := evaluator.Eval(program, object.NewEnvironment()); evaluated.Type() == object.ERROR_OBJ {
				fmt.Println(evaluated.Inspect())
			}
			return
		}
	}

	fmt.Println("Usage: go run main.go [-f [path]]")
}
