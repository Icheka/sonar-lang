package errors

import (
	"fmt"
)

func NewReferenceError(msg string, conf ErrorConfig) Error {
	conf.Message = msg
	return NewError(conf, REFERENCE_ERROR)
}

func IdentifierNotDefinedError(id string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Identifier '%s' has not been defined", id)
	return NewReferenceError(msg, conf)
}

func IndexOperatorNotAllowed(t string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Index operator not supported for %s", t)
	return NewReferenceError(msg, conf)
}

func OutOfRangeError(idx int, max int, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Index '%d' out of range [%d]", idx, max)
	return NewReferenceError(msg, conf)
}

func InvalidRangeError(start, end int, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Invalid range %d:%d", start, end)
	return NewReferenceError(msg, conf)
}
