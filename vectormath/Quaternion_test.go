package vectormath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	q := AngleAxis(1.57, Vector3{0, 1, 0})
	v := q.Apply(Vector3{0, 0, 1})
	expectV := Vector3{1, 0, 0}
	fmt.Printf("Expect %v to be %v", v, expectV)
	assert.True(t, v.ApproxEqual(expectV, 0.1))
}
