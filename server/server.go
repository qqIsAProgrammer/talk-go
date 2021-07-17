package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"talk-go/serial"
	"talk-go/transport"
)

type RPCServer struct {
	addr string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *RPCServer {
	return &RPCServer{
		addr,
		make(map[string]reflect.Value),
	}
}

func (s *RPCServer) Register(fnName string, fFunc interface{}) {
	if _, ok := s.funcs[fnName]; ok {
		return
	}
	s.funcs[fnName] = reflect.ValueOf(fFunc)
}

func (s *RPCServer) execute(req serial.RPCData) serial.RPCData {
	f, ok := s.funcs[req.Name]
	if !ok {
		e := fmt.Sprintf("func %s not registered", req.Name)
		log.Println(e)
		return serial.RPCData{Name: req.Name, Args: nil, Err: e}
	}

	log.Printf("func %s is called\n", req.Name)
	inArgs := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		inArgs[i] = reflect.ValueOf(req.Args[i])
	}

	out := f.Call(inArgs)

	resArgs := make([]interface{}, len(out)-1)
	for i := 0; i < len(out)-1; i++ {
		resArgs[i] = out[i].Interface()
	}

	var err string
	if _, ok := out[len(out)-1].Interface().(error); ok {
		err = out[len(out)-1].Interface().(error).Error()
	}
	return serial.RPCData{Name: req.Name, Args: resArgs, Err: err}
}

func (s *RPCServer) Run() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("listen on %s err: %v\n", s.addr, err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}
		go func() {
			connTransport := transport.NewTransport(conn)
			for {
				// read request
				req, err := connTransport.Read()
				if err != nil {
					if err != io.EOF {
						log.Printf("read err: %v\n", err)
						return
					}
				}

				// decode the data and pass it to execute
				decReq, err := serial.Decode(req)
				if err != nil {
					log.Printf("Error Decoding the Payload err: %v\n", err)
					return
				}

				// get the executed result
				resp := s.execute(decReq)

				// encode the data back
				b, err := serial.Encode(resp)
				if err != nil {
					log.Printf("Error Encoding the Payload for response err: %v\n", err)
					return
				}

				// send response to client
				err = connTransport.Send(b)
				if err != nil {
					log.Printf("transport write err: %v\n", err)
				}
			}
		}()
	}
}