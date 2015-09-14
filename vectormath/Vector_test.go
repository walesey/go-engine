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

func TestNormalize(t *testing.T) {
    v := Vector3{4,2,-1}
    v = v.Normalize()
    assert.EqualValues(t,  0.87287156094396952506438994166248, v.X, "NormalizeX")
    assert.EqualValues(t,  0.43643578047198476253219497083124, v.Y, "NormalizeY")
    assert.EqualValues(t, -0.21821789023599238126609748541562, v.Z, "NormalizeZ")
}

func TestAdd(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v = v.Add(u);
    assert.EqualValues(t, 6, v.X, "AddX")
    assert.EqualValues(t, 10, v.Y, "AddY")
    assert.EqualValues(t, 15, v.Z, "AddZ")
}

func TestSubtract(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v = v.Subtract(u);
    assert.EqualValues(t, 0, v.X, "SubtractX")
    assert.EqualValues(t, 6, v.Y, "SubtractY")
    assert.EqualValues(t, 5, v.Z, "SubtractZ")
}

func TestMultiply(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v = v.Multiply(u);
    assert.EqualValues(t, 9, v.X, "MultiplyX")
    assert.EqualValues(t, 16, v.Y, "MultiplyY")
    assert.EqualValues(t, 50, v.Z, "MultiplyZ")
}

func TestMultiplyScalar(t *testing.T) {
    v := Vector3{10,15,20}
    s := 5.0
    v = v.MultiplyScalar(s);
    assert.EqualValues(t, 50, v.X, "MultiplyScalarX")
    assert.EqualValues(t, 75, v.Y, "MultiplyScalarY")
    assert.EqualValues(t, 100, v.Z, "MultiplyScalarZ")
}

func TestDivide(t *testing.T) {
    v := Vector3{3,8,10}
    u := Vector3{3,2,5}
    v = v.Divide(u);
    assert.EqualValues(t, 1, v.X, "DivideX")
    assert.EqualValues(t, 4, v.Y, "DivideY")
    assert.EqualValues(t, 2, v.Z, "DivideZ")
}

func TestDivideScalar(t *testing.T) {
    v := Vector3{10,15,20}
    s := 5.0
    v = v.DivideScalar(s);
    assert.EqualValues(t, 2, v.X, "DivideScalarX")
    assert.EqualValues(t, 3, v.Y, "DivideScalarY")
    assert.EqualValues(t, 4, v.Z, "DivideScalarZ")
}

func TestCross(t *testing.T) {
    v := Vector3{0,1,0}
    u := Vector3{1,0,0}
    v = v.Cross(u);
    assert.EqualValues(t, 0, v.X, "CrossX")
    assert.EqualValues(t, 0, v.Y, "CrossY")
    assert.EqualValues(t, 1, v.Z, "CrossZ")
}

func TestLerp(t *testing.T) {
    v := Vector3{10,1,0}.Lerp(Vector3{-20,8,10}, 0.0);
    assert.EqualValues(t, Vector3{10,1,0}, v, "Lerp amount 0")

    v = Vector3{10,1,0}.Lerp(Vector3{-20,8,10}, 0.5);
    assert.EqualValues(t, Vector3{-5,4.5,5}, v, "Lerp amount 0.5")

    v = Vector3{10,1,0}.Lerp(Vector3{-20,8,10}, 1.0);
    assert.EqualValues(t, Vector3{-20,8,10}, v, "Lerp amount 1")
}

func TestDot(t *testing.T) {
    v := Vector3{2,2,2}
    u := Vector3{0,1,0}
    dot := v.Dot(u);
    assert.EqualValues(t, 2, dot, "Dot")
}
