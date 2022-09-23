package evaluator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/errors"
	"github.com/icheka/sonar-lang/sonar-lang/object"
	"github.com/icheka/sonar-lang/sonar-lang/utils"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "len"))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			default:
				return NewError(errors.TypeOfArgumentNotAllowed("len", "iterable", string(arg.Type()), []string{"ITERABLE"}))
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
				return NewError(errors.RequiresAtLeastXArgumentsError("slice", len(args), 1))
			}
			obj := args[0]
			if obj.Type() != object.ARRAY_OBJ && obj.Type() != object.STRING_OBJ {
				return NewError(errors.ArgumentToXAtYMustBeZError(0, "slice", "ARRAY or STRING", string(args[0].Type())))
			}
			if len(args) > 4 {
				return NewError(errors.RequiresAtMostXArgumentsError("slice", len(args), 4))
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
				return NewError(errors.RequiresXArgumentsError(2, len(args), "contains"))
			}
			obj := args[0]
			elm := args[1]

			switch obj.Type() {
			case object.ARRAY_OBJ:
				return &object.Boolean{Value: arrayContains(obj.(*object.Array), elm)}
			case object.STRING_OBJ:
				return &object.Boolean{Value: stringContains(obj.(*object.String), elm)}
			default:
				return NewError(errors.ArgumentToXAtYMustBeZError(0, "contains", "ARRAY or STRING", string(args[0].Type())))
			}
		},
	},
	"copy": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "copy"))
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
				return &object.Array{Elements: obj.(*object.Array).Elements}

			case object.HASH_OBJ:
				return &object.Hash{Pairs: obj.(*object.Hash).Pairs}

			default:
				return NewError(errors.TypeCannotBeCopiedError(string(obj.Type())))
			}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "type"))
			}
			return &object.String{
				Value: string(args[0].Type()),
			}
		},
	},
	"index": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError(errors.RequiresXArgumentsError(2, len(args), "index"))
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
				return NewError(errors.ArgumentToXAtYMustBeZError(0, "index", "ARRAY or STRING", string(args[0].Type())))
			}
		},
	},
	"sort": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "sort"))
			}
			switch args[0].Type() {
			case object.ARRAY_OBJ:
				return utils.SortObjectArray(args[0].(*object.Array))

			default:
				return NewError(errors.ArgumentToXMustBeYError("array", "sort", "ARRAY", string(args[0].Type())))
			}
		},
	},
	"reverse": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "reverse"))
			}
			fmt.Println(args[0].Inspect())
			switch args[0].Type() {
			case object.ARRAY_OBJ:
				return &object.Array{Elements: utils.ReverseSlice(args[0].(*object.Array).Elements)}

			default:
				return NewError(errors.ArgumentToXMustBeYError("array", "reverse", "ARRAY", string(args[0].Type())))
			}
		},
	},
	"range": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return NewError(errors.RequiresAtLeastXArgumentsError("range", len(args), 2))
			}
			if len(args) > 3 {
				return NewError(errors.RequiresAtMostXArgumentsError("range", len(args), 3))
			}
			for i, arg := range args {
				if _, ok := arg.(*object.Integer); !ok {
					return NewError(errors.ArgumentToXAtYMustBeZError(i, "range", object.INTEGER_OBJ, arg.Inspect()))
				}
			}

			start := args[0].(*object.Integer)
			end := args[1].(*object.Integer)
			step := &object.Integer{Value: 1}
			if len(args) == 3 {
				step = args[2].(*object.Integer)
			}

			S := start.Value
			E := end.Value
			St := step.Value

			// if start == end, return []
			// if step is a negative number, start must be > end
			// if step is a positive number, start must be < end
			if S == E {
				return &object.Array{Elements: []object.Object{}}
			}
			if St < 0 && S < E {
				return &object.Array{Elements: []object.Object{}}
			}
			if St > 0 && S > E {
				return &object.Array{Elements: []object.Object{}}
			}

			result := []int64{}
			if St > 0 {
				for i := S; i < E; i += St {
					result = append(result, i)
				}
			} else {
				for i := S; i > E; i += St {
					result = append(result, i)
				}
			}

			arr := []object.Object{}
			for _, v := range result {
				arr = append(arr, &object.Integer{Value: v})
			}

			return &object.Array{Elements: arr}
		},
	},
}

func InitStdlib() {
	var stdlibFunctions = []map[string]*object.Builtin{
		ArrayBuiltins,
		MapBuiltins,
		TypesBuiltins,
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
			return NewError(errors.ArgumentToXAtYMustBeZError(i, "slice", object.INTEGER_OBJ, string(arg.Type())))
		}
	}
	start := int(args[1].(*object.Integer).Value)
	originalArray := obj.(*object.Array).Elements
	if start >= len(originalArray) {
		return NewError(errors.OutOfRangeError(start, len(originalArray), errors.ErrorConfig{}))
	}
	if start < 0 {
		start = len(originalArray) + start
	}
	if len(args) == 2 {
		return &object.Array{Elements: originalArray[start:]}
	}

	end := int(args[2].(*object.Integer).Value)
	if end > len(originalArray) {
		return NewError(errors.OutOfRangeError(end, len(originalArray), errors.ErrorConfig{}))
	}
	if end < 0 {
		end = len(originalArray) + end
	}
	if end < start {
		return NewError(errors.InvalidRangeError(start, end, errors.ErrorConfig{}))
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
