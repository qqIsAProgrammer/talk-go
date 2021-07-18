package serial

import (
	"bytes"
	"encoding/gob"
)

// RPCData represents the serializing format of structured data.
type RPCData struct {
	Name string        // name of the function
	Args []interface{} // request's or response's body expect error
	Err  string        // error any executing remote server
}

// Encode the RPCData in binary format which can
// be sent over the network.
func Encode(data RPCData) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode the binary data into the Go RPC struct.
func Decode(b []byte) (RPCData, error) {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	var data RPCData
	if err := dec.Decode(&data); err != nil {
		return RPCData{}, err
	}
	return data, nil
}
