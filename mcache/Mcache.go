package mcache

import (
	"MayCache/mcache/bloomfilter"
	"MayCache/mcache/singleflight"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name         string
	getter       Getter
	mainCache    cache
	peers        PeerPicker
	once         sync.Once
	singleFlight *singleflight.Group
	filter       *bloomfilter.BloomFilter
}

var (
	groups = map[string]*Group{}
	mu     sync.RWMutex
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	g := &Group{
		name:         name,
		getter:       getter,
		mainCache:    cache{cacheBytes: cacheBytes},
		singleFlight: &singleflight.Group{},
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func (g *Group) AddFilter() {
	g.filter = bloomfilter.New(10000000, 10)
}

func (g *Group) RegisterPeerPicker(peerPicker PeerPicker) {
	g.once.Do(func() {
		g.peers = peerPicker
	})
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if g.filter != nil && g.filter.MightContain(key) {
		return ByteView{}, fmt.Errorf("no such data")
	}
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[MayCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	value, err := g.singleFlight.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			peerGetter, ok := g.peers.PickPeer(key)
			if ok {
				value, err := g.getFromPeer(peerGetter, key)
				if err == nil {
					return value, nil
				}
			}
		}
		return g.getLocally(key)
	})

	if err != nil {
		return ByteView{}, err
	}
	return value.(ByteView), nil
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, err
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{bytes}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	if g.filter != nil {
		g.filter.Add(key)
	}
	g.mainCache.add(key, value)
}
