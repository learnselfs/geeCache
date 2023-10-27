// Package utils @Author Bing
// @Date 2023/10/27 15:44:00
// @Desc
package utils

func ExcludeEmpty(li []string) []string {
	var l []string
	for _, v := range li {
		if v != "" {
			l = append(l, v)
		}
	}
	return l
}
