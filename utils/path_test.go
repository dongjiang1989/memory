package Utils

import (
	"log"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_util_realpath(t *testing.T) {

	assert := assert.New(t)

	abc := RealPath()
	log.Println(abc)

	assert.NotEqual(abc, ".", "is not equal")
}

func Test_util_abspath(t *testing.T) {

	assert := assert.New(t)

	abc := AbsPath()
	log.Println(abc)

	assert.NotEqual(abc, ".", "is not equal")
}

func Test_util_NormPath(t *testing.T) {

	assert := assert.New(t)

	abc := NormPath()
	log.Println(abc)

	assert.NotEqual(abc, ".", "is not equal")
}

func Test_util_Dirname(t *testing.T) {

	assert := assert.New(t)

	abc := Dirname()
	log.Println(abc)

	assert.NotEqual(abc, ".", "is not equal")
}

func Test_util_Filename(t *testing.T) {

	assert := assert.New(t)

	abc := Filename()
	log.Println(abc)

	assert.NotEqual(abc, ".", "is not equal")
}
