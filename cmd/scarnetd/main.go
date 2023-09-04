package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/byhowe/scarnet/src/scarnet"
	"golang.org/x/exp/slog"
)

type Scarnet interface {
	AddUser(username string, password string) bool
	CheckCredentials(username string, password string) bool
}

type Server struct {
	conns map[net.Conn]bool
}

func main() {
	listener, err := net.Listen("tcp", ":20058")
	if err != nil {
		log.Fatal("create tcp listener error:", err)
	}

	fmt.Printf("listening on %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("failed to accept new conn:", err)
		}
		slog.Info("connection accepted from:", conn.RemoteAddr().String())

		go func(conn net.Conn) {
			defer conn.Close()

			request, err := scarnet.ReadRequest(conn)
			if err != nil {
				slog.Error("read request error:", err)
				return
			}

			if val, ok := request.(*scarnet.SignupRequest); ok {
				fmt.Printf("Signup request: %+v\n", val)
			}

			if val, ok := request.(*scarnet.LoginRequest); ok {
				fmt.Printf("Login request: %+v\n", val)
			}

			if val, ok := request.(*scarnet.MessageRequest); ok {
				fmt.Printf("Message request: %+v\n", val)
			}
		}(conn)
	}
}
