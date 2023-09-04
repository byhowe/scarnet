package main

import (
	"fmt"
	"log"
	"net"

	"github.com/byhowe/scarnet/src/scarnet"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:20058")
	if err != nil {
		log.Fatal("dial tcp error:", err)
	}
	defer conn.Close()

	fmt.Printf("connection to %s\n", conn.RemoteAddr().String())

	err = scarnet.WriteExchange(conn, &scarnet.SignupRequest{
		Username: "username",
		Password: "password",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = scarnet.WriteExchange(conn, &scarnet.LoginRequest{
		Username: "username",
		Password: "password",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = scarnet.WriteExchange(conn, &scarnet.LoginRequest{
		Username: "username",
		Password: "passwo",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = scarnet.WriteExchange(conn, &scarnet.LoginRequest{
		Username: "userna",
		Password: "password",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = scarnet.WriteExchange(conn, &scarnet.SignupRequest{
		Username: "username",
		Password: "password",
	})
	if err != nil {
		fmt.Println(err)
	}
	err = scarnet.WriteExchange(conn, &scarnet.MessageRequest{
		Receiver: "username",
		Message:  "hello from sender",
	})
	if err != nil {
		fmt.Println(err)
	}
}
