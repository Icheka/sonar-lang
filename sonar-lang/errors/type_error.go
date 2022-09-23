package errors

import (
	"fmt"
)

func NewTypeError(msg string, conf *ErrorConfig) Error {
	if conf != nil {
		conf.Message = msg
		return NewError(*conf, REFERENCE_ERROR)
	}
	return NewError(ErrorConfig{Message: msg}, REFERENCE_ERROR)
}

func TypeError(id string, t string) Error {
	msg := fmt.Sprintf("'%s' is not of type %s", id, t)
	return NewTypeError(msg, nil)
}
