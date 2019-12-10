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

func SerializeArgs(args ...interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			err = Stringbytes(buf, v)
		case int8:
			err = UInt8Bytes(buf, uint8(v))
		case uint8:
			err = UInt8Bytes(buf, uint8(v))
		case int16:
			err = UInt16Bytes(buf, uint16(v))
		case uint16:
			err = UInt16Bytes(buf, uint16(v))
		case int32:
			err = UInt32Bytes(buf, uint32(v))
		case uint32:
			err = UInt32Bytes(buf, uint32(v))
		case int64:
			err = UInt64Bytes(buf, uint64(v))
		case uint64:
			err = UInt64Bytes(buf, uint64(v))
		case int:
			err = UInt32Bytes(buf, uint32(v))
		case float32:
			err = Float32bytes(buf, v)
		case float64:
			err = Float64bytes(buf, v)
		case mgl32.Vec2:
			err = Vector2bytes(buf, v)
		case mgl32.Vec3:
			err = Vector3bytes(buf, v)
		case mgl64.Vec2:
			err = Vector2bytes64(buf, v)
		case mgl64.Vec3:
			err = Vector3bytes64(buf, v)
		case bool:
			err = BoolBytes(buf, v)
		case []byte:
			err = UInt32Bytes(buf, uint32(len(v)))
			if err == nil {
				_, err = buf.Write(v)
			}
		default:
			fmt.Printf("Unknown typed used in SerializeArgs: %T \n", v)
		}
	}
	return buf.Bytes(), err
}

func Stringfrombytes(r io.Reader) (string, error) {
	var strLenData [1]byte
	if _, err := r.Read(strLenData[:]); err != nil {
		return "", err
	}
	strLen := int(strLenData[0])
	strData := make([]byte, strLen)
	_, err := r.Read(strData)
	return string(strData), err
}

func Stringbytes(w io.Writer, str string) error {
	w.Write([]byte{byte(len(str))})
	_, err := io.WriteString(w, str)
	return err
}

func Float64frombytes(r io.Reader) (float64, error) {
	var bytes [8]byte
	_, err := r.Read(bytes[:])
	bits := binary.LittleEndian.Uint64(bytes[:])
	float := math.Float64frombits(bits)
	return float, err
}

func Float64bytes(w io.Writer, float float64) error {
	bits := math.Float64bits(float)
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], bits)
	_, err := w.Write(bytes[:])
	return err
}

func Float32frombytes(r io.Reader) (float32, error) {
	var bytes [4]byte
	_, err := r.Read(bytes[:])
	bits := binary.LittleEndian.Uint32(bytes[:])
	float := math.Float32frombits(bits)
	return float, err
}

func Float32bytes(w io.Writer, float float32) error {
	bits := math.Float32bits(float)
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], bits)
	_, err := w.Write(bytes[:])
	return err
}

func UInt8frombytes(r io.Reader) (uint8, error) {
	var bytes [1]byte
	_, err := r.Read(bytes[:])
	return uint8(bytes[0]), err
}

func UInt8Bytes(w io.Writer, i uint8) error {
	_, err := w.Write([]byte{byte(i)})
	return err
}

func UInt16frombytes(r io.Reader) (uint16, error) {
	var bytes [2]byte
	_, err := r.Read(bytes[:])
	return binary.LittleEndian.Uint16(bytes[:]), err
}

func UInt16Bytes(w io.Writer, i uint16) error {
	var bytes [2]byte
	binary.LittleEndian.PutUint16(bytes[:], i)
	_, err := w.Write(bytes[:])
	return err
}

func UInt32frombytes(r io.Reader) (uint32, error) {
	var bytes [4]byte
	_, err := r.Read(bytes[:])
	return binary.LittleEndian.Uint32(bytes[:]), err
}

func UInt32Bytes(w io.Writer, i uint32) error {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], i)
	_, err := w.Write(bytes[:])
	return err
}

func UInt64frombytes(r io.Reader) (uint64, error) {
	var bytes [8]byte
	_, err := r.Read(bytes[:])
	return binary.LittleEndian.Uint64(bytes[:]), err
}

func UInt64Bytes(w io.Writer, i uint64) error {
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], i)
	_, err := w.Write(bytes[:])
	return err
}

func Vector2frombytes(r io.Reader) (mgl32.Vec2, error) {
	x, err := Float32frombytes(r)
	if err != nil {
		return mgl32.Vec2{}, err
	}
	y, err := Float32frombytes(r)
	return mgl32.Vec2{x, y}, err
}

func Vector2bytes(w io.Writer, vector mgl32.Vec2) error {
	if err := Float32bytes(w, vector.X()); err != nil {
		return err
	}
	return Float32bytes(w, vector.Y())
}

func Vector2frombytes64(r io.Reader) (mgl64.Vec2, error) {
	x, err := Float64frombytes(r)
	if err != nil {
		return mgl64.Vec2{}, err
	}
	y, err := Float64frombytes(r)
	return mgl64.Vec2{x, y}, err
}

func Vector2bytes64(w io.Writer, vector mgl64.Vec2) error {
	if err := Float64bytes(w, vector.X()); err != nil {
		return err
	}
	return Float64bytes(w, vector.Y())
}

func Vector3frombytes(r io.Reader) (mgl32.Vec3, error) {
	x, err := Float32frombytes(r)
	if err != nil {
		return mgl32.Vec3{}, err
	}
	y, err := Float32frombytes(r)
	if err != nil {
		return mgl32.Vec3{}, err
	}
	z, err := Float32frombytes(r)
	return mgl32.Vec3{x, y, z}, err
}

func Vector3bytes(w io.Writer, vector mgl32.Vec3) error {
	if err := Float32bytes(w, vector.X()); err != nil {
		return err
	}
	if err := Float32bytes(w, vector.Y()); err != nil {
		return err
	}
	return Float32bytes(w, vector.Z())
}

func Vector3frombytes64(r io.Reader) (mgl64.Vec3, error) {
	x, err := Float64frombytes(r)
	if err != nil {
		return mgl64.Vec3{}, err
	}
	y, err := Float64frombytes(r)
	if err != nil {
		return mgl64.Vec3{}, err
	}
	z, err := Float64frombytes(r)
	return mgl64.Vec3{x, y, z}, err
}

func Vector3bytes64(w io.Writer, vector mgl64.Vec3) error {
	if err := Float64bytes(w, vector.X()); err != nil {
		return err
	}
	if err := Float64bytes(w, vector.Y()); err != nil {
		return err
	}
	return Float64bytes(w, vector.Z())
}

func BoolFromBytes(r io.Reader) (bool, error) {
	var bytes [1]byte
	_, err := r.Read(bytes[:])
	return bytes[0] == 1, err
}

func BoolBytes(w io.Writer, b bool) error {
	var err error
	if b {
		_, err = w.Write([]byte{1})
	} else {
		_, err = w.Write([]byte{0})
	}
	return err
}
