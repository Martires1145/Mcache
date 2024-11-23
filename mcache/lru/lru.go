package lru

import "container/list"

type Cache struct {
	cache     map[string]*list.Element
	list      *list.List
	maxBytes  int64
	nowBytes  int64
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		cache:     map[string]*list.Element{},
		list:      list.New(),
		maxBytes:  maxBytes,
		nowBytes:  0,
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, ok
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.list.Back()
	if ele != nil {
		c.list.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.list.PushFront(&entry{key: key, value: value})
		c.nowBytes += int64(len(key)) + int64(value.Len())
		c.cache[key] = ele
	}

	for ; c.maxBytes != 0 && c.nowBytes > c.maxBytes; c.RemoveOldest() {
	}
}

func (c *Cache) Len() int {
	return c.list.Len()
}
