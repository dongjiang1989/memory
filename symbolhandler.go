package Parse

import (
	"parse/symbol"
	"time"

	log "github.com/cihub/seelog"
)

type Cfg struct {
	Path       string
	BlockCount int
	LineCount  int
}

var cfg Cfg

// memory cache handle  ---- LoaderFunc
func LoaderSymbolHandler(key interface{}) (interface{}, error) {

	startTime := time.Now()
	log.Debug("Begin to load symbol file. uuid:", key)
	a, err := Symbol.LoadSymbolFile(cfg.Path, key.(string), &Symbol.ConfigInfo{BlockCount: cfg.BlockCount, LineCount: cfg.LineCount})
	endTime := time.Now()

	if err != nil {
		log.Error("Symbol file is be load! err:", err.Error(), " Use time :", (endTime.UnixNano() - startTime.UnixNano()), "ns")
	} else {
		log.Info("Load symbol file. uuid:", key, " Use time :", (endTime.UnixNano() - startTime.UnixNano()), "ns")
	}
	return a, err
}

// memory cache handle  ---- addedFunc to cache success
func AddedSymbolHandler(key, value interface{}) {
	log.Info("Load symbol to cache success. uuid:", key.(string))
}

// memory cache handle  ---- EvictedFunc from cache
func EvictedSymbolHandler(key, value interface{}) {
	log.Info("symbol is evicted from cache. uuid:", key.(string))
}
