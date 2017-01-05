package Memory

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func evictedFuncForLFU(key, value interface{}) {
	fmt.Printf("[LFU] Key:%v Value:%v will evicted.\n", key, value)
}

func buildLFUCache(size int) MemoryPlugin {
	return New(size).
		LoaderFunc(loader).
		LFU().
		EvictedFunc(evictedFuncForLFU).
		Build()
}

func buildLoadingLFUCache(size int, loader LoaderFunc) MemoryPlugin {
	return New(size).
		LFU().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForLFU).
		Expiration(time.Second).
		Build()
}
func TestLFUGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000

	gc := buildLFUCache(size)
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

func TestLFUGetWithTimeout(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000

	gc := buildLoadingLFUCache(size, loader)
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

func TestLoadingLFUGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000

	gc := buildLoadingLFUCache(size, loader)

	//get
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLFULength(t *testing.T) {
	assert := assert.New(t)

	gc := buildLFUCache(1000)
	gc.Get("test1")
	gc.Get("test2")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(1 * time.Second)
	length = gc.Len()
	assert.Equal(length, 2)
}

func TestLFULengthWithTimeout(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingLFUCache(1000, loader)
	gc.Get("test1")
	gc.Get("test2")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(1 * time.Second)
	length = gc.Len()
	assert.Equal(length, 0)
}

func TestLFUEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLFUCache(cacheSize)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get(fmt.Sprintf("Key-%d", i))
		assert.Nil(err)
	}
}
func TestLFUEvictItemWithTimeout(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingLFUCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get(fmt.Sprintf("Key-%d", i))
		assert.Nil(err)
	}
}

func TestLFUGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		LFU().
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

func TestLFUGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		LFU().
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
