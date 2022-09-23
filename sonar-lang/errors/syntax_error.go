package errors

import (
	"fmt"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/token"
)

func NewSyntaxError(msg string, conf ErrorConfig) Error {
	conf.Message = msg
	return NewError(conf, SYNTAX_ERROR)
}

func PeekError(expected token.TokenType, got string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Illegal or unexpected token. Expected token to be '%s', got '%s'.", expected, got)
	return NewSyntaxError(msg, conf)
}

func IllegalTokenError(t string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Illegal or unexpected token '%s'.", t)
	return NewSyntaxError(msg, conf)
}

func NoPrefixParseFnError(s string, t token.TokenType, conf ErrorConfig) Error {
	conf.Hint = fmt.Sprintf("No prefix parse function for %s found", t)
	return IllegalTokenError(s, conf)
}

func CouldNotParseAsIntegerError(t string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Could not parse %s as integer", t)
	return NewSyntaxError(msg, conf)
}

func CouldNotParseAsFloatError(t string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Could not parse %s as float", t)
	return NewSyntaxError(msg, conf)
}

func ExpectedIdentifierInAssignmentError(t string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Expected identifier in assignment expression, got %s", t)
	return NewSyntaxError(msg, conf)
}

func IdentifierAlreadyDefinedError(id string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Identifier '%s' has already been defined", id)
	return NewSyntaxError(msg, conf)
}

// for arrays and strings
func UnacceptableIndexError(idx string, idxObjectType, objectType string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unacceptable index '%s' for %s. Index must be %s, %s given.", idx, objectType, "INTEGER", idxObjectType)
	return NewSyntaxError(msg, conf)
}

func UnacceptableTypeInKeyAssignmentError(idx string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unacceptable type '%s' in key-assignment operation", idx)
	return NewSyntaxError(msg, conf)
}

func NonIterableInForLoopError(conf ErrorConfig) Error {
	return NewSyntaxError("", conf)
}

func UnknownPrefixOperatorError(operator string, right string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unknown prefix '%s%s'", operator, right)
	return NewSyntaxError(msg, conf)
}

func UnknownPostfixOperatorError(operator string, right string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unknown postfix '%s%s'", operator, right)
	return NewSyntaxError(msg, conf)
}

func TypeMismatchError(operator string, left, right string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Type mismatch: '%s %s %s'", left, operator, right)
	return NewSyntaxError(msg, conf)
}

func UnknownOperatorError(operator string, left, right string, conf ErrorConfig) Error {
	v := fmt.Sprintf("%s %s %s", left, operator, right)
	msg := fmt.Sprintf("Unknown operator: '%s'", strings.TrimSpace(v))
	return NewSyntaxError(msg, conf)
}

func UnusableAsHashKeyError(right string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unusable as hash key. '%s' is not hashable.", right)
	return NewSyntaxError(msg, conf)
}

func UnacceptableRHSInArrayInfixExpressionError(operator string, right string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unacceptable type on right-hand side of array infix expression: '%s %s'", operator, right)
	return NewSyntaxError(msg, conf)
}

func UnacceptableLHSInPostfixExpression(operator, left string, conf ErrorConfig) Error {
	msg := fmt.Sprintf("Unacceptable type on left-hand side of postfix expression: '%s%s'", operator, left)
	return NewSyntaxError(msg, conf)
}
