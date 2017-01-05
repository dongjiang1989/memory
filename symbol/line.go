package Symbol

import (
	"bytes"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

type SymbolLine struct {
	begin  int64
	end    int64
	size   uint64
	detail string
}

func (s *SymbolLine) Begin() int64 {
	return s.begin
}

func (s *SymbolLine) Detail() string {
	return s.detail
}

func (s *SymbolLine) End() int64 {
	return s.end
}

func (s *SymbolLine) Size() uint64 {
	return s.size
}

func (s *SymbolLine) IsEmpty() bool {
	return s == nil || (s.begin == s.end && s.begin == 0) || s.size == 0
}

// demo : e94 ee0 +[VMUArchitecture architectureWithCpuType:cpuSubtype:]
func LoadSymbolLine(line string) *SymbolLine {
	if line == "" {
		return nil
	} else {
		info := strings.SplitN(line, "\t", 3)
		if len(info) != 3 {
			log.Error("one line to Split is not Symbol line Temp! str:", line, "array:", info)
			return nil
		} else {
			begin, end, infoStr := info[0], info[1], info[2]
			size := uint64(len([]byte(line)))

			bg, err := strconv.ParseInt(begin, 16, 64)
			if err != nil {
				log.Error("16 to 10 is err: ", err.Error())
				return nil
			}

			ed, err := strconv.ParseInt(end, 16, 64)
			if err != nil {
				log.Error("16 to 10 is err: ", err.Error())
				return nil
			}

			return &SymbolLine{begin: bg, end: ed, size: size, detail: infoStr}
		}

	}
}

func LoadSymbolLineByte(line []byte) *SymbolLine {
	if bytes.Equal(line, []byte{}) {
		return nil
	} else {
		info := bytes.SplitN(line, []byte{'\t'}, 3)
		if len(info) != 3 {
			log.Error("one line to Split is not Symbol line Temp! str:", line, "array:", info)
			return nil
		} else {
			begin, end, infoStr := info[0], info[1], info[2]
			size := uint64(len(line))

			bg, err := strconv.ParseInt(string(begin), 16, 64)
			if err != nil {
				log.Error("16 to 10 is err: ", err.Error())
				return nil
			}

			ed, err := strconv.ParseInt(string(end), 16, 64)
			if err != nil {
				log.Error("16 to 10 is err: ", err.Error())
				return nil
			}

			return &SymbolLine{begin: bg, end: ed, size: size, detail: string(infoStr)}
		}

	}
}
