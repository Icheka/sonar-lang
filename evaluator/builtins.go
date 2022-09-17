package evaluator

import (
	"bytes"
	"fmt"
	"sonar/v2/object"
	"sonar/v2/utils"
	"strings"
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
			arr := []string{}
			for _, arg := range args {
				if arg.Type() == object.STRING_OBJ {
					arr = append(arr, arg.(*object.String).FormattedInspect())
				} else {
					arr = append(arr, arg.Inspect())
				}
			}
			fmt.Println(strings.Join(arr, ", "))

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
				return SliceArray(args...)
			case object.STRING_OBJ:
				arr := strings.Split(obj.(*object.String).Value, "")
				objs := []object.Object{}
				for _, v := range arr {
					objs = append(objs, &object.String{Value: v})
				}
				args[0] = &object.Array{Elements: objs}
				// newArr := SliceArray(&object.Array{Elements: objs}, args[])
				newArr := SliceArray(args...)
				if newArr.Type() == object.ERROR_OBJ {
					return newArr
				}
				var strs bytes.Buffer
				for _, v := range newArr.(*object.Array).Elements {
					strs.WriteString(v.(*object.String).Value)
				}

				return &object.String{Value: strs.String()}
			default:
				return &object.Array{Elements: []object.Object{}}
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("`contains` requires 2 arguments, got=%d", len(args))
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
	"index": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("index() takes 2 arguments, %d given", len(args))}
			}

			switch args[0].Type() {
			case object.ARRAY_OBJ:
				return ArrayIndexOf(args[0].(*object.Array), args[1])
			case object.STRING_OBJ:
				elements := []object.Object{}
				for _, v := range strings.Split(args[0].(*object.String).Value, "") {
					elements = append(elements, &object.String{Value: v})
				}
				return ArrayIndexOf(&object.Array{Elements: elements}, args[1])

			default:
				return NewError("type of argument 1 to index() must be ARRAY or STRING")
			}
		},
	},
	"sort": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("sort() takes 1 argument, %d given", len(args))}
			}
			switch args[0].Type() {
			case object.ARRAY_OBJ:
				return utils.SortObjectArray(args[0].(*object.Array))

			default:
				return NewError("argument to sort() must be of type ARRAY")
			}
		},
	},
	"reverse": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("reverse() takes 1 argument, %d given", len(args))}
			}
			switch args[0].Type() {
			case object.ARRAY_OBJ:
				return &object.Array{Elements: utils.ReverseSlice(args[0].(*object.Array).Elements)}

			default:
				return NewError("argument to reverse() must be of type ARRAY")
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

func ArrayIndexOf(arr *object.Array, element object.Object) *object.Integer {
	elements := arr.Elements
	for i, v := range elements {
		if v.Inspect() == element.Inspect() && v.Type() == element.Type() {
			return &object.Integer{Value: int64(i)}
		}
	}
	return &object.Integer{Value: -1}
}

func SliceArray(args ...object.Object) object.Object {
	obj := args[0]
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
	if start >= len(originalArray) {
		return NewError("'start' argument out of bounds")
	}
	if start < 0 {
		start = len(originalArray) + start
	}
	if len(args) == 2 {
		return &object.Array{Elements: originalArray[start:]}
	}

	end := int(args[2].(*object.Integer).Value)
	if end > len(originalArray) {
		return NewError("'end' argument out of bounds")
	}
	if end < 0 {
		end = len(originalArray) + end
	}
	if end < start {
		return NewError("invalid range %d:%d", start, end)
	}

	slicedArr := originalArray[start:end]
	newArray := []object.Object{}

	if len(args) == 3 {
		newArray = slicedArr
	} else {
		step := int(args[3].(*object.Integer).Value)
		if step < 0 {
			step = len(originalArray) + step
		}
		for k, v := range slicedArr {
			if k >= end {
				break
			} else if k%step == 0 {
				newArray = append(newArray, v)
			}
		}
	}

	return &object.Array{Elements: newArray}
}
