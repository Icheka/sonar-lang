package evaluator

import (
	"strconv"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/object"
)

var TypesBuiltins = map[string]*object.Builtin{
	"convertable": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("convertable() takes 2 arguments, %d given",
					len(args))
			}

			value := args[0]
			to := args[1]
			result := false

			switch strings.ToUpper(to.Inspect()) {
			case object.STRING_OBJ:
				result = true
			case object.INTEGER_OBJ:
				r := toInteger(value)
				if r.Type() == object.INTEGER_OBJ {
					result = true
				}
			case object.FLOAT_OBJ:
				r := toFloat(value)
				if r.Type() == object.FLOAT_OBJ {
					result = true
				}
			case object.HASH_OBJ:
				if value.Type() == object.ARRAY_OBJ {
					result = true
				}
			}

			return &object.Boolean{Value: result}
		},
	},
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("string() takes 1 argument, %d given",
					len(args))
			}

			from := args[0]

			switch from.Type() {
			case object.STRING_OBJ:
				return from
			default:
				return &object.String{Value: from.Inspect()}
			}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("int() takes 1 argument, %d given",
					len(args))
			}

			return toInteger(args[0])
		},
	},
	"float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("float() takes 1 argument, %d given",
					len(args))
			}

			return toFloat(args[0])
		},
	},
	"map": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("map() takes 1 argument, %d given",
					len(args))
			}

			switch args[0].Type() {
			case object.ARRAY_OBJ:
				elements := args[0].(*object.Array).Elements
				pairs := map[object.HashKey]object.HashPair{}

				for i, v := range elements {
					pairs[object.HashKey{Type: object.INTEGER_OBJ, Value: uint64(i)}] = object.HashPair{
						Key:   &object.Integer{Value: int64(i)},
						Value: v,
					}
				}

				return &object.Hash{Pairs: pairs}
			default:
				return NewError("Expected argument to ARRAY, got %s", args[0].Type())
			}
		},
	},
}

func toFloat(from object.Object) object.Object {
	switch from.Type() {
	case object.STRING_OBJ:
		f := from.(*object.String).Value
		fl, err := strconv.ParseFloat(f, 64)
		if err == nil {
			return &object.Float{Value: float64(fl)}
		}

	case object.FLOAT_OBJ:
		return from

	case object.INTEGER_OBJ:
		return &object.Float{Value: float64(from.(*object.Integer).Value)}
	}

	return illegalConversion(from, object.FLOAT_OBJ)
}

func toInteger(from object.Object) object.Object {
	switch from.Type() {
	case object.STRING_OBJ:
		f := from.(*object.String).Value
		fl, err := strconv.ParseFloat(f, 64)
		if err == nil {
			return &object.Integer{Value: int64(fl)}
		}

	case object.INTEGER_OBJ:
		return from

	case object.FLOAT_OBJ:
		return &object.Integer{Value: int64(from.(*object.Float).Value)}
	}

	return illegalConversion(from, object.INTEGER_OBJ)
}

func illegalConversion(from object.Object, to object.ObjectType) *object.Error {
	return NewError("IllegalConversionError: %s to %s", from.Type(), to)
}

var ConvertableMap map[object.ObjectType][]object.ObjectType = map[object.ObjectType][]object.ObjectType{
	FALSE.Type(): {object.STRING_OBJ},

	object.INTEGER_OBJ: {object.STRING_OBJ, object.FLOAT_OBJ},
	object.FLOAT_OBJ:   {object.STRING_OBJ, object.INTEGER_OBJ},

	object.STRING_OBJ: {object.INTEGER_OBJ, object.FLOAT_OBJ},

	object.ARRAY_OBJ:    {},
	object.HASH_OBJ:     {object.ARRAY_OBJ},
	object.FUNCTION_OBJ: {},
	object.BUILTIN_OBJ:  {},
	object.ERROR_OBJ:    {},
	object.NULL_OBJ:     {},
}
