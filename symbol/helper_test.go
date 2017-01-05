package Symbol

import (
	"os"
	"strconv"
	"testing"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
)

func Test_symbol(t *testing.T) {
	assert := assert.New(t)

	log.Debug(os.Getwd())
	pwd, _ := os.Getwd()
	a, err := LoadSymbolFile(pwd+"/utdata/", "fffeb18603dd3abda51d0066bedb3c2c", &ConfigInfo{BlockCount: 1, LineCount: 9})

	log.Debug("========================================", a.Header().File)

	assert.Nil(err)
	assert.NotNil(a)
	assert.Equal(a.Header().Uuid, "fffeb18603dd3abda51d0066bedb3c2c")

	assert.Equal(a.Len(), 6)

}

func Test_symbol1(t *testing.T) {
	assert := assert.New(t)

	log.Debug(os.Getwd())
	pwd, _ := os.Getwd()
	a, err := LoadSymbolFile(pwd+"/utdata/", "fffeb18603dd3abda51d0066bedb3c2c", &ConfigInfo{BlockCount: 1000, LineCount: 9})

	log.Debug("========================================", a.Header().File)

	assert.Nil(err)
	assert.NotNil(a)
	assert.Equal(a.Header().Uuid, "fffeb18603dd3abda51d0066bedb3c2c")

	assert.Equal(a.Len(), 142)

}

func Test_symbol64(t *testing.T) {
	assert := assert.New(t)

	log.Debug(os.Getwd())
	pwd, _ := os.Getwd()
	a, err := LoadSymbolFile(pwd+"/utdata/", "fffeb18603dd3abda51d0066bedb3c2c", &ConfigInfo{BlockCount: 64, LineCount: 9})

	log.Debug("========================================", a.Header().File)

	assert.Nil(err)
	assert.NotNil(a)
	assert.Equal(a.Header().Uuid, "fffeb18603dd3abda51d0066bedb3c2c")

	log.Critical("--=-=-=-=-=-=-", a.Len())

	assert.Equal(a.Len(), 40)

}

func Test_symbol_128(t *testing.T) {
	assert := assert.New(t)

	log.Debug(os.Getwd())
	pwd, _ := os.Getwd()
	a, err := LoadSymbolFile(pwd+"/utdata/", "fffeb18603dd3abda51d0066bedb3c2c", &ConfigInfo{BlockCount: 254, LineCount: 9})

	log.Debug("========================================", a.Header().File)

	assert.Nil(err)
	assert.NotNil(a)
	assert.Equal(a.Header().Uuid, "fffeb18603dd3abda51d0066bedb3c2c")

	assert.Equal(a.Len(), 80)

	index, _ := strconv.ParseInt("a00", 16, 64)
	ab, err := a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.NotNil(err)

	index, _ = strconv.ParseInt("cd0", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture initialize]")

	index, _ = strconv.ParseInt("cd4", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture initialize]")

	index, _ = strconv.ParseInt("e04", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture currentArchitecture]")

	index, _ = strconv.ParseInt("e94", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]")

	index, _ = strconv.ParseInt("e95", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.Nil(err)
	assert.Equal(ab, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]")

	//TODO maybe bug
	index, _ = strconv.ParseInt("618ff90", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.NotNil(err)
	//assert.Equal(ab, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]")

	index, _ = strconv.ParseInt("618ff91", 16, 64)
	ab, err = a.BinSearch(index)
	log.Debug(index, a.begin, a.end)
	assert.NotNil(err)

}
