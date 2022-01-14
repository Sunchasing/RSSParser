package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExample(t *testing.T){
	assert.Equal(t, true, true)
	require.True(t, true)

}