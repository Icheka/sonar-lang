package evaluator

import (
	"sonar/v2/object"
	"testing"
)

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

func TestSliceBuiltin(t *testing.T) {
	input := `let a = [1, 2, 3]; a = slice();`
	evaluated := testEval(input)
	_, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected evaluated to be object.Error, got=%s", evaluated.Type())
	}

	input = `let a = [1, 2, 3]; a = slice(a);`
	evaluated = testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected arr.Elements to be 3, got=%d", len(arr.Elements))
	}

	input = `let a = [1, 2, 3]; a = slice(a, 0);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected arr.Elements to be 3, got=%d", len(arr.Elements))
	}

	input = `let a = [1, 2, 3]; a = slice(a, 0, 1);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 1 {
		t.Fatalf("expected arr.Elements to be 1, got=%d", len(arr.Elements))
	}

	input = `let a = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; a = slice(a, 0, 10, 2);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 5 {
		t.Fatalf("expected arr.Elements to be 5, got=%d", len(arr.Elements))
	}
	elm, ok := arr.Elements[0].(*object.Integer)
	if !ok {
		t.Fatalf("expected first element of arr to be INTEGER, got=%s", arr.Elements[0].Type())
	}

	if elm.Value != 1 {
		t.Fatalf("expected elm.Value to be 5, got=%d", elm.Value)
	}

	input = `let a = [1, 2, 3]; a = slice(a, -1, -1);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 0 {
		t.Fatalf("expected arr.Elements to be 0, got=%d", len(arr.Elements))
	}
}

func TestContainsBuiltin(t *testing.T) {
	input := `let a = [1, 2, 3]; contains(a, 1);`
	evaluated := testEval(input)
	b, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Fatalf("expected evaluated to be object.Boolean, got=%s", evaluated.Type())
	}
	if b.Value != true {
		t.Fatalf("expected b.Value to be true, got=%t", b.Value)
	}
}
