package evaluator

import (
	"testing"

	"github.com/icheka/sonar-lang/object"
)

func TestReverseBuiltin(t *testing.T) {
	input := `reverse([1, 2, 3])`
	testEvalType[*object.Array](t, input, `[3, 2, 1]`)
}

func TestPushBuiltin(t *testing.T) {
	input := `let a = [1, 2, 3]; a = push(a, 4, 5); a;`
	evaluated := testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 5 {
		t.Fatalf("expected evaluated.Elements to be 5, got=%d", len(arr.Elements))
	}

	input = `let a = [1, 2, 3]; a = push(4, 5); a;`
	evaluated = testEval(input)
	if _, ok = evaluated.(*object.Array); ok {
		t.Fatalf("expected evaluated to be object.Error, got=%s", evaluated.Type())
	}
}

func TestPopBuiltin(t *testing.T) {
	input := `let a = [1, 2, 3]; pop(a);`
	evaluated := testEval(input)
	res, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("expected evaluated to be object.Integer, got=%s", evaluated.Type())
	}
	if res.Value != 3 {
		t.Fatalf("expected res.Value to be 3, got=%d", res.Value)
	}

	input = `let a = [1,2,3]; pop(a); a;`
	evaluated = testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 2 {
		t.Fatalf("expected arr.Elements to be 2, got=%d", len(arr.Elements))
	}
}
