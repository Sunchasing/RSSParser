package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T){
	assert.Equal(t, true, true)
	require.True(t, true)

}