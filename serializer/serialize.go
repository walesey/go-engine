package serializer

import (
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

var GlobalSerializer = Serializer{
	OnError: func(err error) {
		fmt.Println("serializer error: ", err)
	},
}

func SerializeArgs(args ...interface{}) []byte {
	return GlobalSerializer.SerializeArgs(args...)
}

func Stringfrombytes(r io.Reader) string {
	return GlobalSerializer.Stringfrombytes(r)
}

func Stringbytes(w io.Writer, str string) {
	GlobalSerializer.Stringbytes(w, str)
}

func Float64frombytes(r io.Reader) float64 {
	return GlobalSerializer.Float64frombytes(r)
}

func Float64bytes(w io.Writer, float float64) {
	GlobalSerializer.Float64bytes(w, float)
}

func Float32frombytes(r io.Reader) float32 {
	return GlobalSerializer.Float32frombytes(r)
}

func Float32bytes(w io.Writer, float float32) {
	GlobalSerializer.Float32bytes(w, float)
}

func UInt8frombytes(r io.Reader) uint8 {
	return GlobalSerializer.UInt8frombytes(r)
}

func UInt8Bytes(w io.Writer, i uint8) {
	GlobalSerializer.UInt8Bytes(w, i)
}

func UInt16frombytes(r io.Reader) uint16 {
	return GlobalSerializer.UInt16frombytes(r)
}

func UInt16Bytes(w io.Writer, i uint16) {
	GlobalSerializer.UInt16Bytes(w, i)
}

func UInt32frombytes(r io.Reader) uint32 {
	return GlobalSerializer.UInt32frombytes(r)
}

func UInt32Bytes(w io.Writer, i uint32) {
	GlobalSerializer.UInt32Bytes(w, i)
}

func UInt64frombytes(r io.Reader) uint64 {
	return GlobalSerializer.UInt64frombytes(r)
}

func UInt64Bytes(w io.Writer, i uint64) {
	GlobalSerializer.UInt64Bytes(w, i)
}

func Vector2frombytes(r io.Reader) mgl32.Vec2 {
	return GlobalSerializer.Vector2frombytes(r)
}

func Vector2bytes(w io.Writer, vector mgl32.Vec2) {
	GlobalSerializer.Vector2bytes(w, vector)
}

func Vector2frombytes64(r io.Reader) mgl64.Vec2 {
	return GlobalSerializer.Vector2frombytes64(r)
}

func Vector2bytes64(w io.Writer, vector mgl64.Vec2) {
	GlobalSerializer.Vector2bytes64(w, vector)
}

func Vector3frombytes(r io.Reader) mgl32.Vec3 {
	return GlobalSerializer.Vector3frombytes(r)
}

func Vector3bytes(w io.Writer, vector mgl32.Vec3) {
	GlobalSerializer.Vector3bytes(w, vector)
}

func Vector3frombytes64(r io.Reader) mgl64.Vec3 {
	return GlobalSerializer.Vector3frombytes64(r)
}

func Vector3bytes64(w io.Writer, vector mgl64.Vec3) {
	GlobalSerializer.Vector3bytes64(w, vector)
}

func BoolFromBytes(r io.Reader) bool {
	return GlobalSerializer.BoolFromBytes(r)
}

func BoolBytes(w io.Writer, b bool) {
	GlobalSerializer.BoolBytes(w, b)
}
