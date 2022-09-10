package utils

import (
	"sonar/v2/object"
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
