package utils

import (
	"sort"

	"github.com/icheka/sonar-lang/sonar-lang/object"
)

func SliceContains[T comparable](slice []T, elm T) bool {
	for _, v := range slice {
		if v == elm {
			return true
		}
	}
	return false
}

func SliceChunk[T interface{}](arr []T, size int) [][]T {
	var chunks [][]T

	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}

	return chunks
}

func SliceChunkAsArrayObject(arr *object.Array, size int) *object.Array {
	elements := arr.Elements
	chunks := SliceChunk(elements, size)

	result := []object.Object{}
	for _, v := range chunks {
		result = append(result, &object.Array{Elements: v})
	}

	return &object.Array{Elements: result}
}

func SortObjectArray(arr *object.Array) *object.Array {
	elements := make([]object.Object, len(arr.Elements))
	copy(elements, arr.Elements)

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Inspect() < elements[j].Inspect()
	})
	return &object.Array{Elements: elements}
}

func SortObjectArrayWithFunction(arr *object.Array) [][]object.Object {
	args := [][]object.Object{}

	for i, v := range arr.Elements {
		ownArr := interface{}(arr.Elements).(*object.Object)
		arg := []object.Object{}

		if i == len(arr.Elements)-1 {
			arg = []object.Object{
				v,
				v,
				*ownArr,
			}
		} else {
			arg = []object.Object{
				&object.Integer{Value: int64(i)},
				v,
				arr.Elements[i+1],
				*ownArr,
			}
		}
		args = append(args, arg)
	}

	return args
}

func ObjectArrayEqual(arr1, arr2 *object.Array) bool {
	len1 := len(arr1.Elements)
	len2 := len(arr2.Elements)
	if len1 != len2 {
		return false
	}

	for i, v := range arr1.Elements {
		elm := arr2.Elements[i]
		if elm.Type() != v.Type() || elm.Inspect() != v.Inspect() {
			return false
		}
	}
	return true
}

func ReverseSlice[T ~[]E, E any](arr T) T {
	S := make(T, len(arr))
	copy(S, arr)

	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		S[i], S[j] = S[j], S[i]
	}

	return S
}
