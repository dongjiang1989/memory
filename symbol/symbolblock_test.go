package Symbol

import (
	"strconv"
	"strings"
	"testing"

	log "github.com/cihub/seelog"

	"github.com/stretchr/testify/assert"
)

func Test_SymbolBlock(t *testing.T) {
	assert := assert.New(t)

	sb := NewSymbolBlock(100)
	//	head := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
	//Format: Mach-O/32-Bit
	//Arch:   armv7s
	//Symbols:    1570
	//Tool Version:   1.0
	//File Version:   undefine
	//Built Time: 2016-08-05 11:51:32
	//UUID:   fffeb18603dd3abda51d0066bedb3c2c`
	//	sb.Header = LoadSymbolHeader(head)
	//	log.Debug(sb.Header)

	str1 := `
cd0	e04	+[VMUArchitecture initialize]
e04	e14	+[VMUArchitecture currentArchitecture]
e14	e24	+[VMUArchitecture anyArchitecture]
e24	e34	+[VMUArchitecture ppcArchitecture]
e34	e44	+[VMUArchitecture ppc32Architecture]
e44	e54	+[VMUArchitecture ppc64Architecture]
`

	for _, line := range strings.Split(str1, "\n") {
		aa := LoadSymbolLine(line)
		log.Debug(aa, str1)
		sb.Append(aa)
	}

	index, _ := strconv.ParseInt("e19", 16, 64)

	ab, err := sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture anyArchitecture]")

	index, _ = strconv.ParseInt("cd0", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture initialize]")

	index, _ = strconv.ParseInt("e04", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture currentArchitecture]")

	index, _ = strconv.ParseInt("ad0", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("e44", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture ppc64Architecture]")

	index, _ = strconv.ParseInt("e54", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("e94", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("e54", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("e52", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture ppc64Architecture]")

}

func Test_SymbolBlockOneLine11(t *testing.T) {
	assert := assert.New(t)

	sb := NewSymbolBlock(100)
	//	head := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
	//Format: Mach-O/32-Bit
	//Arch:   armv7s
	//Symbols:    1570
	//Tool Version:   1.0
	//File Version:   undefine
	//Built Time: 2016-08-05 11:51:32
	//UUID:   fffeb18603dd3abda51d0066bedb3c2c`
	//	sb.Header = LoadSymbolHeader(head)
	//	log.Debug(sb.Header)

	str1 := `cd0	e04	+[VMUArchitecture initialize]`

	for _, line := range strings.Split(str1, "\n") {
		aa := LoadSymbolLine(line)
		log.Debug(aa, str1)
		sb.Append(aa)
	}

	index, _ := strconv.ParseInt("a00", 16, 64)
	ab, err := sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("cd0", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture initialize]")

	index, _ = strconv.ParseInt("cd4", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture initialize]")

	index, _ = strconv.ParseInt("e04", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("e94", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)
}

func Test_SymbolBlockOneLine12(t *testing.T) {
	assert := assert.New(t)

	sb := NewSymbolBlock(100)
	//	head := `File:   ./SystemSymbol/8.1.1 (12B436)/Symbols/System/Library/PrivateFrameworks/Symbolication.framework/Symbolication
	//Format: Mach-O/32-Bit
	//Arch:   armv7s
	//Symbols:    1570
	//Tool Version:   1.0
	//File Version:   undefine
	//Built Time: 2016-08-05 11:51:32
	//UUID:   fffeb18603dd3abda51d0066bedb3c2c`
	//	sb.Header = LoadSymbolHeader(head)
	//	log.Debug(sb.Header)

	str1 := `cd0	cd0	+[VMUArchitecture initialize]`

	for _, line := range strings.Split(str1, "\n") {
		aa := LoadSymbolLine(line)
		log.Debug(aa, str1)
		sb.Append(aa)
	}

	index, _ := strconv.ParseInt("a00", 16, 64)
	ab, err := sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("cd0", 16, 64)
	ab, err = sb.BinSearch(index)
	log.Debug(index, sb.begin, sb.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture initialize]")

}
