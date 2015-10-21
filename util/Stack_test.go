package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	mstack := CreateStack()
	mstack.Push("test")
	assert.EqualValues(t, 1, mstack.Len(), "Len MatStack should be 1")
	assert.EqualValues(t, "test", mstack.Get(0), "MatStack should return value on pop")
	popVal := mstack.Pop()
	assert.EqualValues(t, "test", popVal, "MatStack should return value on pop")
	assert.EqualValues(t, 0, mstack.Len(), "Len MatStack should be 0")

	mstack.Push("test1")
	mstack.Push("test2")
	assert.EqualValues(t, 2, mstack.Len(), "Len MatStack should be 1")
	assert.EqualValues(t, "test1", mstack.Get(0), "MatStack should return value on get")
	assert.EqualValues(t, "test2", mstack.Get(1), "MatStack should return value on get")
}
