package utils

import (
	"strings"
)

func SplitAndTrim(str, s string) []string {
	return SplitAndTrimN(str, s, -1)
}

func SplitAndTrimN(str, s string, n int) []string {
	strs := strings.SplitN(str, s, n)
	for k, v := range strs {
		strs[k] = strings.TrimSpace(v)
	}
	return strs
}
