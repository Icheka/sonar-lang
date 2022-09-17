package main

import (
	"fmt"
	"os"
	"path"
	"sonar/v2/evaluator"
	"sonar/v2/inputs"
	"sonar/v2/lexer"
	"sonar/v2/object"
	"sonar/v2/parser"
	"sonar/v2/repl"
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
				for _, err := range p.Errors() {
					fmt.Printf("%s\n", err)
				}
			}
			if evaluated := evaluator.Eval(program, object.NewEnvironment()); evaluated.Type() == object.ERROR_OBJ {
				fmt.Println(evaluated.Inspect())
			}
			return
		}
	}

	fmt.Println("Usage: go run main.go [-f [path]]")
}
