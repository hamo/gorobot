package utils

import (
	"strings"
)

func SplitAndTrim(str, s string) []string {
	strs := strings.Split(str, s)
	for k, v := range strs {
		strs[k] = strings.TrimSpace(v)
	}
	return strs
}
