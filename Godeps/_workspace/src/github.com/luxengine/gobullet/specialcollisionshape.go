package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

func NewIcosahedronShape(r float32) CollisionShape {
	icosahedron := NewConvexHullShape()
	icosahedron.AddVertex(r*-0.26286500, r*0.0000000, r*0.42532500)
	icosahedron.AddVertex(r*-0.26286500, r*0.0000000, r*0.42532500)
	icosahedron.AddVertex(r*0.26286500, r*0.0000000, r*0.42532500)
	icosahedron.AddVertex(r*-0.26286500, r*0.0000000, r*-0.42532500)
	icosahedron.AddVertex(r*0.26286500, r*0.0000000, r*-0.42532500)
	icosahedron.AddVertex(r*0.0000000, r*0.42532500, r*0.26286500)
	icosahedron.AddVertex(r*0.0000000, r*0.42532500, r*-0.26286500)
	icosahedron.AddVertex(r*0.0000000, r*-0.42532500, r*0.26286500)
	icosahedron.AddVertex(r*0.0000000, r*-0.42532500, r*-0.26286500)
	icosahedron.AddVertex(r*0.42532500, r*0.26286500, r*0.0000000)
	icosahedron.AddVertex(r*-0.42532500, r*0.26286500, r*0.0000000)
	icosahedron.AddVertex(r*0.42532500, r*-0.26286500, r*0.0000000)
	icosahedron.AddVertex(r*-0.42532500, r*-0.26286500, r*0.0000000)
	return CollisionShape(icosahedron)
}
