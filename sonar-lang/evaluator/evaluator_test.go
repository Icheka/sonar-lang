package evaluator

import (
	"fmt"
	"testing"

	"github.com/icheka/sonar-lang/sonar-lang/lexer"
	"github.com/icheka/sonar-lang/sonar-lang/object"
	"github.com/icheka/sonar-lang/sonar-lang/parser"
	"github.com/icheka/sonar-lang/sonar-lang/utils"
)

func TestContinueStatement(t *testing.T) {
	input := `
let j = []
let i = 0
for (i, v in range(1, 4)) {
	if (v == 2) {
		continue
	}
	j = push(j, i)
}
len(j)
`
	testEvalInteger(t, input, 2)
}

func TestBreakStatement(t *testing.T) {
	input := `
let i = 0
let j = 0
while (i < 5) {
	i++
	j = i
	if (i == 2) {
		break
	}
}
j
`
	testEvalInteger(t, input, 2)

	input = `
let j = 0
for (i,_ in range(1, 6)) {
	j = i
	if (i == 2) {
		break
	}
}
j
`
	testEvalInteger(t, input, 2)
}

func TestDeleteMapKeyExpression(t *testing.T) {
	input := `
let m  = {1: 2}
m - 1
`
	testEvalType[*object.Hash](t, testEval(input).Inspect(), "{}")
}

func TestSquareBracketAssignmentExpression(t *testing.T) {
	input := "[1,2,3][0] = 4"
	testEvalType[*object.Array](t, testEval(input).Inspect(), "[4, 2, 3]")

	// an out-of-bounds index operation	returns error
	input = "[1,2,3][9] = 4"
	evaluated := testEval(input)
	if _, ok := evaluated.(*object.Error); !ok {
		t.Fatalf("expected evaluated to be ERROR, got=%T", evaluated)
	}

	input = `
let m = {1: 1}
m[1] = 10
`
	testEvalType[*object.Hash](t, testEval(input).Inspect(), `{1: 10}`)
}

func TestAssignmentExpression(t *testing.T) {
	input := `
let a = 1
a = 2
`

	testIntegerObject(t, testEval(input), 2)

	// test that attempting to assign to a variable before it is declared throws error
	errorTest := "a = 2"
	evaluated := testEval(errorTest)
	if _, isErr := evaluated.(*object.Error); !isErr {
		t.Fatalf("expected evaluated to be error, got=%T", evaluated)
	}

	// test double-operator assignments
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"a += 1", 1},
		{"a -= 1", -1},
		{"a *= 2", 0},
		{"a /= 1", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(fmt.Sprintf("let a = 0; %s", tt.input))
		switch v := tt.expectedValue.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case int64:
			testIntegerObject(t, evaluated, v)
		case float64:
			testFloatObject(t, evaluated, v)
		case string:
			testStringObject(t, evaluated, v)
		}
	}
}

func TestWhileStatement(t *testing.T) {
	input := `
let i = 0
while (i < 3) {
	i++
}
i
`
	// test that i is a variable with value = 2
	testIntegerObject(t, testEval(input), 3)
}

func TestForStatement(t *testing.T) {
	input := `
let j = 0
for (i, v in [0, 2]) {
	j = 1
}
j
`
	testIntegerObject(t, testEval(input), 1)

	input = `
let j = 0
for (i, v in "Icheka") {
	j = v
}
j
`
	testStringObject(t, testEval(input), "a")

	input = `
let j = 0
for (k, v in {"name": "Icheka"}) {
	j = k
}
j
`
	testStringObject(t, testEval(input), "name")

	input = `
let j = 0
for (k, v in {"name": "Icheka"}) {
	j = v
}
j
`
	testStringObject(t, testEval(input), "Icheka")
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1 + 1", 2},
		{"1 - 1", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !(testIntegerObject(t, evaluated, tt.expected) || testStringObject(t, evaluated, tt.expected)) {
			t.Fatalf("Unknown type, got %q, expected %q", evaluated.Type(), tt.expected)
		}
	}
}

func TestEvalPostfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1++", 2},
		{"0++", 1},
		{"1--", 0},
		{"0--", -1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1.0", 1},
		{"-1.1", -1.1},
		{"2.1 + 1.0", 3.1},
		{"2.1 + 1", 3.1},
		{"2.1 + 1.2", 3.3},
		{"1 + 1.2", 2.2},
		{"50 / 2 * 2 + 10.1", 60.1},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10.1", 49.9},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"1 > 1.0", false},
		{"1 > 1.2", false},
		{"1 < 1.2", true},
		{"1.3 < 1.2", false},
		{"1.3 > 1.2", true},
		{"1 and 2", true},
		{"false and 2", false},
		{"false or 2", true},
		{"false or 0", true},
		{"false or false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
		{
			`
let f = func(x) {
  return x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
let f = func(x) {
   let result = x + 10;
   return result;
   return 10;
};
f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`{"name": "Monkey"}[func(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
		{
			`999[1]`,
			"index operator not supported: INTEGER",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}

	// test that calling let again on a variable in a scope will throw error
	errorTest := "let a = 1; let a = 2;"
	evaluated := testEval(errorTest)
	if _, isErr := evaluated.(*object.Error); !isErr {
		t.Fatalf("expected evaluated to be error, got=%T", evaluated)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { x; }; identity(5);", 5},
		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { x * 2; }; double(5);", 10},
		{"let add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
let first = 10;
let second = 10;
let third = 10;

let ourFunction = func(first) {
  let second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = func(x) {
  func(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "len() takes 1 argument, 2 given"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "'array' argument to `push` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			if _, ok := evaluated.(*object.Error); !ok {
				t.Fatalf("expected evaluated to be object.Error, got %T", evaluated)
			}
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestStringIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Icheka"[0]`, "I"},
		{`""[0]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if _, ok := tt.expected.(string); ok {
			testStringObject(t, evaluated, tt.expected)
			return
		}
		testNullObject(t, evaluated)
	}
}

func testStringObject(t *testing.T, obj object.Object, expected interface{}) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got %s, expected %s", result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	InitStdlib()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	if result, ok := obj.(*object.Float); ok {
		if result.Value != expected {
			t.Errorf("object has wrong value. got=%f, want=%f",
				result.Value, expected)
			return false
		}
		return true
	}

	t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
	return false
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEvalInteger(t *testing.T, input string, expected int) bool {
	evaluated := testEval(input)
	obj, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("expected evaluated to be INTEGER, got=%s", evaluated.Type())
		return false
	}
	if obj.Value != int64(expected) {
		t.Errorf("expected obj.Value to be %d, got=%d", expected, obj.Value)
		return false
	}
	return true
}

func testEvalFloat(t *testing.T, input string, expected float64) bool {
	evaluated := testEval(input)
	obj, ok := evaluated.(*object.Float)
	if !ok {
		t.Errorf("expected evaluated to be FLOAT, got=%s", evaluated.Type())
		return false
	}
	if obj.Value != expected {
		t.Errorf("expected obj.Value to be %f, got=%f", expected, obj.Value)
		return false
	}
	return true
}

func testEvalType[Type *object.Integer | *object.Float | *object.Boolean | *object.String | *object.Function | *object.Builtin | *object.Array | *object.Hash, Expected int | string | bool](t *testing.T, input string, expected Expected) bool {
	evaluated := testEval(input)
	_, ok := evaluated.(*object.Error)
	if ok {
		t.Errorf("expected evaluated to be T, got=%s(%+v)", evaluated.Type(), evaluated)
		return false
	}

	if evaluated.Type() == object.INTEGER_OBJ {
		return testEvalInteger(t, input, any(expected).(int))
	}
	if evaluated.Type() == object.FLOAT_OBJ {
		return testEvalFloat(t, input, any(expected).(float64))
	}

	compareValue := any(expected).(string)
	var passed bool

	switch evaluated.Type() {
	case object.BOOLEAN_OBJ:
		b := any(evaluated).(*object.Boolean)
		passed = b.Inspect() == compareValue

	case object.STRING_OBJ:
		s := any(evaluated).(*object.String)
		passed = s.Inspect() == compareValue

	case object.FUNCTION_OBJ:
		compareValue = utils.StripWhitespace(compareValue)
		f := any(evaluated).(*object.Function)
		inspected := utils.StripWhitespace(f.Inspect())
		passed = inspected == compareValue

	case object.ARRAY_OBJ:
		a := any(evaluated).(*object.Array)
		passed = a.Inspect() == compareValue

	case object.HASH_OBJ:
		h := any(evaluated).(*object.Hash)
		passed = h.Inspect() == compareValue
	}

	if !passed {
		t.Fatalf("%s is not equal to %s", evaluated.Inspect(), compareValue)
		return false
	}
	return true
}

func TestEvalArrayInfixExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue string
	}{
		{"[1,2,3,4,5] + [6,7,8]", "[1, 2, 3, 4, 5, 6, 7, 8]"},
		{"[1,2,3,4,5,6,7] / 2", "[[1, 2], [3, 4], [5, 6], [7]]"},
		{"[1,2] - 0", "[2]"},
		{"[1,2] * 2", "[[1, 2], [1, 2]]"},
	}

	for _, tt := range tests {
		testEvalType[*object.Array](t, tt.input, tt.expectedValue)
	}

	tests = []struct {
		input         string
		expectedValue string
	}{
		{"[1] == [1]", "true"},
		{"[1] == [2]", "false"},
		{"[1] != [2]", "true"},
		{"[1] != [1]", "false"},
	}

	for _, tt := range tests {
		testEvalType[*object.Boolean](t, tt.input, tt.expectedValue)
	}
}
