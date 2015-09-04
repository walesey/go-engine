package vectorMath

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLengthSquared(t *testing.T) {
    v := Vector3{2,2,2}
    lenSq := v.LengthSquared()
    assert.EqualValues(t, 12, lenSq, "LengthSquared")
}

func TestLength(t *testing.T) {
    v := Vector3{2,2,2}
    lenSq := RoundHalfUp(v.Length() * v.Length())
    assert.EqualValues(t, 12, lenSq, "Length")
}


func TestCopy(t *testing.T) {
    v := Vector3{3,8,10}
    u := v.Copy();
    assert.EqualValues(t, v.X, u.X, "CopyX")
    assert.EqualValues(t, v.Y, u.Y, "CopyY")
    assert.EqualValues(t, v.Z, u.Z, "CopyZ")
}

func TestAdd(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v.Add(u);
    assert.EqualValues(t, 6, v.X, "AddX")
    assert.EqualValues(t, 10, v.Y, "AddY")
    assert.EqualValues(t, 15, v.Z, "AddZ")
}

func TestSubtract(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v.Subtract(u);
    assert.EqualValues(t, 0, v.X, "SubtractX")
    assert.EqualValues(t, 6, v.Y, "SubtractY")
    assert.EqualValues(t, 5, v.Z, "SubtractZ")
}

func TestMultiply(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v.Multiply(u);
    assert.EqualValues(t, 9, v.X, "MultiplyX")
    assert.EqualValues(t, 16, v.Y, "MultiplyY")
    assert.EqualValues(t, 50, v.Z, "MultiplyZ")
}

func TestMultiplyScalar(t *testing.T) {
    v := Vector3{10,15,20}
    s := 5.0
    v.MultiplyScalar(s);
    assert.EqualValues(t, 50, v.X, "MultiplyScalarX")
    assert.EqualValues(t, 75, v.Y, "MultiplyScalarY")
    assert.EqualValues(t, 100, v.Z, "MultiplyScalarZ")
}

func TestDivide(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v.Divide(u);
    assert.EqualValues(t, 1, v.X, "DivideX")
    assert.EqualValues(t, 4, v.Y, "DivideY")
    assert.EqualValues(t, 2, v.Z, "DivideZ")
}

func TestDivideScalar(t *testing.T) {
    v := Vector3{10,15,20}
    s := 5.0
    v.DivideScalar(s);
    assert.EqualValues(t, 2, v.X, "DivideScalarX")
    assert.EqualValues(t, 3, v.Y, "DivideScalarY")
    assert.EqualValues(t, 4, v.Z, "DivideScalarZ")
}

func TestCross(t *testing.T) {
    v := &Vector3{0,1,0}
    u := Vector3{1,0,0}
    v = v.Cross(u);
    assert.EqualValues(t, 0, v.X, "CrossX")
    assert.EqualValues(t, 0, v.Y, "CrossY")
    assert.EqualValues(t, 1, v.Z, "CrossZ")
}

func TestDot(t *testing.T) {
    v := Vector3{2,2,2}
    u := Vector3{0,1,0}
    dot := v.Dot(u);
    assert.EqualValues(t, 2, dot, "Dot")
}
