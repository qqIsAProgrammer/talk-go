package serial

import (
	"bytes"
	"encoding/gob"
)

type RPCData struct {
	Name string
	Args []interface{}
	Err  string
}

func Encode(data RPCData) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (RPCData, error) {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	var data RPCData
	if err := dec.Decode(&data); err != nil {
		return RPCData{}, err
	}
	return data, nil
}
