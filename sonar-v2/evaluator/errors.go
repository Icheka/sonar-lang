package evaluator

import (
	"sonar/v2/object"
)

func WrongArityError(got, expected int) *object.Error {
	return NewError("wrong number of arguments. got=%d, want=%d", got, expected)
}
