package errors

import (
	"fmt"
	"strings"
)

func NewArityError(msg string) Error {
	conf := ErrorConfig{Message: msg}
	return NewError(conf, REFERENCE_ERROR)
}

func RequiresXArgumentsError(expected, given int, fn string) Error {
	msg := fmt.Sprintf("Function '%s' requires %d %s, %d given", fn, expected, arguments(expected), given)
	return NewArityError(msg)
}

func ArgumentToXMustBeYError(arg, fn, expectedType, given string) Error {
	msg := fmt.Sprintf("'%s' argument to '%s' must be %s, '%s' given", arg, fn, expectedType, given)
	return NewArityError(msg)
}

func ArgumentToXAtYMustBeZError(idx int, fn, expectedType, given string) Error {
	msg := fmt.Sprintf("Argument to '%s' at index %d must be %s, %s given", fn, idx, expectedType, given)
	return NewArityError(msg)
}

func RequiresAtLeastXArgumentsError(fn string, given, expected int) Error {
	msg := fmt.Sprintf("'%s' requires at least %d %s, %d given", fn, expected, arguments(expected), given)
	return NewArityError(msg)
}

func RequiresAtMostXArgumentsError(fn string, given, expected int) Error {
	msg := fmt.Sprintf("'%s' requires at most %d %s, %d given", fn, expected, arguments(expected), given)
	return NewArityError(msg)
}

func TypeOfArgumentNotAllowed(fn, arg, argType string, allowed []string) Error {
	msg := fmt.Sprintf("Argument '%s' of type %s to '%s' not allowed. '%s' expects %s", arg, argType, fn, fn, strings.Join(allowed, ", "))
	return NewArityError(msg)
}

func arguments(expected int) string {
	if expected == 1 {
		return "argument"
	}
	return "arguments"
}
