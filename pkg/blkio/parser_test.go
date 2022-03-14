package blkio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	dur, err := ParseDuration("2s")
	assert.Equal(t, nil, err)
	assert.Equal(t, "2s", dur.String())

	dur, err = ParseDuration("s")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1s", dur.String())
}
