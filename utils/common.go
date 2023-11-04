// Package utils @Author Bing
// @Date 2023/10/27 15:44:00
// @Desc
package utils

import "strings"

func ExcludeEmpty(li []string) []string {
	var l []string
	for _, v := range li {
		if v != "" {
			l = append(l, v)
		}
	}
	return l
}

func StringJoin(s ...string) string {
	var result strings.Builder
	for _, v := range s {
		result.WriteString(v)
	}
	return result.String()
}
