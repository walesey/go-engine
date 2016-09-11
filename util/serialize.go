package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	vmath "github.com/walesey/go-engine/vectormath"
)

func SerializeArgs(args ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			Stringbytes(buf, v)
		case int:
			UInt32Bytes(buf, uint32(v))
		case float64:
			Float64bytes(buf, v)
		case vmath.Vector2:
			Vector2bytes(buf, v)
		case vmath.Vector3:
			Vector3bytes(buf, v)
		case bool:
			buf.WriteByte(BoolByte(v))
		case []byte:
			UInt32Bytes(buf, uint32(len(v)))
			buf.Write(v)
		default:
			fmt.Println("Unknown typed used in SerializeArgs")
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

func Vector2frombytes(r io.Reader) vmath.Vector2 {
	x := Float64frombytes(r)
	y := Float64frombytes(r)
	return vmath.Vector2{X: x, Y: y}
}

func Vector2bytes(w io.Writer, vector vmath.Vector2) {
	Float64bytes(w, vector.X)
	Float64bytes(w, vector.Y)
}

func Vector3frombytes(r io.Reader) vmath.Vector3 {
	x := Float64frombytes(r)
	y := Float64frombytes(r)
	z := Float64frombytes(r)
	return vmath.Vector3{X: x, Y: y, Z: z}
}

func Vector3bytes(w io.Writer, vector vmath.Vector3) {
	Float64bytes(w, vector.X)
	Float64bytes(w, vector.Y)
	Float64bytes(w, vector.Z)
}

func BoolFromByte(b byte) bool {
	return b == 1
}

func BoolByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
