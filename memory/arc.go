package Memory

import (
	"common/memory/util"
	"common/utils"
	"errors"
	"time"
)

// NewARCPlugin returns a new plugin.
func NewARCPlugin(cb *CacheBuilder) MemoryPlugin {
	c := &ARCPlugin{}

	buildCache(&c.Base, cb)
	c.init()
	c.loadGroup.plugin = c

	return c
}

type ARCPlugin struct {
	Base
	items map[interface{}]*PUtil.ArcItem

	part int
	t1   *PUtil.ArcList
	t2   *PUtil.ArcList
	b1   *PUtil.ArcList
	b2   *PUtil.ArcList
}

func (c *ARCPlugin) init() {
	c.items = make(map[interface{}]*PUtil.ArcItem)
	c.t1 = PUtil.NewARCList()
	c.t2 = PUtil.NewARCList()
	c.b1 = PUtil.NewARCList()
	c.b2 = PUtil.NewARCList()
}

func (c *ARCPlugin) replace(key interface{}) {
	var old interface{}
	if (c.t1.Len() > 0 && c.b2.Has(key) && c.t1.Len() == c.part) || (c.t1.Len() > c.part) {
		old = c.t1.RemoveTail()
		c.b1.PushFront(old)
	} else if c.t2.Len() > 0 {
		old = c.t2.RemoveTail()
		c.b2.PushFront(old)
	} else {
		return
	}
	item, ok := c.items[old]
	if ok {
		delete(c.items, old)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(item.Key, item.Value)
		}
	}
}

func (c *ARCPlugin) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *ARCPlugin) set(key, value interface{}) (interface{}, error) {
	item, ok := c.items[key]
	if ok {
		item.Value = value
	} else {
		item = &PUtil.ArcItem{
			Key:   key,
			Value: value,
		}
		c.items[key] = item
	}

	if c.expiration != nil {
		t := time.Now().Add(*c.expiration)
		item.Expiration = &t
	}

	if elt := c.b1.Lookup(key); elt != nil {
		c.part = minInt(c.size, c.part+maxInt(c.b2.Len()/c.b1.Len(), 1))
		c.replace(key)
		c.b1.Remove(key, elt)
		c.t2.PushFront(key)
		return item, nil
	}

	if elt := c.b2.Lookup(key); elt != nil {
		c.part = maxInt(0, c.part-maxInt(c.b1.Len()/c.b2.Len(), 1))
		c.replace(key)
		c.b2.Remove(key, elt)
		c.t2.PushFront(key)
		return item, nil
	}

	if c.t1.Len()+c.b1.Len() == c.size {
		if c.t1.Len() < c.size {
			c.b1.RemoveTail()
			c.replace(key)
		} else {
			pop := c.t1.RemoveTail()
			item, ok := c.items[pop]
			if ok {
				delete(c.items, pop)
				if c.evictedFunc != nil {
					(*c.evictedFunc)(item.Key, item.Value)
				}
			}
		}
	} else {
		total := c.t1.Len() + c.b1.Len() + c.t2.Len() + c.b2.Len()
		if total >= c.size {
			if total == (2 * c.size) {
				c.b2.RemoveTail()
			}
			c.replace(key)
		}
	}

	c.t1.PushFront(key)

	if c.addedFunc != nil {
		(*c.addedFunc)(key, value)
	}

	return item, nil
}

// Get a value from cache pool using key if it exists.
// If not exists and it has LoaderFunc, it will generate the value using you have specified LoaderFunc method returns value.
func (c *ARCPlugin) Get(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, true)
	}
	return v, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key, returns KeyNotFoundError.
// And send a request which refresh value for specified key if cache object has LoaderFunc.
func (c *ARCPlugin) GetIFPresent(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, false)
	}
	return v, nil
}

func (c *ARCPlugin) get(key interface{}, onLoad bool) (interface{}, error) {
	rl := false
	c.mu.RLock()
	if elt := c.t1.Lookup(key); elt != nil {
		c.mu.RUnlock()
		rl = true
		c.mu.Lock()
		c.t1.Remove(key, elt)
		item := c.items[key]
		if !item.IsExpired(nil) {
			c.t2.PushFront(key)
			c.mu.Unlock()
			if !onLoad {
				c.stats.IncrHitCount()
			}
			return item, nil
		}
		c.b2.PushFront(key)
		delete(c.items, key)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, elt.Value)
		}
		c.mu.Unlock()
	}
	if elt := c.t2.Lookup(key); elt != nil {
		c.mu.RUnlock()
		rl = true
		c.mu.Lock()
		item := c.items[key]
		if !item.IsExpired(nil) {
			c.t2.MoveToFront(elt)
			c.mu.Unlock()
			if !onLoad {
				c.stats.IncrHitCount()
			}
			return item, nil
		}
		c.t2.Remove(key, elt)
		c.b2.PushFront(key)
		delete(c.items, key)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, elt.Value)
		}
		c.mu.Unlock()
	}

	if !rl {
		c.mu.RUnlock()
	}
	if !onLoad {
		c.stats.IncrMissCount()
	}
	return nil, errors.New("key is not find!")
}

func (c *ARCPlugin) getValue(key interface{}) (interface{}, error) {
	it, err := c.get(key, false)
	if err != nil {
		return nil, err
	}
	return it.(*PUtil.ArcItem).Value, nil
}

func (c *ARCPlugin) getWithLoader(key interface{}, isWait bool) (interface{}, error) {
	if c.loaderFunc == nil {
		return nil, errors.New("key is not find!")
	}
	item, _, err := c.load(key, func(v interface{}, e error) (interface{}, error) {
		if e == nil {
			c.mu.Lock()
			defer c.mu.Unlock()
			return c.set(key, v)
		}
		return nil, e
	}, isWait)
	if err != nil {
		return nil, err
	}
	return item.(*PUtil.ArcItem).Value, nil
}

// Remove removes the provided key from the cache.
func (c *ARCPlugin) Remove(key interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.remove(key)
}

func (c *ARCPlugin) remove(key interface{}) bool {
	if elt := c.t1.Lookup(key); elt != nil {
		v := elt.Value.(*PUtil.ArcItem).Value
		c.t1.Remove(key, elt)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, v)
		}
		return true
	}

	if elt := c.t2.Lookup(key); elt != nil {
		v := elt.Value.(*PUtil.ArcItem).Value
		c.t2.Remove(key, elt)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, v)
		}
		return true
	}

	return false
}

func (c *ARCPlugin) keys() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]interface{}, len(c.items))
	var i = 0
	for k := range c.items {
		keys[i] = k
		i++
	}
	return keys
}

// Keys returns a slice of the keys in the cache.
func (c *ARCPlugin) Keys() []interface{} {
	keys := []interface{}{}
	for _, k := range c.keys() {
		_, err := c.GetIFPresent(k)
		if err == nil {
			keys = append(keys, k)
		}
	}
	return keys
}

// Returns all key-value pairs in the cache.
func (c *ARCPlugin) GetALL() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, k := range c.keys() {
		v, err := c.GetIFPresent(k)
		if err == nil {
			m[k] = v
		}
	}
	return m
}

// Len returns the number of items in the cache.
func (c *ARCPlugin) Len() int {
	return len(c.GetALL())
}

// Purge is used to completely clear the cache
func (c *ARCPlugin) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.init()
}

func (c *ARCPlugin) HasKey(key interface{}) bool {
	return Utils.InSliceIface(key, c.Keys())
}

//init
func init() {
	Register(ARC, NewARCPlugin)
}
