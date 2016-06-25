package vectormath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func expectApproxEqual2D(t *testing.T, v1, v2 Vector2) {
	assert.True(t, v1.ApproxEqual(v2, 0.001), fmt.Sprintf("Expect %v to approx equal %v", v1, v2))
}

func expectApproxEqual(t *testing.T, v1, v2 Vector3) {
	assert.True(t, v1.ApproxEqual(v2, 0.001), fmt.Sprintf("Expect %v to approx equal %v", v1, v2))
}

func TestRotate(t *testing.T) {
	v := Vector2{1, 1}
	rv := v.Rotate(1.57)
	expectApproxEqual2D(t, rv, Vector2{-1, 1})
}

func TestLengthSquared(t *testing.T) {
	v := Vector3{2, 2, 2}
	lenSq := v.LengthSquared()
	assert.EqualValues(t, 12, lenSq, "LengthSquared")
}

func TestLength(t *testing.T) {
	v := Vector3{2, 2, 2}
	lenSq := RoundHalfUp(v.Length() * v.Length())
	assert.EqualValues(t, 12, lenSq, "Length")
}

func TestNormalize(t *testing.T) {
	v := Vector3{4, 2, -1}
	v = v.Normalize()
	expectApproxEqual(t, v, Vector3{0.872, 0.436, -0.218})
}

func TestAdd(t *testing.T) {
	v := Vector3{3, 8, 10}
	u := Vector3{3, 2, 5}
	v = v.Add(u)
	assert.EqualValues(t, 6, v.X, "AddX")
	assert.EqualValues(t, 10, v.Y, "AddY")
	assert.EqualValues(t, 15, v.Z, "AddZ")
}

func TestSubtract(t *testing.T) {
	v := Vector3{3, 8, 10}
	u := Vector3{3, 2, 5}
	v = v.Subtract(u)
	assert.EqualValues(t, 0, v.X, "SubtractX")
	assert.EqualValues(t, 6, v.Y, "SubtractY")
	assert.EqualValues(t, 5, v.Z, "SubtractZ")
}

func TestMultiply(t *testing.T) {
	v := Vector3{3, 8, 10}
	u := Vector3{3, 2, 5}
	v = v.Multiply(u)
	assert.EqualValues(t, 9, v.X, "MultiplyX")
	assert.EqualValues(t, 16, v.Y, "MultiplyY")
	assert.EqualValues(t, 50, v.Z, "MultiplyZ")
}

func TestMultiplyScalar(t *testing.T) {
	v := Vector3{10, 15, 20}
	s := 5.0
	v = v.MultiplyScalar(s)
	assert.EqualValues(t, 50, v.X, "MultiplyScalarX")
	assert.EqualValues(t, 75, v.Y, "MultiplyScalarY")
	assert.EqualValues(t, 100, v.Z, "MultiplyScalarZ")
}

func TestDivide(t *testing.T) {
	v := Vector3{3, 8, 10}
	u := Vector3{3, 2, 5}
	v = v.Divide(u)
	assert.EqualValues(t, 1, v.X, "DivideX")
	assert.EqualValues(t, 4, v.Y, "DivideY")
	assert.EqualValues(t, 2, v.Z, "DivideZ")
}

func TestDivideScalar(t *testing.T) {
	v := Vector3{10, 15, 20}
	s := 5.0
	v = v.DivideScalar(s)
	assert.EqualValues(t, 2, v.X, "DivideScalarX")
	assert.EqualValues(t, 3, v.Y, "DivideScalarY")
	assert.EqualValues(t, 4, v.Z, "DivideScalarZ")
}

func TestCross(t *testing.T) {
	v := Vector3{0, 1, 0}
	u := Vector3{1, 0, 0}
	v = v.Cross(u)
	assert.EqualValues(t, 0, v.X, "CrossX")
	assert.EqualValues(t, 0, v.Y, "CrossY")
	assert.EqualValues(t, 1, v.Z, "CrossZ")
}

func TestLerp(t *testing.T) {
	v := Vector3{10, 1, 0}.Lerp(Vector3{-20, 8, 10}, 0.0)
	assert.EqualValues(t, Vector3{10, 1, 0}, v, "Lerp amount 0")

	v = Vector3{10, 1, 0}.Lerp(Vector3{-20, 8, 10}, 0.5)
	assert.EqualValues(t, Vector3{-5, 4.5, 5}, v, "Lerp amount 0.5")

	v = Vector3{10, 1, 0}.Lerp(Vector3{-20, 8, 10}, 1.0)
	assert.EqualValues(t, Vector3{-20, 8, 10}, v, "Lerp amount 1")
}

func TestDot(t *testing.T) {
	v := Vector3{2, 2, 2}
	u := Vector3{0, 1, 0}
	dot := v.Dot(u)
	assert.EqualValues(t, 2, dot, "Dot")
}

func TestApproxEqualVector(t *testing.T) {
	assert.True(t, Vector3{2.7654, -2.7654, 0}.ApproxEqual(Vector3{2.7154, -2.7154, 0.09}, 0.1), "ApproxEqual")
	assert.False(t, Vector3{2.7654, -2.7654, 0}.ApproxEqual(Vector3{2.7154, -2.7154, 0.09}, 0.01), "ApproxEqual")
}
