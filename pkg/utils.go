package pkg

import "strings"

func IsEmptyOrBlankString(str string) bool {
	return strings.TrimSpace(str) == ""
}
