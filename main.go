package main

import (
	"os"
	"sonar/v2/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
