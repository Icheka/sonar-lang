package evaluator

import (
	"sonar/v2/object"
	"testing"
)

func TestSliceBuiltin(t *testing.T) {
	// calling slice without args returns error
	input := `let a = [1, 2, 3]; a = slice();`
	evaluated := testEval(input)
	_, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected evaluated to be object.Error, got=%s", evaluated.Type())
	}

	// calling slice with one arg of type ARRAY or SLICE returns a copy of arg[0]
	input = `let a = [1, 2, 3]; a = slice(a);`
	evaluated = testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected arr.Elements to be 3, got=%d", len(arr.Elements))
	}

	// calling slice with start only returns a copy of arg[0] sliced from an offset of start
	input = `let a = [1, 2, 3]; a = slice(a, 0);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected arr.Elements to be 3, got=%d", len(arr.Elements))
	}
	input = `let a = [1, 2, 3, 4, 5]; a = slice(a, 2);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected arr.Elements to be 3, got=%d", len(arr.Elements))
	}

	// calling slice with start and end only returns arg[0][start:end]
	input = `let a = [1, 2, 3]; a = slice(a, 0, 1);`
	evaluated = testEval(input)
	arr, ok = evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	}
	if len(arr.Elements) != 1 {
		t.Fatalf("expected arr.Elements to be 1, got=%d", len(arr.Elements))
	}

	// calling slice with three indices returns a copy of arg[start:end:slice]
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

	// input = `let a = [1, 2, 3]; a = slice(a, -1, -1);`
	// evaluated = testEval(input)
	// arr, ok = evaluated.(*object.Array)
	// if !ok {
	// 	t.Fatalf("expected evaluated to be object.Array, got=%s", evaluated.Type())
	// }
	// if len(arr.Elements) != 0 {
	// 	t.Fatalf("expected arr.Elements to be 0, got=%d", len(arr.Elements))
	// }
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

	input = `let a = "Icheka"; contains(a, "che");`
	evaluated = testEval(input)
	b, ok = evaluated.(*object.Boolean)
	if !ok {
		t.Fatalf("expected evaluated to be object.Boolean, got=%s", evaluated.Type())
	}
	if b.Value != true {
		t.Fatalf("expected b.Value to be true, got=%t", b.Value)
	}
}

func TestCopyBuiltin(t *testing.T) {
	testEvalInteger(t, `copy(1);`, 1)
	testEvalFloat(t, `copy(1.1);`, 1.1)
	testEvalType[*object.Boolean](t, `copy(true)`, "true")
	testEvalType[*object.String](t, `copy("Icheka")`, "Icheka")
	testEvalType[*object.Function](t, `copy(fn() { return 1 })`, "fn() { return 1; }")
	testEvalType[*object.Array](t, `copy([1,2,3])`, "[1, 2, 3]")
	testEvalType[*object.Hash](t, `copy({1:1})`, "{1: 1}")
}
