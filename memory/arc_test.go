package Memory

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func buildARCache(size int) MemoryPlugin {
	return New(size).
		ARC().
		EvictedFunc(evictedFuncForARC).
		Build()
}

func buildLoadingARCache(size int) MemoryPlugin {
	return New(size).
		ARC().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForARC).
		Build()
}

func buildLoadingARCacheWithExpiration(size int, ep time.Duration) MemoryPlugin {
	return New(size).
		ARC().
		Expiration(ep).
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForARC).
		Build()
}

func evictedFuncForARC(key, value interface{}) {
	fmt.Printf("[ARC] Key:%v Value:%v will evicted.\n", key, value)
}

func TestARCGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildARCache(size)

	//set
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	//get
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestARCGetBig(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildARCache(size)

	//set
	for i := 0; i < size+10; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		assert.Nil(err)
		gc.Set(key, value)
		gc.Get(fmt.Sprintf("Key-1"))
	}

	//get
	assert.Equal(gc.Len(), size)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("Key-%d", i+size)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}

	for i := 11; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingARCGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLoadingARCache(size)

	//set
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		assert.Nil(err)
		gc.Set(key, value)
	}

	//get
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingARCGetWithExpiration(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLoadingARCacheWithExpiration(size, 1*time.Nanosecond)

	//set
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		assert.Nil(err)
		gc.Set(key, value)
	}

	//get
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestARCLength(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingARCacheWithExpiration(2, time.Millisecond)
	gc.Get("test1")
	gc.Get("test2")
	gc.Get("test3")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(time.Millisecond)
	gc.Get("test4")
	length = gc.Len()
	assert.Equal(length, 1)
}

func TestARCEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingARCache(cacheSize)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get(fmt.Sprintf("Key-%d", i))
		assert.Nil(err)
	}
}

func TestARCGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		ARC().
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

func TestARCGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		ARC().
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
