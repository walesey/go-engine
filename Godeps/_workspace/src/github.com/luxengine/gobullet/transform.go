package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

/********************************
************Transform************
********************************/
type Transform struct {
	handle C.plTransform
}

//UNTESTED
func NewTransform() Transform {
	return Transform{C.plNewTransform()}
}

//UNTESTED
func (this Transform) Delete() {
	C.plTransformDelete(this.handle)
}

//UNTESTED
func (this Transform) GetOrigin(dest *[3]float32) {
	C.plTransformGetOrigin(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Transform) GetRotation(dest *[4]float32) {
	C.plTransformGetRotation(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Transform) SetFromOpenGLMatrix(dest [16]float32) {
	C.plTransformSetFromOpenGLMatrix(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Transform) GetOpenGLMatrix(dest [16]float32) {
	C.plTransformGetOpenGLMatrix(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Transform) SetOrigin(dest *[3]float32) {
	C.plTransformSetOrigin(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Transform) SetRotation(dest *[4]float32) {
	C.plTransformSetRotation(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Transform) SetIdentity() {
	C.plTransformSetIdentity(this.handle)
}
