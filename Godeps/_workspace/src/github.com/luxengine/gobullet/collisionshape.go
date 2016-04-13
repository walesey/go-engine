package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

import (
	"unsafe"
)

/**************************
******Collision Shape******
**************************/

type CollisionShape struct {
	handle C.plCollisionShape
}

type ConvexHull CollisionShape
type CompoundShape CollisionShape
type TriangleMeshShape CollisionShape

//Free memory of this CollisionShape
//DONT DELETE IF SOME RIGIDBODY ARE STILL USING IT
func (this CollisionShape) Delete() {
	C.plDeleteShape(this.handle)
}

//overload for ConvexHull
func (this ConvexHull) Delete() {
	C.plDeleteShape(this.handle)
}

//overload for CompoundShape
func (this CompoundShape) Delete() {
	C.plDeleteShape(this.handle)
}

func (this TriangleMeshShape) Delete() {
	C.plDeleteShape(this.handle)
}

//UNTESTED
func (this CollisionShape) RigidBody(userdata unsafe.Pointer, mass float32) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), this.handle)}
}

//UNTESTED
func (this ConvexHull) RigidBody(userdata unsafe.Pointer, mass float32) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), this.handle)}
}

//UNTESTED
func (this CompoundShape) RigidBody(userdata unsafe.Pointer, mass float32) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), this.handle)}
}

//UNTESTED
func (this TriangleMeshShape) RigidBody(userdata unsafe.Pointer, mass float32) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), this.handle)}
}

//UNTESTED
func (this CollisionShape) SetScaling(scale *[3]float32) {
	C.plSetScaling(this.handle, (*C.plReal)(&scale[0]))
}

//UNTESTED
func (this CompoundShape) SetScaling(scale *[3]float32) {
	C.plSetScaling(this.handle, (*C.plReal)(&scale[0]))
}

//UNTESTED
func (this ConvexHull) SetScaling(scale *[3]float32) {
	C.plSetScaling(this.handle, (*C.plReal)(&scale[0]))
}

//UNTESTED
func (this TriangleMeshShape) SetScaling(scale *[3]float32) {
	C.plSetScaling(this.handle, (*C.plReal)(&scale[0]))
}

//Creates a Sphere shape
//radius: the radius of the sphere...//kinda self explanatory
func NewSphereShape(radius float32) CollisionShape {
	return CollisionShape{C.plNewSphereShape(C.plReal(radius))}
}

//Create a box shape
//x,y,z is the width, height and length of the box, the origin is the center of the box, aka if it is located at (0,0,0) it will
//covers from (-x/2, -y/2, -z/2) to (x/2, y/2, z/2)
func NewBoxShape(x, y, z float32) CollisionShape {
	return CollisionShape{C.plNewBoxShape(C.plReal(x), C.plReal(y), C.plReal(z))}
}

//Create a Capsule shape
//the top and bottom are semi spheres
func NewCapsuleShape(radius, height float32) CollisionShape {
	return CollisionShape{C.plNewCapsuleShape(C.plReal(radius), C.plReal(height))}
}

//Create a cone shape
func NewConeShape(radius, height float32) CollisionShape {
	return CollisionShape{C.plNewConeShape(C.plReal(radius), C.plReal(height))}
}

//Create a cylinder shape
func NewCylinderShape(radius, height float32) CollisionShape {
	return CollisionShape{C.plNewCylinderShape(C.plReal(radius), C.plReal(height))}
}

//Create a Compound shape, a shape made of shapes
func NewCompoundShape() CompoundShape {
	return CompoundShape{C.plNewCompoundShape()}
}

//childPos = {x,y,z}
//childOri = quaternion {w,{x,y,z}}
//UNTESTED
func (this CompoundShape) AddChildShape(other CollisionShape, childPos *[3]float32, childOri *[4]float32) {
	C.plAddChildShape(this.handle, other.handle, (*C.plReal)(&childPos[0]), (*C.plReal)(&childOri[0]))
}

//Convex mesh
//see http://en.wikipedia.org/wiki/Convex_hull
func NewConvexHullShape() ConvexHull {
	return ConvexHull{C.plNewConvexHullShape()}
}

//Add a single vertex to the Convex Hull
func (this ConvexHull) AddVertex(x, y, z float32) {
	C.plAddVertex(this.handle, C.plReal(x), C.plReal(y), C.plReal(z))
}

type TriangleMesh struct {
	handle C.plTriangleMesh
}

type TriangleIndexedMesh struct {
	handle C.plTriangleIndexedMesh
}

//UNTESTED
func NewTriangleMesh(use32bitindices, use4componentvertices bool) TriangleMesh {
	return TriangleMesh{C.plNewTriangleMesh(C.bool(use32bitindices), C.bool(use4componentvertices))}
}

//int numTriangles ,int* triangleIndexBase, int triangleIndexStride, int numVertices, float* vertexBase, int vertexStride
func NewTriangleMeshIndexed(NumTriangles int, TriangleIndexBase []int, TriangleIndexStride, NumVertices int, VertexBase []float32, VertexStride int) TriangleIndexedMesh {
	return TriangleIndexedMesh{C.plNewTriangleMeshIndexed((C.int)(NumTriangles),
		(*C.int)(unsafe.Pointer(&TriangleIndexBase[0])), (C.int)(TriangleIndexStride), (C.int)(NumVertices),
		(*C.float)(unsafe.Pointer(&VertexBase[0])), (C.int)(VertexStride))}
}

//UNTESTED
func (this TriangleMesh) AddTriangle(vec1, vec2, vec3 *[3]float32, removeDuplicates bool) {
	C.plTriangleMeshAddTriangle(this.handle, (*C.plReal)(&vec1[0]), (*C.plReal)(&vec2[0]), (*C.plReal)(&vec3[0]), C.bool(removeDuplicates))
}

//UNTESTED
func (this TriangleMesh) NewTriangleMeshShape(useQuantizedAabbCompression, buildBvh bool) TriangleMeshShape {
	return TriangleMeshShape{C.plNewTriangleShape(this.handle, C.bool(useQuantizedAabbCompression), C.bool(buildBvh))}
}

//UNTESTED
func (this TriangleIndexedMesh) NewTriangleMeshShape(useQuantizedAabbCompression, buildBvh bool) TriangleMeshShape {
	return TriangleMeshShape{C.plNewTriangleShapeIndexed(this.handle, C.bool(useQuantizedAabbCompression), C.bool(buildBvh))}
}

//UNTESTED
func (this TriangleMesh) GetNumTriangles() int {
	return int(C.plTriangleMeshGetNumTriangles(this.handle))
}
