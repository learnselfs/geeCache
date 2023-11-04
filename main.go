// Package main @Author Bing
// @Date 2023/10/25 16:40:00
// @Desc
package main

import (
	"flag"
	"github.com/learnselfs/geeCache/core"
)

var db = map[string]interface{}{"user": "admin", "pass": "12345"}

func loadDb(k string) (any, bool) {
	v, ok := db[k]
	return v, ok
}
func main() {
	var host string
	var port string
	var api bool
	flag.StringVar(&host, "host", "localhost", "http server listening address for host")
	flag.StringVar(&port, "port", "3000", "http server listening address for port")
	flag.BoolVar(&api, "api", false, "api http server")
	flag.Parse()

	nodes := []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082"}

	node := core.NewNode(host, port, api, nodes, loadDb, "")

	node.Run()
}
