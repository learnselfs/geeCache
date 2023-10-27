// Package main @Author Bing
// @Date 2023/10/25 16:40:00
// @Desc
package main

import "github.com/learnselfs/geeCache/core"

var db = map[string]interface{}{"user": "admin", "pass": "12345"}

func loadDb(k string) (any, bool) {
	v, ok := db[k]
	return v, ok
}

func main() {
	group := core.NewGroup(10, "test", loadDb)
	core.Run(group, "localhost", "9999", "/_geeCache/")
}
