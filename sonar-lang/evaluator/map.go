package evaluator

import (
	"github.com/icheka/sonar-lang/sonar-lang/errors"
	"github.com/icheka/sonar-lang/sonar-lang/object"
)

var MapBuiltins = map[string]*object.Builtin{
	"mapKeys": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "mapKeys"))
			}

			if args[0].Type() != object.HASH_OBJ {
				return NewError(errors.ArgumentToXMustBeYError("map", "mapKeys", "MAP", string(args[0].Type())))
			}

			mapObj := args[0].(*object.Hash)
			keys := []object.Object{}

			for _, pair := range mapObj.Pairs {
				keys = append(keys, pair.Key)
			}

			return &object.Array{Elements: keys}
		},
	},
	"mapValues": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "mapValues"))
			}

			if args[0].Type() != object.HASH_OBJ {
				return NewError(errors.ArgumentToXMustBeYError("map", "mapValues", "MAP", string(args[0].Type())))
			}

			mapObj := args[0].(*object.Hash)
			values := []object.Object{}

			for _, pair := range mapObj.Pairs {
				values = append(values, pair.Value)
			}

			return &object.Array{Elements: values}
		},
	},
	"mapEntries": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(errors.RequiresXArgumentsError(1, len(args), "mapEntries"))
			}

			if args[0].Type() != object.HASH_OBJ {
				return NewError(errors.ArgumentToXMustBeYError("map", "mapEntries", "MAP", string(args[0].Type())))
			}

			mapObj := args[0].(*object.Hash)
			entries := [][]object.Object{}

			for _, pair := range mapObj.Pairs {
				entries = append(entries, []object.Object{pair.Key, pair.Value})
			}

			result := []object.Object{}
			for _, v := range entries {
				result = append(result, &object.Array{Elements: v})
			}
			return &object.Array{Elements: result}
		},
	},
}
