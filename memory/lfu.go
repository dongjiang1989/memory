package Memory

import (
	"common/memory/util"
	"common/utils"
	"container/list"
	"errors"
	"time"
)

// NewLFUPlugin returns a new plugin.
func NewLFUPlugin(cb *CacheBuilder) MemoryPlugin {
	c := &LFUPlugin{}
	buildCache(&c.Base, cb)

	c.init()
	c.loadGroup.plugin = c
	return c
}

type freqEntry struct {
	freq  uint
	items map[*PUtil.LfuItem]byte
}

// Discards the least frequently used items first.
type LFUPlugin struct {
	Base
	items    map[interface{}]*PUtil.LfuItem
	freqList *list.List // list for freqEntry
}

func (c *LFUPlugin) init() {
	c.freqList = list.New()
	c.items = make(map[interface{}]*PUtil.LfuItem, c.size+1)
	c.freqList.PushFront(&freqEntry{
		freq:  0,
		items: make(map[*PUtil.LfuItem]byte),
	})
}

// set a new key-value pair
func (c *LFUPlugin) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *LFUPlugin) set(key, value interface{}) (interface{}, error) {
	// Check for existing item
	item, ok := c.items[key]
	if ok {
		item.Value = value
	} else {
		// Verify size not excepted
		if len(c.items) >= c.size {
			c.evict(1)
		}
		item = &PUtil.LfuItem{
			Key:         key,
			Value:       value,
			FreqElement: nil,
		}
		el := c.freqList.Front()
		fe := el.Value.(*freqEntry)
		fe.items[item] = 1

		item.FreqElement = el
		c.items[key] = item
	}

	if c.expiration != nil {
		t := time.Now().Add(*c.expiration)
		item.Expiration = &t
	}

	// run addedFunc
	if c.addedFunc != nil {
		(*c.addedFunc)(key, value)
	}

	return item, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key and has LoaderFunc,
// generate a value using `LoaderFunc` method returns value.
func (c *LFUPlugin) Get(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, true)
	}
	return v, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key, returns KeyNotFoundError.
// And send a request which refresh value for specified key if cache object has LoaderFunc.
func (c *LFUPlugin) GetIFPresent(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, false)
	}
	return v, nil
}

func (c *LFUPlugin) get(key interface{}, onLoad bool) (interface{}, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if ok {
		if !item.IsExpired(nil) {
			c.mu.Lock()
			c.increment(item)
			c.mu.Unlock()
			if !onLoad {
				c.stats.IncrHitCount()
			}
			return item, nil
		}
		c.mu.Lock()
		c.removeItem(item)
		c.mu.Unlock()
	}
	if !onLoad {
		c.stats.IncrMissCount()
	}
	return nil, errors.New("key is not find!")
}

func (c *LFUPlugin) getValue(key interface{}) (interface{}, error) {
	it, err := c.get(key, false)
	if err != nil {
		return nil, err
	}
	return it.(*PUtil.LfuItem).Value, nil
}

func (c *LFUPlugin) getWithLoader(key interface{}, isWait bool) (interface{}, error) {
	if c.loaderFunc == nil {
		return nil, errors.New("key is not find!")
	}
	it, called, err := c.load(key, func(v interface{}, e error) (interface{}, error) {
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
	li := it.(*PUtil.LfuItem)
	if !called {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.increment(li)
	}
	return li.Value, nil
}

func (c *LFUPlugin) increment(item *PUtil.LfuItem) {
	currentFreqElement := item.FreqElement
	currentFreqEntry := currentFreqElement.Value.(*freqEntry)
	nextFreq := currentFreqEntry.freq + 1
	delete(currentFreqEntry.items, item)

	nextFreqElement := currentFreqElement.Next()
	if nextFreqElement == nil {
		nextFreqElement = c.freqList.InsertAfter(&freqEntry{
			freq:  nextFreq,
			items: make(map[*PUtil.LfuItem]byte),
		}, currentFreqElement)
	}
	nextFreqElement.Value.(*freqEntry).items[item] = 1
	item.FreqElement = nextFreqElement
}

// evict removes the least frequence item from the cache.
func (c *LFUPlugin) evict(count int) {
	entry := c.freqList.Front()
	for i := 0; i < count; {
		if entry == nil {
			return
		} else {
			for item, _ := range entry.Value.(*freqEntry).items {
				if i >= count {
					return
				}
				c.removeItem(item)
				i++
			}
			entry = entry.Next()
		}
	}
}

// Removes the provided key from the cache.
func (c *LFUPlugin) Remove(key interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.remove(key)
}

func (c *LFUPlugin) remove(key interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.removeItem(item)
		return true
	}
	return false
}

// removeElement is used to remove a given list element from the cache
func (c *LFUPlugin) removeItem(item *PUtil.LfuItem) {
	delete(c.items, item.Key)
	delete(item.FreqElement.Value.(*freqEntry).items, item)
	if c.evictedFunc != nil {
		(*c.evictedFunc)(item.Key, item.Value)
	}
}

func (c *LFUPlugin) keys() []interface{} {
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
func (c *LFUPlugin) Keys() []interface{} {
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
func (c *LFUPlugin) GetALL() map[interface{}]interface{} {
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
func (c *LFUPlugin) Len() int {
	return len(c.GetALL())
}

// Completely clear the cache
func (c *LFUPlugin) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.init()
}

func (c *LFUPlugin) HasKey(key interface{}) bool {
	return Utils.InSliceIface(key, c.Keys())
}

//init
func init() {
	Register(LFU, NewLFUPlugin)
}
