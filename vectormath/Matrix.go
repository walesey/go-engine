package vectormath

type Matrix interface {
	Determinant() float64
}

type Matrix3 struct {
	M00, M01, M02 float64
	M10, M11, M12 float64
	M20, M21, M22 float64
}

func (m Matrix3) Determinant() float64 {
	return m.M00*m.M11*m.M22 + m.M01*m.M12*m.M20 + m.M02*m.M10*m.M21 - m.M02*m.M11*m.M20 - m.M01*m.M10*m.M22 - m.M00*m.M12*m.M21
}

func (m Matrix3) Inverse() Matrix3 {
	det := m.Determinant()
	if ApproxEqual(det, 0.0, 0.000001) {
		return Matrix3{}
	}
	retMat := Matrix3{
		M00: m.M11*m.M22 - m.M12*m.M21,
		M01: m.M02*m.M21 - m.M01*m.M22,
		M02: m.M01*m.M12 - m.M02*m.M11,
		M10: m.M12*m.M20 - m.M10*m.M22,
		M11: m.M00*m.M22 - m.M02*m.M20,
		M12: m.M02*m.M10 - m.M00*m.M12,
		M20: m.M10*m.M21 - m.M11*m.M20,
		M21: m.M01*m.M20 - m.M00*m.M21,
		M22: m.M00*m.M11 - m.M01*m.M10,
	}
	return retMat.MultiplyScalar(1 / det)
}

func (m Matrix3) MultiplyScalar(scalar float64) Matrix3 {
	return Matrix3{
		m.M00 * scalar, m.M01 * scalar, m.M02 * scalar,
		m.M10 * scalar, m.M11 * scalar, m.M12 * scalar,
		m.M20 * scalar, m.M21 * scalar, m.M22 * scalar,
	}
}

func (m Matrix3) MultiplyVector(v Vector3) Vector3 {
	return Vector3{
		X: m.M00*v.X + m.M10*v.Y + m.M20*v.Z,
		Y: m.M01*v.X + m.M11*v.Y + m.M21*v.Z,
		Z: m.M02*v.X + m.M12*v.Y + m.M22*v.Z,
	}
}

func (m Matrix3) Multiply(other Matrix3) Matrix3 {
	return Matrix3{
		M00: m.M00*other.M00 + m.M10*other.M01 + m.M20*other.M02,
		M01: m.M01*other.M00 + m.M11*other.M01 + m.M21*other.M02,
		M02: m.M02*other.M00 + m.M12*other.M01 + m.M22*other.M02,
		M10: m.M00*other.M10 + m.M10*other.M11 + m.M20*other.M12,
		M11: m.M01*other.M10 + m.M11*other.M11 + m.M21*other.M12,
		M12: m.M02*other.M10 + m.M12*other.M11 + m.M22*other.M12,
		M20: m.M00*other.M20 + m.M10*other.M21 + m.M20*other.M22,
		M21: m.M01*other.M20 + m.M11*other.M21 + m.M21*other.M22,
		M22: m.M02*other.M20 + m.M12*other.M21 + m.M22*other.M22,
	}
}

func (m Matrix3) Transform(v Vector3) Vector3 {
	return m.MultiplyVector(v)
}

type Matrix4 struct {
	M00, M01, M02, M03 float64
	M10, M11, M12, M13 float64
	M20, M21, M22, M23 float64
	M30, M31, M32, M33 float64
}

func (m Matrix4) Determinant() float64 {
	val1 := m.M11*m.M22*m.M33 + m.M12*m.M23*m.M31 + m.M13*m.M21*m.M32 - m.M13*m.M22*m.M31 - m.M12*m.M21*m.M33 - m.M11*m.M23*m.M32
	val2 := m.M10*m.M22*m.M33 + m.M12*m.M23*m.M30 + m.M13*m.M20*m.M32 - m.M13*m.M22*m.M30 - m.M12*m.M20*m.M33 - m.M10*m.M23*m.M32
	val3 := m.M10*m.M21*m.M33 + m.M11*m.M23*m.M30 + m.M13*m.M20*m.M31 - m.M13*m.M21*m.M30 - m.M11*m.M20*m.M33 - m.M10*m.M23*m.M31
	val4 := m.M10*m.M21*m.M32 + m.M11*m.M22*m.M30 + m.M12*m.M20*m.M31 - m.M12*m.M21*m.M30 - m.M11*m.M20*m.M32 - m.M10*m.M22*m.M31
	return m.M00*val1 - m.M01*val2 + m.M02*val3 - m.M03*val4
}
