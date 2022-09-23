package errors

import "fmt"

func NewAssignmentError(msg string) Error {
	conf := ErrorConfig{Message: msg}
	return NewError(conf, ASSIGNMENT_ERROR)
}

func ConstantAssignmentError(id string) Error {
	msg := fmt.Sprintf("Illegal assignment to constant '%s'", id)
	return NewAssignmentError(msg)
}
