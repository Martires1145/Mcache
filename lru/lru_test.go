package lru

import (
	"testing"
)

type String string

func (s String) Len() int64 {
	return int64(len(s))
}

func TestCacheAddAndGet(t *testing.T) {
	cache := New(0, nil)

	cache.Add("1", String("hello"))
	value, ok := cache.Get("1")
	if !ok || string(value.(String)) != "hello" {
		t.Fatalf("cache hit key=1 failed")
	}

	cache.Add("2", String("world"))
	value, ok = cache.Get("2")
	if !ok || string(value.(String)) != "world" {
		t.Fatalf("cache hit key=2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	cache := New(5, nil)

	cache.Add("1", String("hello"))
	cache.Add("2", String("world"))

	if _, ok := cache.Get("1"); ok {
		t.Fatalf("cache remove key=1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	msg := ""
	cache := New(5, func(key string, value Value) {
		msg = "removed"
	})

	cache.Add("1", String("hello"))
	cache.Add("2", String("world"))

	if msg != "removed" {
		t.Fatalf("cache run onevicted function failed")
	}
}