package Memory

import (
	"time"

	log "github.com/cihub/seelog"
)

func New(size int) *CacheBuilder {
	if size <= 0 {
		log.Critical("Memory: size <= 0")
		panic("Memory: size <= 0")
	}

	return &CacheBuilder{
		tp:   SIMPLE,
		size: size,
	}
}

type CacheBuilder struct {
	tp          MODE // mode : simple \lru \ lfu \ arc
	size        int  // cache size > 0
	loaderFunc  *LoaderFunc
	evictedFunc *EvictedFunc
	addedFunc   *AddedFunc
	expiration  *time.Duration
}

// Set a loader function.
// loaderFunc: create a new value with this function if cached value is expired.
func (cb *CacheBuilder) LoaderFunc(loaderFunc LoaderFunc) *CacheBuilder {
	cb.loaderFunc = &loaderFunc
	return cb
}

func (cb *CacheBuilder) EvictType(tp MODE) *CacheBuilder {
	cb.tp = tp
	return cb
}

func (cb *CacheBuilder) Simple() *CacheBuilder {
	return cb.EvictType(SIMPLE)
}

func (cb *CacheBuilder) LRU() *CacheBuilder {
	return cb.EvictType(LRU)
}

func (cb *CacheBuilder) LFU() *CacheBuilder {
	return cb.EvictType(LFU)
}

func (cb *CacheBuilder) ARC() *CacheBuilder {
	return cb.EvictType(ARC)
}

func (cb *CacheBuilder) EvictedFunc(evictedFunc EvictedFunc) *CacheBuilder {
	cb.evictedFunc = &evictedFunc
	return cb
}

func (cb *CacheBuilder) AddedFunc(addedFunc AddedFunc) *CacheBuilder {
	cb.addedFunc = &addedFunc
	return cb
}

func (cb *CacheBuilder) Expiration(expiration time.Duration) *CacheBuilder {
	cb.expiration = &expiration
	return cb
}

func (cb *CacheBuilder) Build() MemoryPlugin {
	if HasRegister(cb.tp) {
		return PluginInstance(cb)
	} else {
		log.Critical("Memory: Unknown memory plugin: " + cb.tp)
		panic("Memory: Unknown memory plugin: " + cb.tp)
	}
}
