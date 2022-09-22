package evaluator

import (
	"strconv"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/object"
	"github.com/icheka/sonar-lang/sonar-lang/utils"
)

var TypesBuiltins = map[string]*object.Builtin{
	"convertable": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("convertable() takes 2 arguments, %d given",
					len(args))
			}

			from := args[0]
			to := args[1]

			if _, ok := object.ObjectTypes[strings.ToLower(to.Inspect())]; !ok {
				return NewError("%s is not a type", to.Inspect())
			}

			if from.Type() == object.STRING_OBJ {
				fromValue := from.(*object.String).Value
				result := false

				switch to.Inspect() {
				case object.INTEGER_OBJ:
					fallthrough
				case object.FLOAT_OBJ:
					if _, err := strconv.ParseFloat(fromValue, 64); err == nil {
						result = true
					}
				}
				return &object.Boolean{Value: result}
			}

			types := ConvertableMap[from.Type()]
			if !utils.SliceContains(types, to.Type()) {
				return FALSE
			}
			return TRUE
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

			from := args[0]

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
		},
	},
	"float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("int() takes 1 argument, %d given",
					len(args))
			}

			from := args[0]

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
		},
	},
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
	object.HASH_OBJ:     {},
	object.FUNCTION_OBJ: {},
	object.BUILTIN_OBJ:  {},
	object.ERROR_OBJ:    {},
	object.NULL_OBJ:     {},
}
