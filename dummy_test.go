package traytor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddOne(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(6, AddOne(5))
	assert.Equal(0, AddOne(-1), "must work for negative numbers as well")
}
