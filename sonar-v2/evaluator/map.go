package evaluator

import "sonar/v2/object"

var MapBuiltins = map[string]*object.Builtin{
	"mapKeys": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("`mapKeys` requires only 1 argument")
			}

			if args[0].Type() != object.HASH_OBJ {
				return NewError("'map' argument to `mapKeys` must be MAP, got %s",
					args[0].Type())
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
				return NewError("`mapValues` requires only 1 argument")
			}

			if args[0].Type() != object.HASH_OBJ {
				return NewError("'map' argument to `mapValues` must be MAP, got %s",
					args[0].Type())
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
				return NewError("`mapEntries` requires only 1 argument")
			}

			if args[0].Type() != object.HASH_OBJ {
				return NewError("'map' argument to `mapEntries` must be MAP, got %s",
					args[0].Type())
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
