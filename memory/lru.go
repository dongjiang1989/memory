package Memory

import (
	"common/memory/util"
	"common/utils"
	"container/list"
	"errors"
	"time"
)

// NewLRUPlugin returns a new plugin.
func NewLRUPlugin(cb *CacheBuilder) MemoryPlugin {
	c := &LRUPlugin{}
	buildCache(&c.Base, cb)

	c.init()
	c.loadGroup.plugin = c
	return c
}

// Discards the least recently used items first.
type LRUPlugin struct {
	Base
	items     map[interface{}]*list.Element
	evictList *list.List
}

func (c *LRUPlugin) init() {
	c.evictList = list.New()
	c.items = make(map[interface{}]*list.Element, c.size+1)
}

// set a new key-value pair
func (c *LRUPlugin) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *LRUPlugin) set(key, value interface{}) (interface{}, error) {
	// Check for existing item
	var item *PUtil.LruItem
	if it, ok := c.items[key]; ok {
		c.evictList.MoveToFront(it)
		item = it.Value.(*PUtil.LruItem)
		item.Value = value
	} else {
		// Verify size not exceeded
		if c.evictList.Len() >= c.size {
			c.evict(1)
		}
		item = &PUtil.LruItem{
			Key:   key,
			Value: value,
		}
		c.items[key] = c.evictList.PushFront(item)
	}

	if c.expiration != nil {
		t := time.Now().Add(*c.expiration)
		item.Expiration = &t
	}

	if c.addedFunc != nil {
		(*c.addedFunc)(key, value)
	}

	return item, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key and has LoaderFunc,
// generate a value using `LoaderFunc` method returns value.
func (c *LRUPlugin) Get(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, true)
	}
	return v, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key, returns KeyNotFoundError.
// And send a request which refresh value for specified key if cache object has LoaderFunc.
func (c *LRUPlugin) GetIFPresent(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, false)
	}
	return v, nil
}

func (c *LRUPlugin) get(key interface{}, onLoad bool) (interface{}, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if ok {
		it := item.Value.(*PUtil.LruItem)
		if !it.IsExpired(nil) {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.evictList.MoveToFront(item)
			if !onLoad {
				c.stats.IncrHitCount()
			}
			return it, nil
		}
		c.mu.Lock()
		c.removeElement(item)
		c.mu.Unlock()
	}
	if !onLoad {
		c.stats.IncrMissCount()
	}
	return nil, errors.New("key is not find!")
}

func (c *LRUPlugin) getValue(key interface{}) (interface{}, error) {
	it, err := c.get(key, false)
	if err != nil {
		return nil, err
	}
	return it.(*PUtil.LruItem).Value, nil
}

func (c *LRUPlugin) getWithLoader(key interface{}, isWait bool) (interface{}, error) {
	if c.loaderFunc == nil {
		return nil, errors.New("key is not find!")
	}
	it, _, err := c.load(key, func(v interface{}, e error) (interface{}, error) {
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
	return it.(*PUtil.LruItem).Value, nil
}

// evict removes the oldest item from the cache.
func (c *LRUPlugin) evict(count int) {
	for i := 0; i < count; i++ {
		ent := c.evictList.Back()
		if ent == nil {
			return
		} else {
			c.removeElement(ent)
		}
	}
}

// Removes the provided key from the cache.
func (c *LRUPlugin) Remove(key interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.remove(key)
}

func (c *LRUPlugin) remove(key interface{}) bool {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

func (c *LRUPlugin) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	entry := e.Value.(*PUtil.LruItem)
	delete(c.items, entry.Key)
	if c.evictedFunc != nil {
		entry := e.Value.(*PUtil.LruItem)
		(*c.evictedFunc)(entry.Key, entry.Value)
	}
}

func (c *LRUPlugin) keys() []interface{} {
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

// Returns a slice of the keys in the cache.
func (c *LRUPlugin) Keys() []interface{} {
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
func (c *LRUPlugin) GetALL() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, k := range c.keys() {
		v, err := c.GetIFPresent(k)
		if err == nil {
			m[k] = v
		}
	}
	return m
}

// Returns the number of items in the cache.
func (c *LRUPlugin) Len() int {
	return len(c.GetALL())
}

// Completely clear the cache
func (c *LRUPlugin) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.init()
}

func (c *LRUPlugin) HasKey(key interface{}) bool {
	return Utils.InSliceIface(key, c.Keys())
}

//init
func init() {
	Register(LRU, NewLFUPlugin)
}
