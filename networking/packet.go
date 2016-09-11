package networking

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Packet struct {
	Token   string
	Command string
	Data    []byte
}

func Encode(packet Packet) []byte {
	tokenLen := len(packet.Token)
	commandLen := len(packet.Command)
	dataLen := len(packet.Data)
	data := new(bytes.Buffer)
	data.WriteByte(byte(tokenLen))
	data.WriteByte(byte(commandLen))
	binary.Write(data, binary.LittleEndian, uint16(dataLen))
	data.WriteString(packet.Token)
	data.WriteString(packet.Command)
	data.Write(packet.Data)
	return data.Bytes()
}

func Decode(data []byte, i int) (Packet, error, int) {
	if len(data)-i < 4 {
		return Packet{}, fmt.Errorf("No data provided to Decode: len=%v", len(data)), len(data)
	}
	tokenLen := int(data[i])
	commandLen := int(data[i+1])
	dataLen := int(binary.LittleEndian.Uint16(data[i+2 : i+4]))
	i += 4
	token := string(data[i : i+tokenLen])
	i += tokenLen
	command := string(data[i : i+commandLen])
	i += commandLen
	packetData := data[i : i+dataLen]
	i += dataLen
	return Packet{
		Token:   token,
		Command: command,
		Data:    packetData,
	}, nil, i
}
