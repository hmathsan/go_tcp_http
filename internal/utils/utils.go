package utils

import "unicode"

func IsUpper(str string) bool {
	for _, r := range str {
		if (!unicode.IsUpper(r) && unicode.IsLetter(r)) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
