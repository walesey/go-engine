package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func SerializeArgs(args ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			Stringbytes(buf, v)
		case int8:
			UInt8Bytes(buf, uint8(v))
		case uint8:
			UInt8Bytes(buf, uint8(v))
		case int16:
			UInt16Bytes(buf, uint16(v))
		case uint16:
			UInt16Bytes(buf, uint16(v))
		case int32:
			UInt32Bytes(buf, uint32(v))
		case uint32:
			UInt32Bytes(buf, uint32(v))
		case int64:
			UInt64Bytes(buf, uint64(v))
		case uint64:
			UInt64Bytes(buf, uint64(v))
		case int:
			UInt32Bytes(buf, uint32(v))
		case float32:
			Float32bytes(buf, v)
		case float64:
			Float64bytes(buf, v)
		case mgl32.Vec2:
			Vector2bytes(buf, v)
		case mgl32.Vec3:
			Vector3bytes(buf, v)
		case mgl64.Vec2:
			Vector2bytes64(buf, v)
		case mgl64.Vec3:
			Vector3bytes64(buf, v)
		case bool:
			BoolBytes(buf, v)
		case []byte:
			UInt32Bytes(buf, uint32(len(v)))
			buf.Write(v)
		default:
			fmt.Printf("Unknown typed used in SerializeArgs: %T \n", v)
		}
	}
	return buf.Bytes()
}

func Stringfrombytes(r io.Reader) string {
	var strLenData [1]byte
	r.Read(strLenData[:])
	strLen := int(strLenData[0])
	strData := make([]byte, strLen)
	r.Read(strData)
	return string(strData)
}

func Stringbytes(w io.Writer, str string) {
	w.Write([]byte{byte(len(str))})
	io.WriteString(w, str)
}

func Float64frombytes(r io.Reader) float64 {
	var bytes [8]byte
	r.Read(bytes[:])
	bits := binary.LittleEndian.Uint64(bytes[:])
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(w io.Writer, float float64) {
	bits := math.Float64bits(float)
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], bits)
	w.Write(bytes[:])
}

func Float32frombytes(r io.Reader) float32 {
	var bytes [4]byte
	r.Read(bytes[:])
	bits := binary.LittleEndian.Uint32(bytes[:])
	float := math.Float32frombits(bits)
	return float
}

func Float32bytes(w io.Writer, float float32) {
	bits := math.Float32bits(float)
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], bits)
	w.Write(bytes[:])
}

func UInt8frombytes(r io.Reader) uint8 {
	var bytes [1]byte
	r.Read(bytes[:])
	return uint8(bytes[0])
}

func UInt8Bytes(w io.Writer, i uint8) {
	w.Write([]byte{byte(i)})
}

func UInt16frombytes(r io.Reader) uint16 {
	var bytes [2]byte
	r.Read(bytes[:])
	return binary.LittleEndian.Uint16(bytes[:])
}

func UInt16Bytes(w io.Writer, i uint16) {
	var bytes [2]byte
	binary.LittleEndian.PutUint16(bytes[:], i)
	w.Write(bytes[:])
}

func UInt32frombytes(r io.Reader) uint32 {
	var bytes [4]byte
	r.Read(bytes[:])
	return binary.LittleEndian.Uint32(bytes[:])
}

func UInt32Bytes(w io.Writer, i uint32) {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], i)
	w.Write(bytes[:])
}

func UInt64frombytes(r io.Reader) uint64 {
	var bytes [8]byte
	r.Read(bytes[:])
	return binary.LittleEndian.Uint64(bytes[:])
}

func UInt64Bytes(w io.Writer, i uint64) {
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], i)
	w.Write(bytes[:])
}

func Vector2frombytes(r io.Reader) mgl32.Vec2 {
	x := Float32frombytes(r)
	y := Float32frombytes(r)
	return mgl32.Vec2{x, y}
}

func Vector2bytes(w io.Writer, vector mgl32.Vec2) {
	Float32bytes(w, vector.X())
	Float32bytes(w, vector.Y())
}

func Vector2frombytes64(r io.Reader) mgl64.Vec2 {
	x := Float64frombytes(r)
	y := Float64frombytes(r)
	return mgl64.Vec2{x, y}
}

func Vector2bytes64(w io.Writer, vector mgl64.Vec2) {
	Float64bytes(w, vector.X())
	Float64bytes(w, vector.Y())
}

func Vector3frombytes(r io.Reader) mgl32.Vec3 {
	x := Float32frombytes(r)
	y := Float32frombytes(r)
	z := Float32frombytes(r)
	return mgl32.Vec3{x, y, z}
}

func Vector3bytes(w io.Writer, vector mgl32.Vec3) {
	Float32bytes(w, vector.X())
	Float32bytes(w, vector.Y())
	Float32bytes(w, vector.Z())
}

func Vector3frombytes64(r io.Reader) mgl64.Vec3 {
	x := Float64frombytes(r)
	y := Float64frombytes(r)
	z := Float64frombytes(r)
	return mgl64.Vec3{x, y, z}
}

func Vector3bytes64(w io.Writer, vector mgl64.Vec3) {
	Float64bytes(w, vector.X())
	Float64bytes(w, vector.Y())
	Float64bytes(w, vector.Z())
}

func BoolFromBytes(r io.Reader) bool {
	var bytes [1]byte
	r.Read(bytes[:])
	return bytes[0] == 1
}

func BoolBytes(w io.Writer, b bool) {
	if b {
		w.Write([]byte{1})
	} else {
		w.Write([]byte{0})
	}
}
