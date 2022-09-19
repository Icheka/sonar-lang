package utils

import "strings"

func StripWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		ws := []string{
			" ", "\t", "\r", "\n",
		}
		if !SliceContains(ws, string(r)) {
			return r
		}
		return -1
	}, s)
}
