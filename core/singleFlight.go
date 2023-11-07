// Package core @Author Bing
// @Date 2023/11/6 11:19:00
// @Desc
package core

import (
	"sync"
)

type call struct {
	wg    sync.WaitGroup
	value []byte
	error error
}

type SingleFlight struct {
	lock   sync.RWMutex
	buffer map[string]*call
}

func (s *SingleFlight) singleCall(g, k string, f func() ([]byte, error)) ([]byte, error) {
	s.lock.RLock()
	if c, ok := s.buffer[k]; ok {
		c.wg.Wait()
		return c.value, nil
	}
	s.lock.RUnlock()

	c := new(call)
	c.wg.Add(1)
	s.buffer[k] = c
	c.value, c.error = f()
	//time.Sleep(time.Second * 1) // todo: open sleep single flight cache
	c.wg.Done()

	delete(s.buffer, k)
	return c.value, c.error
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		lock:   sync.RWMutex{},
		buffer: make(map[string]*call),
	}
}
