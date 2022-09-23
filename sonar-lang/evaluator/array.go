package evaluator

import (
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/errors"
	"github.com/icheka/sonar-lang/sonar-lang/object"
)

var ArrayBuiltins = map[string]*object.Builtin{
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return NewError(errors.RequiresXArgumentsError(2, len(args), "push"))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError(errors.ArgumentToXMustBeYError("array", "push", object.ARRAY_OBJ, args[0].Inspect()))
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
				return NewError(errors.RequiresAtLeastXArgumentsError("pop", 1, len(args)))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError(errors.ArgumentToXMustBeYError("array", "pop", object.ARRAY_OBJ, args[0].Inspect()))
			}

			arr := args[0].(*object.Array)
			idx := len(arr.Elements) - 1
			if len(args) == 2 {
				if args[1].Type() != object.INTEGER_OBJ {
					return NewError(errors.ArgumentToXMustBeYError("index", "pop", object.INTEGER_OBJ, args[1].Inspect()))
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

func arrayContains(obj *object.Array, elm object.Object) bool {
	elements := obj.Elements
	for _, v := range elements {
		if v.Inspect() == elm.Inspect() && v.Type() == elm.Type() {
			return true
		}
	}
	return false
}

func stringContains(obj *object.String, elm object.Object) bool {
	return strings.Contains(obj.Value, elm.Inspect())
}
