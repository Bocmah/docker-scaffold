package dockercompose

import (
	"fmt"
	"strconv"
)

func DoubleQuotted(str string) string {
	if str == "" {
		return ""
	}

	return strconv.Quote(str)
}

func Mapping(str1, str2 string) string {
	if str1 == "" {
		return str2
	}

	if str2 == "" {
		return str1
	}

	return fmt.Sprintf("%s:%s", str1, str2)
}
