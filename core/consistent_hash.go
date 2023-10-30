// Package core @Author Bing
// @Date 2023/10/27 17:46:00
// @Desc
package core

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hashMap  map[int]string
	keys     []int
	replicas int
	hash     Hash
}

func NewMap(replicas int, hash Hash) *Map {
	m := &Map{
		hashMap:  make(map[int]string),
		keys:     make([]int, 0),
		replicas: replicas,
		hash:     hash,
	}
	if hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[index%len(m.keys)]]
}
