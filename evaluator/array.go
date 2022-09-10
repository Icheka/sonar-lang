package evaluator

import (
	"sonar/v2/object"
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
	"slice": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return NewError("`slice` requires at least 1 argument")
			}
			obj := args[0]
			if obj.Type() != object.ARRAY_OBJ && obj.Type() != object.STRING_OBJ {
				return NewError("first argument to `slice` must be of type ARRAY or STRING, got %s",
					args[0].Type())
			}
			if len(args) > 4 {
				return NewError("`slice` requires at most 4 arguments, %d given", len(args))
			}
			for i, arg := range args[1:] {
				if arg.Type() != object.INTEGER_OBJ {
					return NewError("ordinal argument to `slice` at index %d must be INTEGER, %s given", i, arg.Type())
				}
			}
			start := 0
			if len(args) > 1 {
				start = int(args[1].(*object.Integer).Value)
			}
			end := 0
			if len(args) > 2 {
				end = int(args[2].(*object.Integer).Value)
			}
			step := 0
			if len(args) > 3 {
				step = int(args[3].(*object.Integer).Value)
			}

			switch obj.Type() {
			case object.ARRAY_OBJ:
				elements := obj.(*object.Array).Elements
				newArr := []object.Object{}

				// similarly to Python
				// allow expressions like slice(a, -1, ...)
				if start < 0 {
					// wrap start around length of array
					start = len(elements) + start
				}
				// allow expressions like slice(a, 0, -1, ...)
				if end < 0 {
					// wrap start around length of array
					end = len(elements) + end
				} else if end == 0 {
					end = len(elements)
				}
				if step == 0 {
					newArr = elements[start:end]
				} else {
					for i, element := range elements[start:end] {
						if i%step == 0 {
							newArr = append(newArr, element)
						}
					}
				}
				return &object.Array{Elements: newArr}
			default:
				return &object.Array{Elements: []object.Object{}}
			}
		},
	},
}

func checkArity(args []object.Object, expected int) (bool, *object.Error) {
	if len(args) != expected {
		return false, WrongArityError(len(args), 2)
	}
	return true, nil
}
