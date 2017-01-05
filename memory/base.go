package Memory

import (
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

type Base struct {
	size        int // cache size > 0
	loaderFunc  *LoaderFunc
	evictedFunc *EvictedFunc
	addedFunc   *AddedFunc
	expiration  *time.Duration
	mu          sync.RWMutex
	loadGroup   Group
	*stats
}

type LoaderFunc func(interface{}) (interface{}, error)

type EvictedFunc func(interface{}, interface{})

type AddedFunc func(interface{}, interface{})

func buildCache(c *Base, cb *CacheBuilder) {
	c.size = cb.size
	c.loaderFunc = cb.loaderFunc
	c.expiration = cb.expiration
	c.addedFunc = cb.addedFunc
	c.evictedFunc = cb.evictedFunc
	c.stats = &stats{}
}

// load a new value using by specified key.
func (c *Base) load(key interface{}, cb func(interface{}, error) (interface{}, error), isWait bool) (interface{}, bool, error) {
	v, called, err := c.loadGroup.Do(key, func() (interface{}, error) {
		return cb((*c.loaderFunc)(key))
	}, isWait)

	if err != nil {
		log.Error(err)
		return nil, called, err
	}
	return v, called, nil
}
