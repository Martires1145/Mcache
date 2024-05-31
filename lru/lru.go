package lru

import "container/list"

type Cache struct {
	cache map[string]*list.Element
	list  *list.List
	// maxBytes less than 1 indicates an unlimited size
	maxBytes  int64
	usedBytes int64
	// optional and executed when an entry is purged
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int64
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		cache:     make(map[string]*list.Element),
		list:      list.New(),
		maxBytes:  maxBytes,
		usedBytes: 0,
		onEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	ele := c.list.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		c.list.Remove(ele)
		delete(c.cache, kv.key)
		c.usedBytes -= kv.value.Len()
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += value.Len() - kv.value.Len()
		kv.value = value
	} else {
		ele := c.list.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.usedBytes += value.Len()
	}

	for c.maxBytes > 0 && c.usedBytes > c.maxBytes {
		c.RemoveOldest()
	}
}
