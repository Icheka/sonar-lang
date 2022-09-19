package evaluator

import (
	"github.com/icheka/sonar-lang/object"
)

func WrongArityError(got, expected int) *object.Error {
	return NewError("wrong number of arguments. got=%d, want=%d", got, expected)
}
