package utils

import (
	"sonar/v2/object"
	"sort"
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
	elementsMap := map[int]object.Object{}
	elementsSlice := []string{}

	for i, v := range arr.Elements {
		elementsMap[i] = v
		elementsSlice = append(elementsSlice, v.Inspect())
	}

	sort.Slice(elementsSlice, func(i, j int) bool {
		return elementsSlice[i] < elementsSlice[j]
	})

	sortedElements := []object.Object{}
	for _, v := range elementsSlice {
		for _, mapV := range elementsMap {
			if mapV.Inspect() == v {
				sortedElements = append(sortedElements, mapV)
			}
		}
	}

	return &object.Array{Elements: sortedElements}
}
