package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"talk-go/client"
	"talk-go/server"
	"time"
)

type User struct {
	Name string
	Age  int
}

var userDB = map[int]User{
	1: {"Kim", 20},
	2: {"Jackson", 35},
	3: {"Lily", 41},
}

func QueryUser(id int) (User, error) {
	if u, ok := userDB[id]; ok {
		return u, nil
	}

	return User{}, fmt.Errorf("id %d not in user db", id)
}

func main() {
	// new type needs to be registered
	gob.Register(User{})

	// server part
	addr := "localhost:6060"
	s := server.NewServer(addr)
	s.Register("QueryUser", QueryUser)
	go s.Run()

	// wait for server to start
	time.Sleep(time.Second)

	// client part
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	c := client.NewClient(conn)
	var Query func(int) (User, error)
	c.CallRPC("QueryUser", &Query)

	u1, err := Query(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(u1)

	u2, err := Query(8)
	if err != nil {
		panic(err)
	}
	fmt.Println(u2)
}
