package util

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"strings"
)

func Base64ToBytes(base64String string) []byte {
	decoder := base64.NewDecoder(base64.URLEncoding, strings.NewReader(base64String))
	data, err := ioutil.ReadAll(decoder)
	if err != nil {
		fmt.Printf("Error converting base64 string to bytes: %v\n", err)
		return []byte{}
	}
	return data
}

func Serialize(object interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(object)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
