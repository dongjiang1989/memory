package Memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	assert := assert.New(t)

	var cases = []struct {
		hit  int
		miss int
		rate float64
	}{
		{3, 1, 0.75},
		{0, 1, 0.0},
		{3, 0, 1.0},
		{0, 0, 0.0},
	}

	for _, cs := range cases {
		st := &stats{}
		for i := 0; i < cs.hit; i++ {
			st.IncrHitCount()
		}
		for i := 0; i < cs.miss; i++ {
			st.IncrMissCount()
		}
		assert.Equal(cs.rate, st.HitRate(), "not equal")
	}

}

func TestStatsAllHitCount(t *testing.T) {
	assert := assert.New(t)

	var cases = []struct {
		hit  int
		miss int
		rate float64
	}{
		{3, 1, 0.75},
	}
	for _, cs := range cases {
		st := &stats{}
		for i := 0; i < cs.hit; i++ {
			st.IncrHitCount()
		}
		for i := 0; i < cs.miss; i++ {
			st.IncrMissCount()
		}
		assert.Equal(cs.rate, st.HitRate(), "not equal")
		assert.Equal(st.LookupCount(), uint64(cs.hit+cs.miss), "is not equal")
	}
}
