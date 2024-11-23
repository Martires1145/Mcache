package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 实现不论Do方法被调用多少次 fn函数都只会被调用一次的效果
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// 懒加载g.m
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 如果请求正在进行中,等待请求完成,直接返回结果
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 发起新的请求
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 完成请求解锁返回结果
	c.val, c.err = fn()
	c.wg.Done()

	// 更新g.m
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	// 返回结果
	return c.val, c.err
}
