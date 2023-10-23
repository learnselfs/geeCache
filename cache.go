// Package geeCache @Author Bing
// @Date 2023/10/20 16:37:00
// @Desc
package geeCache

import (
	"container/list"
	"sync"
)

type Caches[k comparable, v any] interface {
	New()
	Put(k, v)
	Get(k)
}

type Cache[K comparable, V any] struct {
	lock  sync.RWMutex
	size  int
	maps  map[K]*list.Element
	lists *list.List
}

type item[K comparable, V any] struct {
	k K
	v V
}

func (c *Cache[K, V]) New(len int) *Cache[K, V] {
	return &Cache[K, V]{
		size:  len,
		maps:  make(map[K]*list.Element),
		lists: list.New(),
	}
}

func (c *Cache[K, V]) Add(k K, v V) bool {
	defer c.lock.Unlock()
	c.lock.Lock()
	i := item[K, V]{k, v}

	if _, ok := c.maps[k]; !ok {
		return false
	}

	e := c.lists.PushFront(i)
	c.maps[k] = e
	c.outOfLen()
	return true
}

func (c *Cache[K, V]) Set(k K, v V) {
	defer c.lock.Unlock()
	c.lock.Lock()
	i := newItem(k, v)
	if v, ok := c.maps[k]; ok {
		c.lists.Remove(v)
	}
	e := c.lists.PushFront(i)
	c.maps[k] = e
	c.outOfLen()
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	v, ok := c.maps[k]
	if ok {
		c.lists.MoveToFront(v)
	}
	return v.Value.(item[K, V]).v, ok
}

func (c *Cache[K, V]) Peek(k K) (V, bool) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	v, ok := c.maps[k]
	return v.Value.(item[K, V]).v, ok
}

func (c *Cache[K, V]) Range(f func(k K, v V)) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	for k, v := range c.maps {
		f(k, v.Value.(item[K, V]).v)
	}
}

func newItem[K comparable, V any](k K, v V) item[K, V] {
	return item[K, V]{k, v}
}

func (c *Cache[K, V]) outOfLen() {
	if c.size >= c.lists.Len() {
		e := c.lists.Back()
		k := e.Value.(item[K, V]).k
		c.lists.Remove(e)
		delete(c.maps, k)
	}
}
