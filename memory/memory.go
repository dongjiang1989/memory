package Memory

import (
	"log"
)

type MemoryPlugin interface {
	// set key-value
	Set(interface{}, interface{})
	// get by key
	Get(interface{}) (interface{}, error)
	GetIFPresent(interface{}) (interface{}, error)
	GetALL() map[interface{}]interface{}
	//private func: get key by serialize
	get(interface{}, bool) (interface{}, error)
	// delete for key
	Remove(interface{}) bool
	// clear Plugin
	Purge()
	// get all keys
	Keys() []interface{}
	// get count of block
	Len() int
	// cache has key
	HasKey(interface{}) bool

	//statistics hit count
	statsAccessor
}

// Instance is a function create a new Memory Instance
type Instance func(*CacheBuilder) MemoryPlugin

var adapters = make(map[MODE]Instance)

func Register(name MODE, adapter Instance) {
	if adapter == nil {
		panic("MemoryPlugin: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("MemoryPlugin: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func PluginInstance(cb *CacheBuilder) (adapter MemoryPlugin) {
	instanceFunc, ok := adapters[cb.tp]
	if !ok {
		log.Fatal("MemoryPlugin: unknown adapter name %q (forgot to import?)", string(cb.tp))
		return
	}
	adapter = instanceFunc(cb)
	return
}

func HasRegister(name MODE) bool {

	if _, ok := adapters[name]; ok {
		return true
	}
	log.Panic("Can not find adapter name:" + name)
	return false
}
