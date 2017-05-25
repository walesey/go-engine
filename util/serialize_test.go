package util

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestUint64(t *testing.T) {
	var i uint64 = 50
	data := SerializeArgs(i)
	result := UInt64frombytes(bytes.NewBuffer(data))
	assert.EqualValues(t, i, result)
}

func TestAbc(t *testing.T) {
	a := []string{"a", "b"}
	data := SerializeArgs(len(a))
	result := UInt32frombytes(bytes.NewBuffer(data))
	assert.EqualValues(t, 2, result)
}
