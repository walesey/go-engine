package renderer

import (
    "testing"
    "github.com/stretchr/testify/assert"
	"github.com/go-gl/mathgl/mgl32"
)

func TestMatStack(t *testing.T) {
	mstack := CreateStack()
	mstack.Push(mgl32.Ident4())
    assert.EqualValues(t, 1, mstack.Len(), "Len MatStack should be 1")
	popVal := mstack.Pop()
    assert.EqualValues(t, mgl32.Ident4(), popVal, "MatStack should return value on pop")
    assert.EqualValues(t, 0, mstack.Len(), "Len MatStack should be 0")
}

func TestMatStackMultiplyAll(t *testing.T) {
	mstack := CreateStack()
    mstack.Push(&GlTransform{mgl32.Ident4()})
    mstack.Push(&GlTransform{mgl32.Ident4()})
	mulVal := mstack.MultiplyAll()
    assert.EqualValues(t, mgl32.Ident4(), mulVal, "")

	v0 := mgl32.Vec4{1, 2, 3, 4}
	v1 := mgl32.Vec4{5, 6, 7, 8}
	v2 := mgl32.Vec4{9, 10, 11, 12}
	v3 := mgl32.Vec4{13, 14, 15, 16}
	m1 := mgl32.Mat4FromCols(v0, v1, v2, v3)

	v4 := mgl32.Vec4{14, 3, 8, 4}
	v5 := mgl32.Vec4{9, 1, 7, 8}
	v6 := mgl32.Vec4{9, 10, 11, 5}
	v7 := mgl32.Vec4{4, 7, 15, 2}
	m2 := mgl32.Mat4FromCols(v4, v5, v6, v7)

    mstack.Push(&GlTransform{m1})
    mstack.Push(&GlTransform{m2})
	mulVal12 := mstack.MultiplyAll()
	expected := mgl32.Ident4().Mul4(m1).Mul4(m2)

    assert.EqualValues(t, expected, mulVal12, "")
}