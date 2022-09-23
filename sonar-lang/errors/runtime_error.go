package errors

import "fmt"

func NewRuntimeError(msg string) Error {
	conf := ErrorConfig{Message: msg}
	return NewError(conf, RUNTIME_ERROR)
}

func IllegalConversionError(from, to string) Error {
	msg := fmt.Sprintf("Illegal conversion: %s -> %s", from, to)
	return NewRuntimeError(msg)
}

func TypeCannotBeCopiedError(t string) Error {
	msg := fmt.Sprintf("Type '%s' cannot be copied", t)
	return NewRuntimeError(msg)
}

func ZeroDivisionError(t string) Error {
	msg := fmt.Sprintf("Division by zero (%s/0)", t)
	return NewRuntimeError(msg)
}
