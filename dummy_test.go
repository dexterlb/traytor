package traytor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddOne(t *testing.T) {
	require := require.New(t)

	require.Equal(AddOne(5), 6)
}
