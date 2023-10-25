// Package core @Author Bing
// @Date 2023/10/20 16:37:00
// @Desc
package geeCache

import (
	"container/list"
	"fmt"
	"sync"
)

type Lru[K comparable, V any] struct {
	lock  sync.RWMutex
	size  int
	maps  map[K]*list.Element
	lists *list.List
}

type item[K comparable, V any] struct {
	k K
	v V
}

func NewLru[K comparable, V any](len int) *Lru[K, V] {
	return &Lru[K, V]{
		size:  len,
		maps:  make(map[K]*list.Element),
		lists: list.New(),
	}
}

func (l *Lru[K, V]) Add(k K, v V) bool {
	defer l.lock.Unlock()
	l.lock.Lock()
	i := item[K, V]{k, v}

	if _, ok := l.maps[k]; ok {
		return false
	}

	e := l.lists.PushFront(i)
	l.maps[k] = e
	l.outOfLen()
	return true
}

func (l *Lru[K, V]) Set(k K, v V) {
	defer l.lock.Unlock()
	l.lock.Lock()
	i := newItem(k, v)
	if v, ok := l.maps[k]; ok {
		l.lists.Remove(v)
	}
	e := l.lists.PushFront(i)
	l.maps[k] = e
	l.outOfLen()
}

func (l *Lru[K, V]) Get(k K) (V, bool) {
	defer l.lock.RUnlock()
	l.lock.RLock()
	_item, ok := l.maps[k]
	fmt.Println(l.maps)
	if ok {
		l.lists.MoveToFront(_item)
		return _item.Value.(item[K, V]).v, ok
	}
	var temp V
	return temp, ok
}

func (l *Lru[K, V]) Peek(k K) (V, bool) {
	defer l.lock.RUnlock()
	l.lock.RLock()
	v, ok := l.maps[k]
	return v.Value.(item[K, V]).v, ok
}

func (l *Lru[K, V]) Range(f func(k, v any) bool) {
	defer l.lock.RUnlock()
	l.lock.RLock()
	for k, v := range l.maps {
		if ok := f(k, v.Value.(item[K, V]).v); !ok {
			return
		}

	}
}

func newItem[K comparable, V any](k K, v V) item[K, V] {
	return item[K, V]{k, v}
}

func (l *Lru[K, V]) outOfLen() {
	if l.size < l.lists.Len() {
		e := l.lists.Back()
		k := e.Value.(item[K, V]).k
		l.lists.Remove(e)
		delete(l.maps, k)
	}
}
