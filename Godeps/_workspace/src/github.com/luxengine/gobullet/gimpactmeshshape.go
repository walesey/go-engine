package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

import (
	"unsafe"
)

type GImpactMeshShape struct {
	handle C.plGImpactMeshShape
}

func NewGImpactMeshShapeFromTriangleMesh(mesh TriangleMesh) GImpactMeshShape {
	return GImpactMeshShape{C.plNewGImpactMeshShapeFromTriangleMesh(mesh.handle)}
}

func NewGImpactMeshShapeFromTriangleMeshIndexed(mesh TriangleIndexedMesh) GImpactMeshShape {
	return GImpactMeshShape{C.plNewGImpactMeshShapeFromTriangleMeshIndexed(mesh.handle)}
}

func (this GImpactMeshShape) Delete() {
	C.plGImpactMeshShapeDelete(this.handle)
}

//UNTESTED
func (this GImpactMeshShape) RigidBody(userdata *interface{}, mass float32) RigidBody {
	return RigidBody{C.plCreateRigidBodyGImpactMeshShape(unsafe.Pointer(userdata), C.float(mass), this.handle)}
}

func (this GImpactMeshShape) UpdateBound() {
	C.plGImpactMeshShapeUpdateBound(this.handle)
}

func (this GImpactMeshShape) PostUpdate() {
	C.plGImpactMeshShapePostUpdate(this.handle)
}
