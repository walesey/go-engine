package physics

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/walesey/go-engine/vectormath"
)


func TestTriangleOverlap(t *testing.T) {
    t1 := Triangle{ vectormath.Vector3{-0.504109,0.540383,1.441891}, vectormath.Vector3{1.004138,-0.547662,0.706127}, vectormath.Vector3{0.499970,0.007280,-1.148018} }
    t2 := Triangle{ vectormath.Vector3{0.615065,-0.037162,-1.208676}, vectormath.Vector3{-0.355948,0.614617,0.413767}, vectormath.Vector3{0.892249,-0.426572,1.579071} }
    assert.EqualValues(t, false, t1.Overlap(t2), "Triangle Overlap should be false")

    t1 = Triangle{ vectormath.Vector3{-1,0,1}, vectormath.Vector3{1,0,1}, vectormath.Vector3{1,0,-1} }
    t2 = Triangle{ vectormath.Vector3{-1,0.3,1}, vectormath.Vector3{1,0.3,1}, vectormath.Vector3{1,0.3,-1} }
    assert.EqualValues(t, false, t1.Overlap(t2), "Triangle Overlap should be false (Parallel triangles)")

    t1 = Triangle{ vectormath.Vector3{9,0,1}, vectormath.Vector3{10,0,1}, vectormath.Vector3{10,0,-1} }
    t2 = Triangle{ vectormath.Vector3{-1,0,1}, vectormath.Vector3{1,0,1}, vectormath.Vector3{1,0,-1} }
    assert.EqualValues(t, false, t1.Overlap(t2), "Triangle Overlap should be false (coplanar triangles)")
	
    t1 = Triangle{ vectormath.Vector3{-0.504109,0.540383,1.441891}, vectormath.Vector3{1.004138,-0.547662,0.706127}, vectormath.Vector3{0.499970,0.007280,-1.148018} }
    t2 = Triangle{ vectormath.Vector3{0.048924,-1.436378,0.544679}, vectormath.Vector3{-0.679617,0.418397,0.715096}, vectormath.Vector3{0.775831,1.099488,-0.475613} }
    assert.EqualValues(t, true, t1.Overlap(t2), "Triangle Overlap should be true")
	
    t1 = Triangle{ vectormath.Vector3{-0.504109,0.540383,1.441891}, vectormath.Vector3{1.004138,-0.547662,0.706127}, vectormath.Vector3{0.499970,0.007280,-1.148018} }
    t2 = Triangle{ vectormath.Vector3{0.081168,0.966970,1.599979}, vectormath.Vector3{0.320969,-0.553781,0.323337}, vectormath.Vector3{0.749229,0.741485,-1.139154} }
    assert.EqualValues(t, true, t1.Overlap(t2), "Triangle Overlap should be true")
}

