package utils

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
