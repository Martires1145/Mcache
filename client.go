package Mcache

import (
	"fmt"
	"github.com/Martires1145/Mcache/singleflight"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

type Client struct {
	mu        sync.Mutex
	getter    Getter
	mainCache cache
	loader    *singleflight.Group
	uMap      map[string]any
}

func New(maxBytes int64, getter Getter) *Client {
	return &Client{
		getter:    getter,
		mainCache: cache{cacheBytes: maxBytes},
		loader:    &singleflight.Group{},
		uMap:      make(map[string]any),
	}
}

func (c *Client) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("empty key")
	}

	value, err := c.loader.Do(key, func() (any, error) {
		c.mu.Lock()
		if _, ok := c.uMap[key]; ok {
			delete(c.uMap, key)
			c.mu.Unlock()
			return c.getFromGetter(key)
		}
		c.mu.Unlock()
		if value, ok := c.mainCache.get(key); ok {
			log.Printf("[MCache] hit key:%s\n", key)
			return value, nil
		}

		return c.getFromGetter(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return value.(ByteView), err
}

func (c *Client) Update(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.uMap[key] = struct{}{}
}

func (c *Client) getFromGetter(key string) (value ByteView, err error) {
	v, err := c.getter.Get(key)
	if err != nil {
		return
	}

	bytes := ByteView{cloneBytes(v)}
	c.loadNewCache(key, bytes)
	return bytes, nil
}

func (c *Client) loadNewCache(key string, value ByteView) {
	c.mainCache.add(key, value)
}
