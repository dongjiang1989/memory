package Memory

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoaderFuncLRU(t *testing.T) {
	assert := assert.New(t)

	size := 2

	var testCaches = []*CacheBuilder{
		New(size).Simple(),
		New(size).LRU(),
		New(size).LFU(),
		New(size).ARC(),
	}

	for _, builder := range testCaches {
		var testCounter int64
		counter := 1000
		cache := builder.
			LoaderFunc(func(key interface{}) (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				log.Println("dongjiang LoaderFunc==", "key:", key)
				return atomic.AddInt64(&testCounter, 1), nil
			}).
			AddedFunc(func(key, value interface{}) {
				log.Println("dongjiang AddedFunc==", "key:", key, "value:", value)
			}).
			EvictedFunc(func(key, value interface{}) {
				log.Println("dongjiang EvictedFunc==", "key:", key, "value:", value)
			}).Build()

		var wg sync.WaitGroup
		for i := 0; i < counter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cache.Get(0)
				assert.Nil(err)
			}()
		}
		wg.Wait()

		assert.Equal(testCounter, int64(1))
	}

}

func TestLoaderFuncSimple(t *testing.T) {
	assert := assert.New(t)

	size := 2

	var testCaches = []*CacheBuilder{
		New(size).Simple(),
		New(size).LRU(),
		New(size).LFU(),
		New(size).ARC(),
	}

	for _, builder := range testCaches {
		var testCounter int64
		counter := 1000
		cache := builder.
			LoaderFunc(func(key interface{}) (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				return atomic.AddInt64(&testCounter, 1), nil
			}).
			AddedFunc(func(key, value interface{}) {
				log.Println("dongjiang AddedFunc==", "key:", key, "value:", value)
			}).
			EvictedFunc(func(key, value interface{}) {
				log.Println("dongjiang EvictedFunc==", "key:", key, "value:", value)
			}).Simple().Build()

		var wg sync.WaitGroup
		for i := 0; i < counter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cache.Get(0)
				assert.Nil(err)
			}()
		}
		wg.Wait()

		assert.Equal(testCounter, int64(1))
	}

}

func TestLoaderFuncLFU(t *testing.T) {
	assert := assert.New(t)

	size := 2

	var testCaches = []*CacheBuilder{
		New(size).Simple(),
		New(size).LRU(),
		New(size).LFU(),
		New(size).ARC(),
	}

	for _, builder := range testCaches {
		var testCounter int64
		counter := 1000
		cache := builder.
			LoaderFunc(func(key interface{}) (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				return atomic.AddInt64(&testCounter, 1), nil
			}).
			AddedFunc(func(key, value interface{}) {
				log.Println("dongjiang AddedFunc==", "key:", key, "value:", value)
			}).
			EvictedFunc(func(key, value interface{}) {
				log.Println("dongjiang EvictedFunc==", "key:", key, "value:", value)
			}).LFU().Build()

		var wg sync.WaitGroup
		for i := 0; i < counter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cache.Get(0)
				assert.Nil(err)
			}()
		}
		wg.Wait()

		assert.Equal(testCounter, int64(1))
	}

}

func TestLoaderFuncARC(t *testing.T) {
	assert := assert.New(t)

	size := 2

	var testCaches = []*CacheBuilder{
		New(size).Simple(),
		New(size).LRU(),
		New(size).LFU(),
		New(size).ARC(),
	}

	for _, builder := range testCaches {
		var testCounter int64
		counter := 1000
		cache := builder.
			LoaderFunc(func(key interface{}) (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				return atomic.AddInt64(&testCounter, 1), nil
			}).
			AddedFunc(func(key, value interface{}) {
				log.Println("dongjiang AddedFunc==", "key:", key, "value:", value)
			}).
			EvictedFunc(func(key, value interface{}) {
				log.Println("dongjiang EvictedFunc==", "key:", key, "value:", value)
			}).ARC().Build()

		var wg sync.WaitGroup
		for i := 0; i < counter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cache.Get(0)
				assert.Nil(err)
			}()
		}
		wg.Wait()

		assert.Equal(testCounter, int64(1))
	}

}
