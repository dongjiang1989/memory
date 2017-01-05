package Memory

import (
	"errors"
	"sync"
)

// callback type
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	plugin MemoryPlugin
	mu     sync.Mutex            // protects m
	m      map[interface{}]*call // lazily initialized
}

// Do executes and returns the results of the given function, only one execution is runing for a given key at a time
// If a duplicate comes in, the duplicate caller waits for the original to complete and receives the same results.
func (g *Group) Do(key interface{}, fn func() (interface{}, error), isWait bool) (interface{}, bool, error) {
	g.mu.Lock()
	v, err := g.plugin.get(key, true)
	if err == nil {
		g.mu.Unlock()
		return v, false, nil
	}
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		if !isWait {
			return nil, false, errors.New("Key is not find!")
		}
		c.wg.Wait()
		return c.val, false, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	if !isWait {
		go g.call(c, key, fn)
		return nil, false, errors.New("Key is not find!")
	}
	v, err = g.call(c, key, fn)
	return v, true, err
}

func (g *Group) call(c *call, key interface{}, fn func() (interface{}, error)) (interface{}, error) {
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
