package evaluator

import (
	"sonar/v2/object"
	"strings"
)

var ArrayBuiltins = map[string]*object.Builtin{
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return NewError("`push` requires at least 2 arguments")
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("'array' argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arrObj := args[0].(*object.Array)
			toBePushed := args[1:]

			newArrElements := append(arrObj.Elements, toBePushed...)
			return &object.Array{Elements: newArrElements}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return NewError("`pop` requires at least 1 argument")
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("'array' argument to `pop` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			idx := len(arr.Elements) - 1
			if len(args) == 2 {
				if args[1].Type() != object.INTEGER_OBJ {
					return NewError("'index' argument to `pop` must be INTEGER, got %s",
						args[1].Type())
				}
				idx = int(args[1].(*object.Integer).Value)
			}

			newArr := arr.Elements[0:idx]
			newArr = append(newArr, arr.Elements[idx+1:]...)

			poppedElement := arr.Elements[idx]

			arr.Elements = newArr

			return poppedElement
		},
	},
}

func checkArity(args []object.Object, expected int) (bool, *object.Error) {
	if len(args) != expected {
		return false, WrongArityError(len(args), 2)
	}
	return true, nil
}

func arrayContains(obj *object.Array, elm object.Object) bool {
	elements := obj.Elements
	for _, v := range elements {
		if v.Inspect() == elm.Inspect() {
			return true
		}
	}
	return false
}

func stringContains(obj *object.String, elm object.Object) bool {
	print(obj.Value, elm.Inspect())
	return strings.Contains(obj.Value, elm.Inspect())
}
