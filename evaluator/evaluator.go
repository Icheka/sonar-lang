package evaluator

import (
	"fmt"
	"sonar/v2/ast"
	"sonar/v2/object"
	"sonar/v2/token"
	"strconv"
	"strings"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if _, ok := env.Get(node.Name.Value); ok {
			return &object.Error{
				Message: fmt.Sprintf("SyntaxError: Identifier '%s' has already been declared", node.Name.Value),
			}
		}
		env.Set(node.Name.Value, val)

	case *ast.AssignmentExpression:
		right := Eval(node.Value, env)
		if isError(right) {
			return right
		}
		oldValue, ok := env.Get(node.Identifier.Value)
		if !ok {
			return &object.Error{
				Message: fmt.Sprintf("SyntaxError: Assignment to undefined identifier '%s'", node.Identifier.Value),
			}
		}

		result := right

		if node.Operator != token.ASSIGN {
			operators := map[token.TokenType]token.TokenType{
				(token.PLUS_ASSIGN):     token.PLUS,
				(token.MINUS_ASSIGN):    token.MINUS,
				(token.ASTERISK_ASSIGN): token.ASTERISK,
				(token.SLASH_ASSIGN):    token.SLASH,
			}

			// re-use evalInfixExpression...
			// ... for example, if node.Operator is token.PLUS_ASSIGN, this will evaluate oldValue + right and return to result
			result = evalInfixExpression(string(operators[token.TokenType(node.Operator)]), oldValue, right)

			// catch errors from evalInfixhere
			if isError(result) {
				return result
			}
		}
		return env.Set(node.Identifier.Value, result)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(env, node)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	// evaluate condition
	// -- if condition throws error, return error
	// -- if condition is not true, break loop
	// -- evaluate consequence

	var result object.Object

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		result := Eval(ws.Consequence, env)
		// if there is a return statement, return immediately
		if !isError(result) && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch operator {
	case token.AND:
		l := isTruthy(left)
		r := isTruthy(right)
		return nativeBoolToBooleanObject(l && r)
	case token.OR:
		l := isTruthy(left)
		r := isTruthy(right)
		return nativeBoolToBooleanObject(l || r)
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

		// if types are mismatched, cast the integer type to float
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		right = &object.Float{
			Value: float64(right.(*object.Integer).Value),
		}
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		left = &object.Float{
			Value: float64(left.(*object.Integer).Value),
		}
		return evalFloatInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)

	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}

	value := right.(*object.Float).Value
	return &object.Float{Value: -value}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.SLASH:
		return evalNumberDivision(leftVal, rightVal)
	case token.LT:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case token.GT:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case token.EQ:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case token.LTE:
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case token.GTE:
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case token.PLUS:
		return &object.Float{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Float{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Float{Value: leftVal * rightVal}
	case token.SLASH:
		return evalNumberDivision(leftVal, rightVal)
	case token.LT:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case token.GT:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case token.EQ:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case token.LTE:
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case token.GTE:
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalNumberDivision[T int64 | float64](leftVal, rightVal T) object.Object {
	result := float64(leftVal) / float64(rightVal)
	if strings.Contains(fmt.Sprintf("%f", result), ".") {
		return &object.Float{Value: result}
	}
	return &object.Integer{Value: int64(result)}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case token.PLUS:
		return &object.String{Value: leftVal + rightVal}

	case token.MINUS:
		return &object.String{Value: strings.ReplaceAll(leftVal, rightVal, "")}

	case token.EQ:
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case token.NOT_EQ:
		return nativeBoolToBooleanObject(leftVal != rightVal)

	case token.LT:
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case token.GT:
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case token.LTE:
		return nativeBoolToBooleanObject(leftVal <= rightVal)

	case token.GTE:
		return nativeBoolToBooleanObject(leftVal >= rightVal)

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringindexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalStringindexExpression(str, index object.Object) object.Object {
	/*
		- cast to *object.String
		- get literal of index
		- test index literal is within bounds such that 0 <= index <= len(str)
		- return value at index
	*/
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value))

	if max == 0 || idx < 0 || idx > max {
		return NULL
	}

	return &object.String{
		Value: string(strObject.Value[idx]),
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalPostfixExpression(env *object.Environment, node *ast.PostfixExpression) object.Object {
	/*
		- get token from env
		- if no token, check if token is int literal
		- if token is integer, set incremented value in env
	*/
	tokenLiteral := node.Token.Literal
	tok, ok := env.Get(tokenLiteral)
	if !ok {
		if intValue, ok := strconv.ParseInt(tokenLiteral, 10, 64); ok == nil {
			tok = &object.Integer{Value: intValue}
		} else {
			return newError("Unknown token %s", tokenLiteral)
		}
	}

	switch node.Operator {
	case token.POST_INCR:
		integer, ok := tok.(*object.Integer)
		if !ok {
			return newError("Left side of post-increment operator must be of type int, got %s", tokenLiteral)
		}

		newInteger := &object.Integer{Value: integer.Value + 1}
		env.Set(tokenLiteral, newInteger)
		return newInteger

	case token.POST_DECR:
		integer, ok := tok.(*object.Integer)
		if !ok {
			return newError("Left side of post-decrement operator must be of type int, got %s", tokenLiteral)
		}

		newInteger := &object.Integer{Value: integer.Value - 1}
		env.Set(tokenLiteral, newInteger)
		return newInteger
	}
	return newError("Unknown operator %s", node.Operator)
}
