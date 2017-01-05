package Memory

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func evictedFuncForLRU(key, value interface{}) {
	fmt.Printf("[LRU] Key:%v Value:%v will evicted.\n", key, value)
}

func buildLRUCache(size int, loader LoaderFunc) MemoryPlugin {
	return New(size).
		LRU().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForLRU).
		Build()
}

func buildLoadingLRUCache(size int, loader LoaderFunc) MemoryPlugin {
	return New(size).
		LRU().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForLRU).
		Expiration(time.Second).
		Build()
}

func TestLRUGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000
	gc := buildLRUCache(size, loader)
	//set
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	//get
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLRUGetWithTimeout(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000
	gc := buildLoadingLRUCache(size, loader)
	//set
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	//get
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingLRUGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLRUCache(size, loader)
	//get
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingLRUGetWithTimeout(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLoadingLRUCache(size, loader)
	//get
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLRULength(t *testing.T) {
	assert := assert.New(t)

	gc := buildLRUCache(1000, loader)
	gc.Get("test1")
	gc.Get("test2")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(time.Second)

	length = gc.Len()
	assert.Equal(length, 2)

}

func TestLRULengthWithTimeout(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingLRUCache(1000, loader)
	gc.Get("test1")
	gc.Get("test2")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(time.Second)

	length = gc.Len()
	assert.Equal(length, 0)

}

func TestLRUEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLRUCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get(fmt.Sprintf("Key-%d", i))
		assert.Nil(err)
	}
}

func TestLRUEvictItemWithTimeout(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingLRUCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get(fmt.Sprintf("Key-%d", i))
		assert.Nil(err)
	}
}

func TestLRUGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		LRU().
		LoaderFunc(
			func(key interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond)
				return "value", nil
			}).
		Build()

	v, err := cache.GetIFPresent("key")
	assert.Equal(err, errors.New("Key is not find!"))

	time.Sleep(2 * time.Millisecond)

	v, err = cache.GetIFPresent("key")
	assert.Nil(err)

	assert.Equal(v, "value")
}

func TestLRUGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		LRU().
		Build()

	for i := 0; i < size; i++ {
		cache.Set(i, i*i)
	}

	m := cache.GetALL()
	for i := 0; i < size; i++ {
		v, ok := m[i]
		assert.True(ok)
		assert.Equal(v, i*i)
	}
	time.Sleep(time.Millisecond)

	cache.Set(size, size*size)
	m = cache.GetALL()

	assert.Equal(len(m), 1)

	v1, ok := m[size]
	assert.True(ok)
	assert.Equal(v1, size*size)
}
