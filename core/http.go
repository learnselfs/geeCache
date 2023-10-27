// Package core @Author Bing
// @Date 2023/10/25 10:38:00
// @Desc
package core

import (
	"encoding/json"
	"fmt"
	"github.com/learnselfs/geeCache/utils"
	"net/http"
	"strings"
)

type Pool struct {
	group    *Group
	address  string
	port     string
	addr     string
	basePath string
	// request
	r *http.Request
	// response
	w http.ResponseWriter
}

func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.w = w
	p.r = r
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		utils.Log("not found %s", r.URL.Path)
	}
	path := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	path = utils.ExcludeEmpty(path)
	groupName := path[0]
	key := path[1]
	g, ok := p.group.Groups[groupName]
	if !ok {
		p.JSON(utils.Failed())
		return
	}
	v, ok := g.Get(key)
	if ok {
		p.JSON(utils.OkWithDetail(http.StatusOK, "selected successfully", utils.H{"groupName": groupName, "key": key, "value": v}))
	} else {
		p.JSON(utils.Failed())
	}
}

func (p *Pool) JSON(obj utils.H) {
	p.SetHeader("Content-Type", "application/json")
	w := json.NewEncoder(p.w)
	err := w.Encode(obj)
	if err != nil {
		p.Log("JSON: %s", obj)
	}
}

func (p *Pool) Log(f string, v ...any) {
	_, err := p.w.Write([]byte(fmt.Sprintf(f, v...)))
	utils.Log(f, v...)
	if err != nil {
		return
	}
}

func (p *Pool) SetHeader(k, v string) {
	p.w.Header().Set(k, v)
}

func NewPool(group *Group, address string, port string, basePath string) *Pool {
	var builder strings.Builder
	builder.WriteString(address)
	builder.WriteString(":")
	builder.WriteString(port)
	return &Pool{
		group:    group,
		address:  address,
		port:     port,
		addr:     builder.String(),
		basePath: basePath,
	}
}

func Run(group *Group, address, port, basePath string) {
	pool := NewPool(group, address, port, basePath)
	err := http.ListenAndServe(pool.addr, pool)
	if err != nil {
		return
	}
}
