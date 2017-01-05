package Symbol

import (
	"errors"
	"sync"

	log "github.com/cihub/seelog"
)

func NewSymbol(maxcount int) *Symbol {
	c := &Symbol{}
	c.initSymbol(maxcount)
	return c
}

type Symbol struct {
	mu       sync.RWMutex
	header   *Header
	begin    int64
	end      int64
	maxCount int
	size     uint64
	symbols  []*SymbolBlock
}

func (s *Symbol) initSymbol(maxcount int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.begin = int64(0)
	s.end = int64(0)
	s.symbols = make([]*SymbolBlock, 0)
	s.maxCount = maxcount
	s.header = &Header{}
	s.size = uint64(0)
}

func (s *Symbol) SetHeader(h *Header) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.header = h
}

func (s *Symbol) Header() *Header {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.header
}

func (s *Symbol) MaxCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxCount
}

func (s *Symbol) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.symbols)
}

func (s *Symbol) Size() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.size
}

func (s *Symbol) Has(index int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s == nil {
		return false
	}

	if s.Len() <= 0 {
		return false
	}

	if s.begin > s.end || (s.begin == s.end && s.begin == 0) {
		log.Warn("Symbol is empty!")
		return false
	}

	// 左开右闭
	if (index == s.begin && s.begin != 0 && s.begin <= s.end) || (index > s.begin && index < s.end) {
		return true
	}

	return false
}

func (s *Symbol) Append(sl *SymbolBlock) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sl == nil || sl.IsEmpty() {
		return false
	}

	if s.begin == 0 || (s.begin > sl.Begin() && sl.Begin() > 0) {
		s.begin = sl.Begin()
	}

	if s.end == 0 || s.end < sl.End() {
		s.end = sl.End()
	}

	s.symbols = append(s.symbols, sl)
	s.size = s.size + sl.Size()
	return true

}

// bin search
func (s *Symbol) BinSearch(index int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Has(index) && s.Len() >= 1 {
		lo, hi := 0, s.Len()-1
		for lo <= hi {
			m := (lo + hi) >> 1
			if s.symbols[m].Begin() < index {
				lo = m + 1
			} else if s.symbols[m].Begin() > index {
				hi = m - 1
			} else {
				return s.symbols[m].BinSearch(index)
			}
		}

		if lo == s.Len() {
			if index == s.symbols[lo-1].Begin() || index < s.symbols[lo-1].End() {
				return s.symbols[lo-1].BinSearch(index)
			} else {
				return "", errors.New("index is not in Symbol!")
			}
		} else {
			return s.symbols[lo-1].BinSearch(index)
		}
	}
	return "", errors.New("index is not in Symbol!")
}
