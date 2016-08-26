package networking

import "fmt"

type Packet struct {
	Token   string
	Command string
	Data    []byte
}

func Encode(packet Packet) []byte {
	tokenLen := len(packet.Token)
	commandLen := len(packet.Command)
	data := make([]byte, 2, 2+tokenLen+commandLen+len(packet.Data))
	data[0] = byte(tokenLen)
	data[1] = byte(commandLen)
	data = append(data, []byte(packet.Token)...)
	data = append(data, []byte(packet.Command)...)
	return append(data, packet.Data...)
}

func Decode(data []byte) (Packet, error) {
	if len(data) < 2 {
		return Packet{}, fmt.Errorf("No data provided to Decode: len=%v", len(data))
	}
	tokenLen := int(data[0])
	commandLen := int(data[1])
	i := 2
	token := string(data[i : i+tokenLen])
	i += tokenLen
	command := string(data[i : i+commandLen])
	i += commandLen
	return Packet{
		Token:   token,
		Command: command,
		Data:    data[i:],
	}, nil
}
