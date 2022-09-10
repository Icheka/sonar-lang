package evaluator

import (
	"fmt"
	"sonar/v2/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		default:
			return NewError("argument to `len` not supported, got %s",
				args[0].Type())
		}
	},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"slice": {
		Fn: func(args ...object.Object) object.Object {
			/*
				slice(arr)
				-- returns a copy of arr
				slice(arr, n)
				-- returns a copy of arr[n:]
				slice(arr, n, m)
				-- returns a copy of arr[n:m]
				slice(arr, n, m, y)
				-- returns a copy of arr[n:m:y]
				-- indices wrap around len(arr) when negative
			*/
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

			switch obj.Type() {
			case object.ARRAY_OBJ:
				if len(args) == 1 {
					return &object.Array{Elements: obj.(*object.Array).Elements}
				}

				for i, arg := range args[1:] {
					if arg.Type() != object.INTEGER_OBJ {
						return NewError("ordinal argument to `slice` at index %d must be INTEGER, %s given", i, arg.Type())
					}
				}
				start := int(args[1].(*object.Integer).Value)
				originalArray := obj.(*object.Array).Elements
				if start < 0 {
					start = len(originalArray) + start
				}
				if len(args) == 2 {
					return &object.Array{Elements: originalArray[start:]}
				}

				end := int(args[2].(*object.Integer).Value)
				if end < 0 {
					end = len(originalArray) + end
				}

				newArray := []object.Object{}
				if len(args) == 3 {
					newArray = originalArray[start:end]
				} else {
					step := int(args[3].(*object.Integer).Value)
					if step < 0 {
						step = len(originalArray) + step
					}
					for k, v := range originalArray {
						if k%step == 0 {
							newArray = append(newArray, v)
						}
					}
				}

				return &object.Array{Elements: newArray}
			default:
				return &object.Array{Elements: []object.Object{}}
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("`contains` requires at 2 arguments, got=%d", len(args))
			}
			obj := args[0]
			elm := args[1]

			switch obj.Type() {
			case object.ARRAY_OBJ:
				return &object.Boolean{Value: arrayContains(obj.(*object.Array), elm)}
			case object.STRING_OBJ:
				return &object.Boolean{Value: stringContains(obj.(*object.String), elm)}
			default:
				return NewError("first argument to `contains` must be of type ARRAY or STRING, got %s",
					args[0].Type())
			}
		},
	},
	"copy": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("copy() takes 1 argument, %d given", len(args))}
			}
			obj := args[0]

			switch obj.Type() {
			case object.INTEGER_OBJ:
				return &object.Integer{Value: obj.(*object.Integer).Value}

			case object.FLOAT_OBJ:
				return &object.Float{Value: obj.(*object.Float).Value}

			case object.BOOLEAN_OBJ:
				return &object.Boolean{Value: obj.(*object.Boolean).Value}

			case object.STRING_OBJ:
				return &object.String{Value: obj.(*object.String).Value}

			case object.FUNCTION_OBJ:
				fn := obj.(*object.Function)
				return &object.Function{
					Parameters: fn.Parameters,
					Body:       fn.Body,
					Env:        fn.Env,
				}

			case object.ARRAY_OBJ:
				arr := obj.(*object.Array)
				return &object.Array{Elements: arr.Elements}

			case object.HASH_OBJ:
				hash := obj.(*object.Hash)
				return &object.Hash{Pairs: hash.Pairs}

			default:
				return &object.Error{Message: fmt.Sprintf("Type %s cannot be copied", obj.Type())}
			}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("type() takes 1 argument, %d given", len(args))}
			}
			return &object.String{
				Value: string(args[0].Type()),
			}
		},
	},
}

func InitStdlib() {
	var stdlibFunctions = []map[string]*object.Builtin{
		ArrayBuiltins,
	}

	for _, f := range stdlibFunctions {
		for k, v := range f {
			builtins[k] = v
		}
	}
}
