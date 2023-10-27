// Package geeCache @Author Bing
// @Date 2023/10/24 17:24:00
// @Desc
package core

import (
	"testing"
)

func TestLru_Add(t *testing.T) {

	l := NewLru[string, int](10)
	l.Add("one", 1)
	l.Add("two", 2)
	l.Add("three", 3)
	if _, ok := l.Get("one"); !ok {
		return
	}
	l.Range(func(k, v any) bool {
		t.Log(k, v)
		return true
	})
}
