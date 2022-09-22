package evaluator

import (
	"fmt"
	"testing"

	"github.com/icheka/sonar-lang/sonar-lang/ast"
	"github.com/icheka/sonar-lang/sonar-lang/object"
)

func TestMapBuiltin(t *testing.T) {
	input := `map([1, 2, "John", true])`
	expected := `{0: 1, 1: 2, 2: 'John', 3: true}`
	testEvalType[*object.Hash](t, input, expected)
}

func TestFloatConvertBuiltin(t *testing.T) {
	tests := []struct {
		obj      object.Object
		expected interface{}
	}{
		{&object.String{Value: "\"1\""}, 1.0},
		{&object.String{Value: "\"1.3\""}, 1.3},
		{&object.String{Value: "\"1.f\""}, false},

		{&object.Float{Value: 1.4}, 1.4},
		{&object.Integer{Value: 1}, 1.0},

		{&object.Array{Elements: []object.Object{}}, false},
		{&object.Hash{Pairs: map[object.HashKey]object.HashPair{}}, false},
	}

	for _, tt := range tests {
		input := fmt.Sprintf("float(%s)", tt.obj.Inspect())
		evaluated := testEval(input)
		switch evaluated.Type() {
		case object.ERROR_OBJ:
			if tt.expected != false {
				t.Errorf("expected evaluated to be ERROR, got=%t", evaluated)
			}
		default:
			i, ok := tt.expected.(float64)
			if !ok {
				t.Errorf("expected tt.expected to be float, got=%T", tt.expected)
			}
			testEvalFloat(t, input, i)
		}
	}
}

func TestIntegerConvertBuiltin(t *testing.T) {
	tests := []struct {
		obj      object.Object
		expected interface{}
	}{
		{&object.String{Value: "\"1\""}, 1},
		{&object.String{Value: "\"1.3\""}, 1},
		{&object.String{Value: "\"1.f\""}, false},

		{&object.Float{Value: 1.4}, 1},
		{&object.Integer{Value: 1}, 1},

		{&object.Array{Elements: []object.Object{}}, false},
		{&object.Hash{Pairs: map[object.HashKey]object.HashPair{}}, false},
	}

	for _, tt := range tests {
		input := fmt.Sprintf("int(%s)", tt.obj.Inspect())
		evaluated := testEval(input)
		switch evaluated.Type() {
		case object.ERROR_OBJ:
			if tt.expected != false {
				t.Errorf("expected evaluated to be ERROR, got=%t", evaluated)
			}
		default:
			i, ok := tt.expected.(int)
			if !ok {
				t.Errorf("expected tt.expected to be int, got=%T", tt.expected)
			}
			testEvalInteger(t, input, i)
		}
	}
}

func TestStringConvertBuiltin(t *testing.T) {
	input := []struct {
		from object.Object
	}{
		{&object.Integer{Value: 1}},
		{&object.Float{Value: 1}},
		{TRUE},
		{FALSE},
		{&object.String{Value: "\"Icheka\""}},
		{&object.Integer{Value: 1}},
		{&object.Array{Elements: []object.Object{&object.Integer{Value: 1}}}},
		{&object.Hash{Pairs: map[object.HashKey]object.HashPair{}}},
	}

	for _, tt := range input {
		evaluated := testEval(fmt.Sprintf("string(%s)", tt.from.Inspect()))
		if evaluated.Type() != object.STRING_OBJ {
			t.Fatalf("expected evaluated to be STRING, got=%T(%+v)", evaluated, evaluated)
		}
	}
}

func TestConvertableBuiltin(t *testing.T) {
	tests := map[object.ObjectType][]object.ObjectType{
		FALSE.Type(): {object.STRING_OBJ},

		object.INTEGER_OBJ: {object.STRING_OBJ, object.FLOAT_OBJ},
		object.FLOAT_OBJ:   {object.STRING_OBJ, object.INTEGER_OBJ},

		object.ARRAY_OBJ:    {},
		object.HASH_OBJ:     {},
		object.FUNCTION_OBJ: {},
		object.BUILTIN_OBJ:  {},
		object.ERROR_OBJ:    {},
		object.NULL_OBJ:     {},
	}

	for k, v := range tests {
		for _, tt := range v {
			from, _ := getValueOfType(k)
			to, _ := getValueOfType(tt)
			testTypeConversion(t, from, to.Type(), TRUE)
		}
	}

	tests = map[object.ObjectType][]object.ObjectType{
		object.STRING_OBJ: {object.INTEGER_OBJ, object.FLOAT_OBJ},
	}

	for k, v := range tests {
		for _, tt := range v {
			from, _ := getValueOfType(k)
			to, _ := getValueOfType(tt)
			testTypeConversion(t, from, to.Type(), FALSE)
		}
	}
}

func testTypeConversion(t *testing.T, value object.Object, to object.ObjectType, expected *object.Boolean) {
	toType, _ := getValueOfType(to)
	vInspected := value.Inspect()
	if value.Type() == object.STRING_OBJ {
		vInspected = fmt.Sprintf("\"%s\"", value.(*object.String).Value)
	}
	input := fmt.Sprintf("convertable(%s, \"%s\")", vInspected, toType.Type())
	evaluated := testEval(input)
	if evaluated.Type() != object.BOOLEAN_OBJ || evaluated.(*object.Boolean).Value != expected.Value {
		t.Errorf("expected '%s' to be %T(%+v), got=%T(%+v)", input, expected, expected, evaluated, evaluated)
	}
}

func getValueOfType(tp object.ObjectType) (object.Object, string) {
	var val object.Object

	switch tp {
	case object.BOOLEAN_OBJ:
		val = TRUE
	case object.INTEGER_OBJ:
		val = &object.Integer{Value: 1}
	case object.FLOAT_OBJ:
		val = &object.Float{Value: 1}
	case object.STRING_OBJ:
		val = &object.String{Value: "Icheka"}
	case object.ARRAY_OBJ:
		val = &object.Array{Elements: []object.Object{&object.String{Value: "Icheka"}, &object.String{Value: "Ozuru"}}}
	case object.HASH_OBJ:
		key := object.HashKey{
			Type:  object.STRING_OBJ,
			Value: 0,
		}
		val = &object.Hash{Pairs: map[object.HashKey]object.HashPair{
			key: {Key: &object.String{Value: "name"}, Value: &object.String{Value: "Icheka"}},
		}}
	case object.FUNCTION_OBJ:
		val = &object.Function{
			Parameters: []*ast.Identifier{{Value: "getName"}},
			Body:       &ast.BlockStatement{Statements: []ast.Statement{&ast.ReturnStatement{ReturnValue: &ast.IntegerLiteral{Value: int64(1)}}}},
			Env:        object.NewEnvironment(),
		}
	case object.BUILTIN_OBJ:
		val = builtins["len"]
	case object.ERROR_OBJ:
		val = NewError("An error")
	case object.NULL_OBJ:
		val = &object.Null{}
	}

	if val.Type() == object.STRING_OBJ {
		return val, "\"" + val.Inspect() + "\""
	}
	return val, val.Inspect()
}
