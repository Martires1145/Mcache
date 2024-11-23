package consistenthash

import (
	"hash/crc32"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     SortSlice
	hashMap  map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys.Add(hash)
			m.hashMap[hash] = key
		}
	}
}

func (m *Map) Get(key string) string {
	if m.keys.Len() == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := m.keys.Index(hash)

	return m.hashMap[m.keys.Get(idx%m.keys.Len())]
}

func (m *Map) Delete(key string) {
	hash := int(m.hash([]byte(key)))
	idx := m.keys.Index(hash)
	m.keys.Delete(idx)
	delete(m.hashMap, hash)
}
