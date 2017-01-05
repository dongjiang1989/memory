package Symbol

import (
	"errors"
	"sync"

	log "github.com/cihub/seelog"
)

func NewSymbolBlock(count int) *SymbolBlock {
	c := &SymbolBlock{}
	c.initSymbolBlock(count)
	return c
}

type SymbolBlock struct {
	mu    sync.RWMutex
	begin int64
	end   int64
	count int
	size  uint64
	sb    []*SymbolLine
}

func (s *SymbolBlock) Begin() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.begin
}

func (s *SymbolBlock) End() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.end
}

func (s *SymbolBlock) Size() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.size
}

func (s *SymbolBlock) initSymbolBlock(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.begin = int64(0)
	s.end = int64(0)
	s.sb = make([]*SymbolLine, 0)
	s.count = count
	s.size = uint64(0)
}

func (s *SymbolBlock) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.sb)
}

func (s *SymbolBlock) IsEmpty() bool {
	return s == nil || (s.begin == s.end && s.begin == 0)
}

func (s *SymbolBlock) Has(index int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s == nil {
		return false
	}

	if s.Len() <= 0 {
		return false
	}

	if s.begin > s.end || (s.begin == s.end && s.begin == 0) {
		log.Warn("SymbolBlock is empty!")
		return false
	}

	// 左开右闭
	if (index == s.begin && s.begin != 0 && s.begin <= s.end) || (index > s.begin && index < s.end) {
		return true
	}

	return false
}

func (s *SymbolBlock) Append(sl *SymbolLine) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sl == nil || sl.IsEmpty() {
		return false
	}

	if len(s.sb) < s.count {

		if s.begin == 0 || (s.begin > sl.Begin() && sl.Begin() > 0) {
			s.begin = sl.Begin()
		}

		if s.end == 0 || s.end < sl.End() {
			s.end = sl.End()
		}

		s.sb = append(s.sb, sl)
		s.size = s.size + sl.Size()
		return true
	} else {
		return false
	}
}

// bin search
func (s *SymbolBlock) BinSearch(index int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Has(index) && s.Len() >= 1 {
		lo, hi := 0, s.Len()-1
		for lo <= hi {
			m := (lo + hi) >> 1
			if s.sb[m].Begin() < index {
				lo = m + 1
			} else if s.sb[m].Begin() > index {
				hi = m - 1
			} else {
				return s.sb[m].Detail(), nil
			}
		}

		if lo == s.Len() {
			if index == s.sb[lo-1].Begin() || index < s.sb[lo-1].End() {
				return s.sb[lo-1].Detail(), nil
			} else {
				return "", errors.New("index is not in Symbol block!")
			}
		} else {
			return s.sb[lo-1].Detail(), nil
		}
	}
	return "", errors.New("index is not in Symbol block!")
}
