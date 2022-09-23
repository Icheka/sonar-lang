package evaluator

import (
	"github.com/icheka/sonar-lang/sonar-lang/errors"
	"github.com/icheka/sonar-lang/sonar-lang/object"
)

func WrongArityError(got, expected int, fn string) *object.Error {
	return NewError(errors.RequiresXArgumentsError(expected, got, fn))
}
