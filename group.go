// Package geeCache @Author Bing
// @Date 2023/10/24 11:07:00
// @Desc global cache group for get cache callback cache
package geeCache

import "sync"

type callback[K comparable, V any] func(k K) (V, bool)

type Group[K comparable, V any] struct {
	lock     sync.RWMutex
	Groups   map[string]*Group[K, V]
	callback callback[K, V]
	lru      *Lru[K, V]
}

var lock sync.RWMutex

func NewGroup[K comparable, V any](size int, name string, f callback[K, V]) *Group[K, V] {
	defer lock.Unlock()
	lock.Lock()
	if f == nil {
		panic("nil callback")
	}
	group := &Group[K, V]{
		Groups:   make(map[string]*Group[K, V]),
		callback: f,
		lru:      NewLru[K, V](size),
	}
	group.Groups[name] = group
	return group
}

func (g *Group[K, V]) GetGroup(name string) (*Group[K, V], bool) {
	defer g.lock.RUnlock()
	g.lock.RLock()
	group, ok := g.Groups[name]
	return group, ok
}

func (g *Group[K, V]) Get(k K) (V, bool) {
	defer g.lock.RUnlock()
	g.lock.RLock()
	v, ok := g.lru.Get(k)
	if ok {
		return v, ok
	}
	return g.load(k)
}

func (g *Group[K, V]) load(k K) (V, bool) {
	v, ok := g.callback(k)
	if ok {
		g.lru.Add(k, v)
	}
	return v, ok

}
