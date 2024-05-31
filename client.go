package Mcache

import (
	"fmt"
	"log"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

type Client struct {
	getter    Getter
	mainCache cache
}

func New(maxBytes int64, getter Getter) *Client {
	return &Client{
		getter:    getter,
		mainCache: cache{cacheBytes: maxBytes},
	}
}

func (c *Client) Get(key string) (value ByteView, err error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("empty key")
	}

	if value, ok := c.mainCache.get(key); ok {
		log.Printf("[MCache] hit key:%s\n", key)
		return value, nil
	}

	return c.getFromGetter(key)
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
