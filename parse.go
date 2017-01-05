package Parse

import (
	"common/memory"
	"errors"
	"time"
)

type Parse struct {
	Syb Memory.MemoryPlugin //symbol
	//cmp *Memory.CacheBuilder // class map
}

func (p *Parse) NewInit(bc int, lc int, path string, size int, mode Memory.MODE, loaderFunc Memory.LoaderFunc, addedFunc Memory.AddedFunc, evictedFunc Memory.EvictedFunc, timeout int) (*Parse, error) {
	cfg.BlockCount = bc
	cfg.LineCount = lc
	cfg.Path = path

	sybObj := Memory.
		New(size).
		Expiration(time.Duration(timeout) * time.Second).
		EvictType(mode).
		LoaderFunc(loaderFunc).
		AddedFunc(addedFunc).
		EvictedFunc(evictedFunc).
		Build()

	if sybObj == nil {
		return nil, errors.New("New Memory is cache object is nil!")
	}
	p.Syb = sybObj

	return p, nil
}
