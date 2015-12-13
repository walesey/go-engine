package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
)

func TestSimplex(t *testing.T) {
	simplex := NewSimplex()
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{1, 1, 1}})
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{2, 2, 2}})
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{3, 3, 3}})
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{4, 4, 4}})
	assert.EqualValues(t, 4, simplex.Len(), "simplex should have 4 entries")

	simplex.Remove(2)
	assert.EqualValues(t, 3, simplex.Len(), "simplex should have 3 entries")
	assert.EqualValues(t, vmath.Vector3{2, 2, 2}, simplex.Get(1).mPoint, "incorrect value in simplex")
	assert.EqualValues(t, vmath.Vector3{4, 4, 4}, simplex.Get(2).mPoint, "incorrect value in simplex")

	simplex.Remove(2)
	assert.EqualValues(t, 2, simplex.Len(), "simplex should have 2 entries")
	assert.EqualValues(t, vmath.Vector3{2, 2, 2}, simplex.Get(1).mPoint, "incorrect value in simplex")
}

func TestSimplexFaces(t *testing.T) {
	simplex := NewSimplex()
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{-0.1, -0.1, -0.05}})
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{1, -0.1, -0.05}})
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{-0.1, 1, -0.05}})
	simplex.Add(SimplexPoint{mPoint: vmath.Vector3{-0.1, -0.1, 1}})

	simplex.AddFace(SimplexFace{p1: 0, p2: 1, p3: 2})
	simplex.AddFace(SimplexFace{p1: 0, p2: 2, p3: 3})
	simplex.AddFace(SimplexFace{p1: 0, p2: 3, p3: 1})
	simplex.AddFace(SimplexFace{p1: 1, p2: 2, p3: 3})

	face, d := simplex.ClosestFace()
	assert.True(t, vmath.ApproxEqual(d, 0.05, 0.001), "distance to closest face should be 0.05")
	assert.EqualValues(t, 0, face.p1, "closest face should contain points 0, 1 and 2")
	assert.EqualValues(t, 1, face.p2, "closest face should contain points 0, 1 and 2")
	assert.EqualValues(t, 2, face.p3, "closest face should contain points 0, 1 and 2")

	norm := simplex.FaceNormal(0).Normalize()
	assert.True(t, norm.ApproxEqual(vmath.Vector3{0, 0, -1}, 0.1), "normal should face away from the origin")
	norm = simplex.FaceNormal(1).Normalize()
	assert.True(t, norm.ApproxEqual(vmath.Vector3{-1, 0, 0}, 0.1), "normal should face away from the origin")
	norm = simplex.FaceNormal(2).Normalize()
	assert.True(t, norm.ApproxEqual(vmath.Vector3{0, -1, 0}, 0.1), "normal should face away from the origin")
	norm = simplex.FaceNormal(3).Normalize()
	assert.True(t, norm.ApproxEqual(vmath.Vector3{0.6, 0.6, 0.6}, 0.2), "normal should face away from the origin")

	simplex.AddPointToFace(SimplexPoint{mPoint: vmath.Vector3{1, 1, 1}}, 3)
	assert.EqualValues(t, 5, simplex.Len(), "simplex should have an extra point")
	assert.EqualValues(t, 6, len(simplex.faces), "simplex should have 6 faces now")
	for i := 0; i < len(simplex.faces); i = i + 1 {
		assert.EqualValues(t, i, simplex.faces[i].index, "simplex should have faces all with correct index fields")
	}
}
