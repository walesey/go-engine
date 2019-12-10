package serializer

import (
	"io"

	"github.com/walesey/go-engine/util"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Serializer struct {
	OnError func(error)
}

func (s Serializer) SerializeArgs(args ...interface{}) []byte {
	result, err := util.SerializeArgs(args...)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Stringfrombytes(r io.Reader) string {
	result, err := util.Stringfrombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Stringbytes(w io.Writer, str string) {
	if err := util.Stringbytes(w, str); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) Float64frombytes(r io.Reader) float64 {
	result, err := util.Float64frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Float64bytes(w io.Writer, float float64) {
	if err := util.Float64bytes(w, float); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) Float32frombytes(r io.Reader) float32 {
	result, err := util.Float32frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Float32bytes(w io.Writer, float float32) {
	if err := util.Float32bytes(w, float); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) UInt8frombytes(r io.Reader) uint8 {
	result, err := util.UInt8frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) UInt8Bytes(w io.Writer, i uint8) {
	if err := util.UInt8Bytes(w, i); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) UInt16frombytes(r io.Reader) uint16 {
	result, err := util.UInt16frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) UInt16Bytes(w io.Writer, i uint16) {
	if err := util.UInt16Bytes(w, i); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) UInt32frombytes(r io.Reader) uint32 {
	result, err := util.UInt32frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) UInt32Bytes(w io.Writer, i uint32) {
	if err := util.UInt32Bytes(w, i); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) UInt64frombytes(r io.Reader) uint64 {
	result, err := util.UInt64frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) UInt64Bytes(w io.Writer, i uint64) {
	if err := util.UInt64Bytes(w, i); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) Vector2frombytes(r io.Reader) mgl32.Vec2 {
	result, err := util.Vector2frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Vector2bytes(w io.Writer, vector mgl32.Vec2) {
	if err := util.Vector2bytes(w, vector); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) Vector2frombytes64(r io.Reader) mgl64.Vec2 {
	result, err := util.Vector2frombytes64(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Vector2bytes64(w io.Writer, vector mgl64.Vec2) {
	if err := util.Vector2bytes64(w, vector); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) Vector3frombytes(r io.Reader) mgl32.Vec3 {
	result, err := util.Vector3frombytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Vector3bytes(w io.Writer, vector mgl32.Vec3) {
	if err := util.Vector3bytes(w, vector); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) Vector3frombytes64(r io.Reader) mgl64.Vec3 {
	result, err := util.Vector3frombytes64(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) Vector3bytes64(w io.Writer, vector mgl64.Vec3) {
	if err := util.Vector3bytes64(w, vector); err != nil {
		s.OnError(err)
	}
}

func (s Serializer) BoolFromBytes(r io.Reader) bool {
	result, err := util.BoolFromBytes(r)
	if err != nil {
		s.OnError(err)
	}
	return result
}

func (s Serializer) BoolBytes(w io.Writer, b bool) {
	if err := util.BoolBytes(w, b); err != nil {
		s.OnError(err)
	}
}
