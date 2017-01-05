package Symbol

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	log "github.com/cihub/seelog"
)

type ConfigInfo struct {
	BlockCount int
	LineCount  int
}

//============================================================================================
func LoadSymbolFile(path, uuid string, cfg *ConfigInfo) (*Symbol, error) {

	fmarkname := fmt.Sprintf("%s/%s.symbol.success", path, uuid)

	if _, err := os.Stat(fmarkname); err != nil {
		log.Debug("Can not os.Stat filepath:", fmarkname, " err:", err.Error())
		return nil, err
	}

	fname := fmt.Sprintf("%s/%s.symbol", path, uuid)
	p, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Error("Can not ioutil.ReadFile filepath:", fname, " err:", err.Error())
		return nil, err
	}

	maxCount := int(math.Sqrt(float64(len(p) / cfg.BlockCount)))

	if maxCount < 1 {
		maxCount = 1
	}

	symbol := NewSymbol(maxCount)
	sb := NewSymbolBlock(symbol.MaxCount())
	for index, line := range bytes.Split(p, []byte{'\n'}) {
		if index < cfg.LineCount {
			symbol.Header().Add(string(line))
			continue
		}

		if string(line) == "" {
			continue
		}

		sl := LoadSymbolLine(string(line))
		if sl == nil {
			log.Warn("symbol readline is err! uuid:", symbol.Header().Uuid, " line:", index, "string:", string(line))
		} else {
			if !sb.Append(sl) {
				symbol.Append(sb)
				sb = NewSymbolBlock(symbol.MaxCount())
			}
		}
	}

	symbol.Append(sb)

	if symbol.Len() <= 0 {
		log.Error("can not load symbol file! uuid:", symbol.Header().Uuid, "symbol len:", symbol.Len())
		return nil, errors.New("can not load symbol file! uuid:" + symbol.Header().Uuid)
	}

	log.Debug("symbol filename:", fname, " file byte size: ", len(p), " symbol obj size:", symbol.Size(), " symbol block count:", symbol.Len(), " one block line:", maxCount)

	return symbol, nil
}
