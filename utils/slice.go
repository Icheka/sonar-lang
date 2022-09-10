package utils

func SliceContains[T comparable](slice []T, elm T) bool {
	for _, v := range slice {
		if v == elm {
			return true
		}
	}
	return false
}
