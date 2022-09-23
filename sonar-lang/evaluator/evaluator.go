package evaluator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/ast"
	"github.com/icheka/sonar-lang/sonar-lang/errors"
	"github.com/icheka/sonar-lang/sonar-lang/object"
	"github.com/icheka/sonar-lang/sonar-lang/token"
	"github.com/icheka/sonar-lang/sonar-lang/utils"
)

var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
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
		if node.ReturnValue == nil {
			node.ReturnValue = &ast.NullValueExpression{}
		}
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
		if _, ok := env.Store[node.Name.Value]; ok {
			return NewError(errors.IdentifierAlreadyDefinedError(node.Name.Value, node.TokenInfo.(errors.ErrorConfig)))
		}
		return env.Set(node.Name.Value, val)

	case *ast.AssignmentExpression:
		right := Eval(node.Value, env)
		if isError(right) {
			return right
		}
		oldValue, ok := env.Get(node.Identifier.Value)
		if !ok {
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.IdentifierNotDefinedError(node.Identifier.Value, r))
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
			t := ast.InfixExpression{TokenInfo: node.TokenInfo}
			result = evalInfixExpression(string(operators[token.TokenType(node.Operator)]), oldValue, right, t)

			// catch errors from evalInfixhere
			if isError(result) {
				return result
			}
		}
		return env.Set(node.Identifier.Value, result)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.BreakStatement:
		return BREAK

	case *ast.ContinueStatement:
		return CONTINUE

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
		return evalPrefixExpression(node.Operator, right, *node)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, *node)

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
		return evalIndexExpression(left, index, node)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.SquareBracketAssignment:
		left := Eval(node.Left, env)
		index := Eval(node.Key, env)
		value := Eval(node.Value, env)

		if isError(left) {
			return left
		}
		if isError(index) {
			return index
		}
		if isError(value) {
			return value
		}

		switch left.Type() {
		case object.ARRAY_OBJ:
			if i, ok := index.(*object.Integer); !ok {
				if r, ok := node.TokenInfo.(errors.ErrorConfig); ok {
					return NewError(errors.UnacceptableIndexError(i.Inspect(), string(i.Type()), object.ARRAY_OBJ, r))
				}
			}
			return evalArraySquareBracketExpression(&left, index, value, node)

		case object.HASH_OBJ:
			return evalMapSquareBracketExpression(&left, index, value, node)

		default:
			if r, ok := node.TokenInfo.(errors.ErrorConfig); ok {
				return NewError(errors.UnacceptableTypeInKeyAssignmentError(string(left.Type()), r))
			}
		}

	case *ast.NullValueExpression:
		return &object.Null{}
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
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
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

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		result := Eval(ws.Consequence, env)
		// if there is a return, error or break, return immediately
		if isError(result) || isBreak(result) || result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return &object.Null{}
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	// iterable is an object.Object that implements the Iterable interface
	iterable, ok := Eval(fs.Iterable, env).(object.Iterable)

	if !ok {
		if r, ok := fs.TokenInfo.(errors.ErrorConfig); ok {
			return NewError(errors.NonIterableInForLoopError(r))
		}
	}

	// iterable.Iters
	iters := iterable.Iters()

	var readonly = make(map[string]bool)
	allowed := []string{
		fs.Counter.String(),
	}
	scope := object.NewEphemeralScope(allowed, readonly, env)

	for i, v := range iters {
		// allow the counter to be mutated only to allow setting counter to iter
		scope.Readonly = make(map[string]bool)

		// set 'counter' to current iter
		if iterable.(object.Object).Type() == object.HASH_OBJ {
			arr := v.(*object.Array).Elements
			val := scope.Set(fs.Counter.String(), &object.String{Value: arr[0].Inspect()})
			if isError(val) {
				return val
			}
			val = scope.Set(fs.Value.String(), &object.String{Value: arr[1].Inspect()})
			if isError(val) {
				return val
			}
		} else {
			val := scope.Set(fs.Counter.String(), &object.Integer{Value: int64(i)})
			if isError(val) {
				return val
			}
			val = scope.Set(fs.Value.String(), v)
			if isError(val) {
				return val
			}
		}

		// make the counter a constant to make it immutable until this iteration concludes
		scope.Readonly = map[string]bool{
			(fs.Counter.String()): true,
			(fs.Value.String()):   true,
		}

		result := Eval(fs.Consequence, scope)

		if isContinue(result) {
			continue
		}

		if isError(result) || isBreak(result) || result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}

	}

	return &object.Null{}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object, node ast.PrefixExpression) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, &node)
	default:
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownPrefixOperatorError(operator, string(right.Type()), r))
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object, node ast.InfixExpression,
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
		return evalIntegerInfixExpression(operator, left, right, &node)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, &node)

	// if number types are mismatched, cast the integer type to float
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		right = &object.Float{
			Value: float64(right.(*object.Integer).Value),
		}
		return evalFloatInfixExpression(operator, left, right, &node)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		left = &object.Float{
			Value: float64(left.(*object.Integer).Value),
		}
		return evalFloatInfixExpression(operator, left, right, &node)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, &node)

	case left.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(operator, left, right, &node)

	case left.Type() == object.HASH_OBJ:
		return evalMapInfixExpression(operator, left, right, &node)

	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)

	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.TypeMismatchError(operator, string(left.Type()), string(right.Type()), r))

	default:
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
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

func evalMinusPrefixOperatorExpression(right object.Object, node *ast.PrefixExpression) object.Object {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownOperatorError("-", "", string(right.Type()), r))
	}

	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}

	value := right.(*object.Float).Value
	return &object.Float{Value: -value}
}

func evalZeroDivision[T int64 | float64](left T) object.Object {
	return NewError(errors.ZeroDivisionError(fmt.Sprint(left)))
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object, node *ast.InfixExpression,
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
		if rightVal == 0 {
			return evalZeroDivision(leftVal)
		}
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
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
	}
}

func evalFloatInfixExpression(
	operator string,
	left, right object.Object, node *ast.InfixExpression,
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
		if rightVal == 0 {
			return evalZeroDivision(leftVal)
		}
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
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
	}
}

func evalNumberDivision[T int64 | float64](leftVal, rightVal T) object.Object {
	// if left is 0 and right is negative,
	// convert right to positive
	// this will ensure that operations like 0/-1 will return 0, not -0
	if leftVal == 0 && rightVal < 0 {
		rightVal = rightVal * -1
	}

	result := float64(leftVal) / float64(rightVal)
	resultStr := fmt.Sprint(result)
	if !strings.Contains(resultStr, ".") {
		return &object.Integer{Value: int64(result)}
	}
	return &object.Float{Value: result}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object, node *ast.InfixExpression,
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
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
	}
}

func evalMapInfixExpression(operator string, left, right object.Object, node *ast.InfixExpression) object.Object {
	switch operator {
	case token.MINUS:
		hashKey, ok := right.(object.Hashable)
		if !ok {
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.UnusableAsHashKeyError(right.Inspect(), r))
		}
		pairs := left.(*object.Hash).Pairs
		for _, key := range pairs {
			if key.Key.Inspect() == right.Inspect() {
				delete(pairs, hashKey.HashKey())
			}
		}

		return &object.Hash{
			Pairs: pairs,
		}

	default:
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
	}
}

func evalArrayInfixExpression(operator string, left, right object.Object, node *ast.InfixExpression) object.Object {
	leftVal := left.(*object.Array).Elements
	if right.Type() == object.ARRAY_OBJ {
		rightVal := right.(*object.Array).Elements

		switch operator {
		case token.PLUS:
			return &object.Array{Elements: append(leftVal, rightVal...)}

		case token.EQ:
			arr1, _ := left.(*object.Array)
			arr2, _ := right.(*object.Array)
			return &object.Boolean{Value: utils.ObjectArrayEqual(arr1, arr2)}

		case token.NOT_EQ:
			arr1, _ := left.(*object.Array)
			arr2, _ := right.(*object.Array)
			return &object.Boolean{Value: !utils.ObjectArrayEqual(arr1, arr2)}

		default:
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
		}
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal := right.(*object.Integer).Value

		switch operator {
		case token.SLASH:
			if arr, ok := left.(*object.Array); ok {
				return utils.SliceChunkAsArrayObject(arr, int(rightVal))
			}

		case token.MINUS:
			if int(rightVal) >= len(leftVal) {
				r, _ := node.TokenInfo.(errors.ErrorConfig)
				return NewError(errors.OutOfRangeError(int(rightVal), len(leftVal), r))
			}
			newArr := append(leftVal[0:rightVal], leftVal[rightVal+1:]...)
			return &object.Array{Elements: newArr}

		case token.ASTERISK:
			newArr := []object.Object{}
			for i := 0; i < int(rightVal); i++ {
				newArr = append(newArr, &object.Array{Elements: leftVal})
			}
			return &object.Array{Elements: newArr}

		default:
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.UnknownOperatorError(operator, string(left.Type()), string(right.Type()), r))
		}
	}

	r, _ := node.TokenInfo.(errors.ErrorConfig)
	return NewError(errors.UnacceptableRHSInArrayInfixExpressionError(operator, string(right.Type()), r))
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

	r, _ := node.TokenInfo.(errors.ErrorConfig)
	return NewError(errors.IdentifierNotDefinedError(node.Value, r))
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

func isContinue(obj object.Object) bool {
	return obj.Type() == object.CONTINUE_OBJ
}

func isBreak(obj object.Object) bool {
	return obj.Type() == object.BREAK_OBJ
}

func NewError(conf errors.Error) *object.Error {
	return &object.Error{Conf: conf}
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
		return NewError(errors.TypeError(fn.Inspect(), object.FUNCTION_OBJ))
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
		} else {
			env.Set(param.Value, &object.Null{})
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object, node *ast.IndexExpression) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, node)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringindexExpression(left, index, node)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index, node)
	default:
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.IndexOperatorNotAllowed(string(left.Type()), r))
	}
}

func evalStringindexExpression(str, index object.Object, node *ast.IndexExpression) object.Object {
	/*
		- cast to *object.String
		- get literal of index
		- test index literal is within bounds such that 0 <= index <= len(str)
		- return value at index
	*/
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value))

	if len(strObject.Value) == 0 || idx > max {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.OutOfRangeError(int(idx), int(max), r))
	}

	if idx < 0 {
		idx = max + 1 + idx
	}

	return &object.String{
		Value: string(strObject.Value[idx]),
	}
}

func evalArrayIndexExpression(array, index object.Object, node *ast.IndexExpression) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if len(arrayObject.Elements) == 0 {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.OutOfRangeError(int(idx), int(max), r))
	}

	if idx > max {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.OutOfRangeError(int(idx), int(max), r))
	} else if idx < 0 {
		idx = max + 1 + idx
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
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.UnusableAsHashKeyError(key.Inspect(), r))
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

func evalHashIndexExpression(hash, index object.Object, node *ast.IndexExpression) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnusableAsHashKeyError(index.Inspect(), r))
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
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.IllegalTokenError(tokenLiteral, r))
		}
	}

	switch node.Operator {
	case token.POST_INCR:
		integer, ok := tok.(*object.Integer)
		if !ok {
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.UnacceptableLHSInPostfixExpression(node.Operator, tokenLiteral, r))
		}

		newInteger := &object.Integer{Value: integer.Value + 1}
		env.Set(tokenLiteral, newInteger)
		return newInteger

	case token.POST_DECR:
		integer, ok := tok.(*object.Integer)
		if !ok {
			r, _ := node.TokenInfo.(errors.ErrorConfig)
			return NewError(errors.UnacceptableLHSInPostfixExpression(node.Operator, tokenLiteral, r))
		}

		newInteger := &object.Integer{Value: integer.Value - 1}
		env.Set(tokenLiteral, newInteger)
		return newInteger
	}
	r, _ := node.TokenInfo.(errors.ErrorConfig)
	return NewError(errors.UnknownOperatorError(node.Operator, "", "", r))
}

func evalArraySquareBracketExpression(left *object.Object, index, value object.Object, node *ast.SquareBracketAssignment) object.Object {
	arr := (*left).(*object.Array)
	idx := index.(*object.Integer).Value

	if int(idx) > len(arr.Elements) {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.OutOfRangeError(int(idx), len(arr.Elements), r))
	}

	arr.Elements[idx] = value

	return arr
}

func evalMapSquareBracketExpression(left *object.Object, index, value object.Object, node *ast.SquareBracketAssignment) object.Object {
	hash := (*left).(*object.Hash)

	hashKey, ok := index.(object.Hashable)
	if !ok {
		r, _ := node.TokenInfo.(errors.ErrorConfig)
		return NewError(errors.UnusableAsHashKeyError(index.Inspect(), r))
	}
	hash.Pairs[hashKey.HashKey()] = object.HashPair{
		Key:   index,
		Value: value,
	}
	return hash
}
