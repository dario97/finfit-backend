package pkg

import (
	"strings"
)

func IsEmptyOrBlankString(str string) bool {
	return strings.TrimSpace(str) == ""
}

func HasMin(str string, min int) bool {
	return len(str) >= min
}

func ExceedsMax(str string, max int) bool {
	return len(str) > max
}
