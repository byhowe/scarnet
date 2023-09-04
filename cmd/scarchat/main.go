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

	scarnet.WriteRequest(conn, &scarnet.SignupRequest{Creds: scarnet.AccountCredentials{Username: "username", Password: "password"}})
}
