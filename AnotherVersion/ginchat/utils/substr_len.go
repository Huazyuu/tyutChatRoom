package utils

import "unicode/utf8"

func MbStrLen(str string) int {
	return utf8.RuneCountInString(str)
}
