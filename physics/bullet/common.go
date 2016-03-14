package bullet

import vmath "github.com/walesey/go-engine/vectormath"

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
