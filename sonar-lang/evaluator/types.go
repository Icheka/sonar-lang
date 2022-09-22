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
