package util

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestUint64(t *testing.T) {
	var i uint64 = 50
	data, err := SerializeArgs(i)
	assert.NoError(t, err)
	result, err := UInt64frombytes(bytes.NewBuffer(data))
	assert.NoError(t, err)
	assert.EqualValues(t, i, result)
}

func TestAbc(t *testing.T) {
	a := []string{"a", "b"}
	data, err := SerializeArgs(len(a))
	assert.NoError(t, err)
	result, err := UInt32frombytes(bytes.NewBuffer(data))
	assert.NoError(t, err)
	assert.EqualValues(t, 2, result)
}
