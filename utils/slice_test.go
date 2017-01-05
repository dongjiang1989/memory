package Utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_InSlice(t *testing.T) {
	assert := assert.New(t)

	assert.True(InSliceIface("aaa", []interface{}{"aaaa", "aaa"}), "解析文件正常，不符合需求！")

	assert.False(InSliceIface("b", []interface{}{"aaaa", "aaa"}), "解析文件正常，不符合需求！")

	assert.False(InSliceIface("aaa", []interface{}{}), "解析文件正常，不符合需求！")

}
