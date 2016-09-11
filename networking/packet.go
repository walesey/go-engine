package networking

import (
	"bytes"
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
	var data bytes.Buffer
	data.WriteByte(byte(tokenLen))
	data.WriteByte(byte(commandLen))
	data.WriteByte(byte(dataLen))
	data.WriteString(packet.Token)
	data.WriteString(packet.Command)
	data.Write(packet.Data)
	return data.Bytes()
}

func Decode(data []byte, i int) (Packet, error, int) {
	if len(data)-i < 3 {
		return Packet{}, fmt.Errorf("No data provided to Decode: len=%v", len(data)), len(data)
	}
	tokenLen := int(data[0])
	commandLen := int(data[1])
	dataLen := int(data[2])
	i += 3
	token := string(data[i : i+tokenLen])
	i += tokenLen
	command := string(data[i : i+commandLen])
	i += commandLen
	return Packet{
		Token:   token,
		Command: command,
		Data:    data[i : i+dataLen],
	}, nil, i + dataLen
}
