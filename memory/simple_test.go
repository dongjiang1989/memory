package Memory

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func loader(key interface{}) (interface{}, error) {
	return fmt.Sprintf("valueFor%s", key), nil
}

func buildSimpleCache(size int) MemoryPlugin {
	return New(size).
		Simple().
		EvictedFunc(evictedFuncForSimple).
		Build()
}

func buildLoadingSimpleCache(size int, loader LoaderFunc) MemoryPlugin {
	return New(size).
		LoaderFunc(loader).
		Simple().
		EvictedFunc(evictedFuncForSimple).
		Build()
}

func evictedFuncForSimple(key, value interface{}) {
	fmt.Printf("[Simple] Key:%v Value:%v will evicted.\n", key, value)
}

func TestSimpleGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildSimpleCache(size)

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

func TestSimpleGetBig(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildSimpleCache(size)

	//set
	for i := 0; i < size+10; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	//get
	assert.Equal(gc.Len(), size)
}

func TestLoadingSimpleGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLoadingSimpleCache(size, loader)

	//get
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestSimpleLength(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingSimpleCache(1000, loader)
	gc.Get("test1")
	gc.Get("test2")

	length := gc.Len()
	expectedLength := 2
	assert.Equal(length, expectedLength)
	log.Println("dongjiang223", gc.GetALL())
}

func TestSimpleLength2(t *testing.T) {
	assert := assert.New(t)

	gc := buildSimpleCache(1000)
	gc.Get("test1")
	gc.Get("test2")

	length := gc.Len()
	expectedLength := 0
	assert.Equal(length, expectedLength)
	log.Println("dongjiang123", gc.GetALL())
}

func TestSimpleEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingSimpleCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get(fmt.Sprintf("Key-%d", i))
		assert.Nil(err)
	}
}

func TestSimpleGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		EvictType(SIMPLE).
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

func TestSimpleGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		EvictType(SIMPLE).
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
