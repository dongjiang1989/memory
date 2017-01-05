package Symbol

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_symbolline(t *testing.T) {
	assert := assert.New(t)

	str := "e94	ee0	+[VMUArchitecture architectureWithCpuType:cpuSubtype:]"

	obj := LoadSymbolLine(str)

	assert.NotNil(obj)

	begin, _ := strconv.ParseInt("e94", 16, 64)
	assert.Equal(obj.begin, begin)
	end, _ := strconv.ParseInt("ee0", 16, 64)
	assert.Equal(obj.end, end)
	assert.Equal(obj.detail, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]")
}

func Test_symbollineforempt(t *testing.T) {
	assert := assert.New(t)

	str := "e94 ee0 +[VMUArchitecture architectureWithCpuType:cpuSubtype:]"

	obj := LoadSymbolLine(str)

	assert.Nil(obj)

}

func Test_symbollineforline4(t *testing.T) {
	assert := assert.New(t)

	str := "e94	ee0	+[VMUArchitecture architectureWithCpuType:cpuSubtype:]	sss"

	obj := LoadSymbolLine(str)

	assert.NotNil(obj)
	assert.Equal(obj.detail, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]	sss")

}

func Test_symbollineforline2(t *testing.T) {
	assert := assert.New(t)

	str := "e94	ee0"

	obj := LoadSymbolLine(str)

	assert.Nil(obj)
}

func Test_symbollineforempt222(t *testing.T) {
	assert := assert.New(t)

	str := "e94	e94	+[VMUArchitecture architectureWithCpuType:cpuSubtype:]"

	obj := LoadSymbolLine(str)

	assert.NotNil(obj)

}

func Test_symbollineb(t *testing.T) {
	assert := assert.New(t)

	str := "e94	ee0	+[VMUArchitecture architectureWithCpuType:cpuSubtype:]"

	obj := LoadSymbolLineByte([]byte(str))

	assert.NotNil(obj)

	begin, _ := strconv.ParseInt("e94", 16, 64)
	assert.Equal(obj.begin, begin)
	end, _ := strconv.ParseInt("ee0", 16, 64)
	assert.Equal(obj.end, end)
	assert.Equal(obj.detail, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]")
}

func Test_symbollineforemptb(t *testing.T) {
	assert := assert.New(t)

	str := "e94 ee0 +[VMUArchitecture architectureWithCpuType:cpuSubtype:]"

	obj := LoadSymbolLineByte([]byte(str))

	assert.Nil(obj)

}

func Test_symbollineforline4b(t *testing.T) {
	assert := assert.New(t)

	str := "e94	ee0	+[VMUArchitecture architectureWithCpuType:cpuSubtype:]	sss"

	obj := LoadSymbolLineByte([]byte(str))

	assert.NotNil(obj)
	assert.Equal(obj.detail, "+[VMUArchitecture architectureWithCpuType:cpuSubtype:]	sss")

}

func Test_symbollineforline2b(t *testing.T) {
	assert := assert.New(t)

	str := "e94	ee0"

	obj := LoadSymbolLineByte([]byte(str))

	assert.Nil(obj)
}

func Test_symbollineforempt222b(t *testing.T) {
	assert := assert.New(t)

	str := "e94	e94	+[VMUArchitecture architectureWithCpuType:cpuSubtype:]"

	obj := LoadSymbolLineByte([]byte(str))

	assert.NotNil(obj)

}
