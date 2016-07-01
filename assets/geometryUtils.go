package assets

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

// IntersectGeometry - returns a list of line segments resulting from the xy plane intersection of the geometry
func IntersectGeometry(geometry *renderer.Geometry) [][2]vmath.Vector2 {
	segments := make([][2]vmath.Vector2, 0)
	for i := 0; i < len(geometry.Indicies); i = i + 3 {
		index := geometry.Indicies[i]
		v1 := vmath.Vector3{float64(geometry.Verticies[index*renderer.VertexStride]), float64(geometry.Verticies[index*renderer.VertexStride+1]), float64(geometry.Verticies[index*renderer.VertexStride+2])}
		index = geometry.Indicies[i+1]
		v2 := vmath.Vector3{float64(geometry.Verticies[index*renderer.VertexStride]), float64(geometry.Verticies[index*renderer.VertexStride+1]), float64(geometry.Verticies[index*renderer.VertexStride+2])}
		index = geometry.Indicies[i+2]
		v3 := vmath.Vector3{float64(geometry.Verticies[index*renderer.VertexStride]), float64(geometry.Verticies[index*renderer.VertexStride+1]), float64(geometry.Verticies[index*renderer.VertexStride+2])}

		va, vb, vc := v1, v2, v3
		if (v1.Y < 0 && v2.Y > 0 && v3.Y > 0) || (v1.Y > 0 && v2.Y < 0 && v3.Y < 0) {
			va, vb, vc = v1, v2, v3
		} else if (v1.Y > 0 && v2.Y < 0 && v3.Y > 0) || (v1.Y < 0 && v2.Y > 0 && v3.Y < 0) {
			va, vb, vc = v2, v1, v3
		} else if (v1.Y > 0 && v2.Y > 0 && v3.Y < 0) || (v1.Y < 0 && v2.Y < 0 && v3.Y > 0) {
			va, vb, vc = v3, v2, v1
		} else {
			continue
		}

		t_ab := -va.Y / (va.Y - vb.Y)
		t_ac := -va.Y / (va.Y - vc.Y)
		segments = append(segments, [2]vmath.Vector2{
			vmath.Vector2{va.X + (va.X-vb.X)*t_ab, va.Z + (va.Z-vb.Z)*t_ab},
			vmath.Vector2{va.X + (va.X-vc.X)*t_ac, va.Z + (va.Z-vc.Z)*t_ac},
		})
	}
	return segments
}

// Converts a geometry directly into points (does threshold culling optimisation)
func PointsFromGeometry(geometry *renderer.Geometry, cullThreshold float64) *[]vmath.Vector3 {
	verticies := make([]vmath.Vector3, 0, len(geometry.Indicies))
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := vmath.Vector3{
			float64(geometry.Verticies[index*renderer.VertexStride]),
			float64(geometry.Verticies[index*renderer.VertexStride+1]),
			float64(geometry.Verticies[index*renderer.VertexStride+2]),
		}
		//do culling
		include := true
		for _, vert := range verticies {
			if vert.Subtract(v).LengthSquared() < cullThreshold*cullThreshold {
				include = false
				break
			}
		}
		if include {
			verticies = append(verticies, v)
		}
	}
	return &verticies
}
