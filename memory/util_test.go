package Memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinInt(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, minInt(1, 2), "not 1")
	assert.Equal(1, minInt(2, 1), "not 1")
	assert.Equal(1, minInt(1, 1), "not 1")
}

func TestMaxInt(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(2, maxInt(1, 2), "not 1")
	assert.Equal(2, maxInt(2, 1), "not 1")
	assert.Equal(1, maxInt(1, 1), "not 1")
}
