package networking

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	packet := Packet{
		Token:   "123",
		Command: "testCommand",
		Data:    []byte("test Data"),
	}
	var bytes [2]byte
	binary.LittleEndian.PutUint16(bytes[:], uint16(len(packet.Data)))
	expectedData := []byte{byte(len(packet.Token)), byte(len(packet.Command)), bytes[0], bytes[1]}
	expectedData = append(expectedData, []byte(packet.Token)...)
	expectedData = append(expectedData, []byte(packet.Command)...)
	expectedData = append(expectedData, packet.Data...)

	data := Encode(packet)
	assert.EqualValues(t, expectedData, data, "Encode packet didn't work")
}

func TestDecode(t *testing.T) {
	expectedPacket := Packet{
		Token:   "123",
		Command: "testCommand",
		Data:    []byte("test Data"),
	}
	data := Encode(expectedPacket)

	packet, err, i := Decode(data, 0)
	assert.Nil(t, err, "decode should not return an error")
	assert.EqualValues(t, len(data), i, "Decode should return the correct read index")
	assert.EqualValues(t, expectedPacket, packet, "Decode packet didn't work")
}

func TestDecodeMultiple(t *testing.T) {
	testPacket1 := Packet{
		Token:   "123",
		Command: "testCommand",
		Data:    []byte("test Data"),
	}
	testPacket2 := Packet{
		Token:   "12345678",
		Command: "cmdTest",
		Data:    []byte("123 4567 abcd"),
	}

	data := append(Encode(testPacket1), Encode(testPacket2)...)

	packet, err, i := Decode(data, 0)
	assert.Nil(t, err, "decode should not return an error")
	assert.EqualValues(t, testPacket1, packet, "Decode first packet didn't work")

	packet, err, i = Decode(data, i)
	assert.EqualValues(t, len(data), i, "Decode should return the correct read index")
	assert.EqualValues(t, testPacket2, packet, "Decode second packet didn't work")
}
