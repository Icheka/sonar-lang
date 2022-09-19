package evaluator

import (
	"testing"

	"github.com/icheka/sonar-lang/object"
)

func TestMapEntriesBuiltin(t *testing.T) {
	input := `mapEntries({"a": 1})`
	testEvalType[*object.Hash](t, input, "[['a', 1]]")
}

func TestMapValuesBuiltin(t *testing.T) {
	input := `mapValues({"a": 1})`
	testEvalType[*object.Hash](t, input, "[1]")
}

func TestMapKeysBuiltin(t *testing.T) {
	input := `mapKeys({"a": 1})`
	testEvalType[*object.Hash](t, input, "['a']")
}
