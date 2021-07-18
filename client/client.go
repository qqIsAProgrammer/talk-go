package client

import (
	"errors"
	"net"
	"reflect"
	"talk-go/serial"
	"talk-go/transport"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

// CallRPC method.
func (c *Client) CallRPC(rpcName string, fPtr interface{}) {
	container := reflect.ValueOf(fPtr).Elem()

	f := func(req []reflect.Value) []reflect.Value {
		// error handle func
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		reqTransport := transport.NewTransport(c.conn)
		// process input params
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}

		// request RPC
		reqRPC := serial.RPCData{Name: rpcName, Args: inArgs}
		b, err := serial.Encode(reqRPC)
		if err != nil {
			panic(err)
		}
		err = reqTransport.Send(b)
		if err != nil {
			return errorHandler(err)
		}

		// receive response from server
		rsp, err := reqTransport.Read()
		if err != nil {
			return errorHandler(err)
		}
		rspDecode, _ := serial.Decode(rsp)
		if rspDecode.Err != "" {
			return errorHandler(errors.New(rspDecode.Err))
		}

		if len(rspDecode.Args) == 0 {
			rspDecode.Args = make([]interface{}, container.Type().NumOut())
		}
		// unpack response arguments
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 {
				if rspDecode.Args[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rspDecode.Args[i])
				}
			} else {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}

		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), f))
}
