// Package core @Author Bing
// @Date 2023/10/30 14:33:00
// @Desc
package core

import (
	"strconv"
	"testing"
)

func TestConsistent(t *testing.T) {
	m := NewMap(3, func(d []byte) uint32 {
		i, _ := strconv.Atoi(string(d))
		return uint32(i)
	})
	keys := []string{"1", "2", "5"}
	m.Add(keys...)
	cases := map[string]string{"2": "2", "11": "1", "23": "5", "27": "1"}
	for k, v := range cases {
		result := m.Get(k)
		if result != v {
			t.Errorf("input value:%s; get value:%s, want:%s", k, result, v)
		} else {
			t.Logf("input value:%s; get value:%s, want:%s", k, result, v)
		}

	}
}
