package vectormath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	mat := Matrix3{
		1, 0, 0,
		0, -1, 0,
		0, 0, 1,
	}
	vec := Vector3{1, 1, 1}
	assert.EqualValues(t, Vector3{1, -1, 1}, mat.Transform(vec))
}
