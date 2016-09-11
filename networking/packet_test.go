package networking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	packet := Packet{
		Token:   "123",
		Command: "testCommand",
		Data:    []byte("test Data"),
	}
	expectedData := []byte{byte(len(packet.Token)), byte(len(packet.Command)), byte(len(packet.Data))}
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
	data := []byte{byte(len(expectedPacket.Token)), byte(len(expectedPacket.Command)), byte(len(expectedPacket.Data))}
	data = append(data, []byte(expectedPacket.Token)...)
	data = append(data, []byte(expectedPacket.Command)...)
	data = append(data, expectedPacket.Data...)

	packet, err, i := Decode(data, 0)
	assert.Nil(t, err, "decode should not return an error")
	assert.EqualValues(t, len(data), i, "Decode should return the correct read index")
	assert.EqualValues(t, expectedPacket, packet, "Decode packet didn't work")
}
