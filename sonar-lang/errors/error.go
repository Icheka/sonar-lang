package errors

import "fmt"

const (
	SYNTAX_ERROR     = "SyntaxError"
	RUNTIME_ERROR    = "RuntimeError"
	REFERENCE_ERROR  = "ReferenceError"
	TYPE_ERROR       = "TypeError"
	ARITY_ERROR      = "ArityError"
	ASSIGNMENT_ERROR = "AssignmentError"
)

type ErrorType string

type Error struct {
	File                  string
	Line                  int
	Column                int
	Message               string
	Type                  ErrorType
	Hint                  string
	LineText              string
	LineTextTokenPosition int
}

func (e *Error) String() string { return fmt.Sprintf("%s: %s", e.Type, e.Message) }

type ErrorConfig struct {
	File                  string
	Line                  int
	Column                int
	Message               string
	Hint                  string
	LineText              string
	LineTextTokenPosition int
}

func NewError(conf ErrorConfig, t ErrorType) Error {
	return Error{
		File:                  conf.File,
		Line:                  conf.Line,
		Column:                conf.Column,
		Type:                  t,
		Message:               conf.Message,
		Hint:                  conf.Hint,
		LineText:              conf.LineText,
		LineTextTokenPosition: conf.LineTextTokenPosition,
	}
}
