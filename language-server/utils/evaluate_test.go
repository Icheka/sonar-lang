package utils

import (
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	input := `print(1)`
	stdOut, stdErr := Evaluate(input)

	if l := len(strings.Trim(stdErr, " ")); l != 0 {
		t.Errorf("expected length of stderr to be 0, got=%d", l)
	}

	if l := len(strings.Trim(stdOut, " ")); l == 0 {
		t.Error("expected length of stdout to be > 0, got=0")
	}
}
