package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, f func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	g.m[key] = &call{}
	g.m[key].wg.Add(1)
	g.m[key].val, g.m[key].err = f()
	g.m[key].wg.Done()
	g.mu.Unlock()
	return g.m[key].val, g.m[key].err
}
