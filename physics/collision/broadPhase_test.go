package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
)

func TestBoundingBox(t *testing.T) {
	bb1 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb2 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	bb2.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	assert.True(t, bb1.Overlap(bb2), "Bounding box overlap")
	bb1.Offset(vmath.Vector3{0, 0, 1.1}, vmath.IdentityQuaternion())
	assert.False(t, bb1.Overlap(bb2), "Bounding box not overlap")
}
