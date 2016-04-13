package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

//UNTESTED
//Transform euler rotation to quaternion
func SetEuler(yaw, pitch, roll float32, quat *[4]float32) {
	C.plSetEuler(C.plReal(yaw), C.plReal(pitch), C.plReal(roll), (*C.plReal)(&quat[0]))
}
