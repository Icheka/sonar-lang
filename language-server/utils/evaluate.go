package utils

import (
	"bytes"
	"log"
	"os/exec"
)

func Evaluate(input string) (string, string) {
	cmd := exec.Command("sonar-lang", "-text", input)

	var outBuf, errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	cmd.Stdout = &outBuf

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	return outBuf.String(), errBuf.String()
}
