// Package geeCache @Author Bing
// @Date 2023/10/24 13:58:00
// @Desc
package geeCache

import (
	"testing"
)

var db = map[string]string{"user": "admin", "password": "pass"}

func call[K string, V string](k K) (V, bool) {
	var v string
	v, ok := db[string(k)]
	if ok {
		return V(v), true
	}
	return V(v), false

}

func TestGroup(t *testing.T) {
	group := NewGroup(2<<10, "test", call[string, string])
	for user, _ := range db {
		r1, ok := group.Get(user)
		if ok {
			t.Log("r1:", r1)
			continue
		}
		r2, ok := group.callback(user)
		if ok {
			t.Log("r2:", r2)
		}

	}
}
