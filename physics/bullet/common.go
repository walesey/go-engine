package bullet

import (
	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

func getVector(getter func(dest *[3]float32)) vmath.Vector3 {
	dest := [3]float32{0, 0, 0}
	getter(&dest)
	return vmath.Vector3{float64(dest[0]), float64(dest[1]), float64(dest[2])}
}

func getQuaternion(getter func(dest *[4]float32)) vmath.Quaternion {
	dest := [4]float32{0, 0, 0, 0}
	getter(&dest)
	return vmath.Quaternion{float64(dest[1]), float64(dest[2]), float64(dest[3]), float64(dest[0])}
}

func setVector(setter func(dest *[3]float32), value vmath.Vector3) {
	setter(&[3]float32{
		float32(value.X),
		float32(value.Y),
		float32(value.Z),
	})
}

func setQuaternion(setter func(dest *[4]float32), value vmath.Quaternion) {
	setter(&[4]float32{
		float32(value.W),
		float32(value.X),
		float32(value.Y),
		float32(value.Z),
	})
}

func TriangleMeshShapeFromGeometry(geometry *renderer.Geometry) gobullet.TriangleMeshShape {
	mesh := gobullet.NewTriangleMesh(false, false)
	for i := 0; i < len(geometry.Indicies); i = i + 3 {
		index := geometry.Indicies[i]
		v1 := [3]float32{geometry.Verticies[index*18], geometry.Verticies[index*18+1], geometry.Verticies[index*18+2]}
		index = geometry.Indicies[i+1]
		v2 := [3]float32{geometry.Verticies[index*18], geometry.Verticies[index*18+1], geometry.Verticies[index*18+2]}
		index = geometry.Indicies[i+2]
		v3 := [3]float32{geometry.Verticies[index*18], geometry.Verticies[index*18+1], geometry.Verticies[index*18+2]}

		mesh.AddTriangle(&v1, &v2, &v3, true)
	}
	return mesh.NewTriangleMeshShape(true, true)
}

func CollisionShapeFromGeometry(geometry *renderer.Geometry, cullThreshold float64) gobullet.CollisionShape {
	verts := assets.PointsFromGeometry(geometry, cullThreshold)
	convexHull := gobullet.NewConvexHullShape()
	for _, vert := range *verts {
		convexHull.AddVertex(float32(vert.X), float32(vert.Y), float32(vert.Z))
	}
	return gobullet.CollisionShape(convexHull)
}
