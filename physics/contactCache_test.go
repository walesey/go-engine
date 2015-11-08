package physics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	contactCache := NewContactCache()
	contactCache.Add(1, 9)
	assert.True(t, contactCache.Contains(1, 9))
	assert.False(t, contactCache.Contains(1, 8))
	assert.False(t, contactCache.Contains(2, 9))
}

func TestRemove(t *testing.T) {
	contactCache := NewContactCache()
	contactCache.Add(1, 9)
	assert.True(t, contactCache.Contains(1, 9))
	contactCache.Remove(1, 9)
	assert.False(t, contactCache.Contains(1, 9))
}

func TestReverse(t *testing.T) {
	contactCache := NewContactCache()
	contactCache.Add(1, 9)
	assert.True(t, contactCache.Contains(1, 9))
	assert.True(t, contactCache.Contains(9, 1))
	contactCache.Remove(1, 9)
	assert.False(t, contactCache.Contains(1, 9))
	assert.False(t, contactCache.Contains(9, 1))
}

func TestClearingOldContacts(t *testing.T) {
	contactCache := NewContactCache()
	contactCache.Add(1, 9)
	contactCache.Add(2, 9)
	contactCache.Add(3, 9)
	contactCache.MarkContactsAsOld()
	contactCache.Add(4, 9)
	contactCache.Add(2, 9)
	contactCache.CleanOldContacts()
	assert.False(t, contactCache.Contains(1, 9))
	assert.True(t, contactCache.Contains(2, 9))
	assert.False(t, contactCache.Contains(3, 9))
	assert.True(t, contactCache.Contains(4, 9))
}
