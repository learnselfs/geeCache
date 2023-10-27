// Package core  @Author Bing
// @Date 2023/10/24 11:07:00
// @Desc global cache group for get cache callback cache
package core

import "sync"

type callback func(k string) (any, bool)

type Group struct {
	lock     sync.RWMutex
	Groups   map[string]*Group
	callback callback
	lru      *Lru[string, any]
}

var lock sync.RWMutex

func NewGroup(size int, name string, f callback) *Group {
	defer lock.Unlock()
	lock.Lock()
	if f == nil {
		panic("nil callback")
	}
	group := &Group{
		Groups:   make(map[string]*Group),
		callback: f,
		lru:      NewLru[string, any](size),
	}
	group.Groups[name] = group
	return group
}

func (g *Group) GetGroup(name string) (*Group, bool) {
	defer g.lock.RUnlock()
	g.lock.RLock()
	group, ok := g.Groups[name]
	return group, ok
}

func (g *Group) Get(k string) (any, bool) {
	defer g.lock.RUnlock()
	g.lock.RLock()
	v, ok := g.lru.Get(k)
	if ok {
		return v, ok
	}
	return g.load(k)
}

func (g *Group) load(k string) (any, bool) {
	v, ok := g.callback(k)
	if ok {
		g.lru.Add(k, v)
	}
	return v, ok
}
