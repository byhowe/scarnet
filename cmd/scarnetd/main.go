package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/byhowe/scarnet/src/scarnet"
	"golang.org/x/exp/slog"
)

type Server struct {
	mu    sync.RWMutex
	users map[string]string // username: password
}

func NewServer() *Server {
	return &Server{
		users: map[string]string{},
	}
}

type Connection struct {
	conn net.Conn
}

func main() {
	listener, err := net.Listen("tcp", ":20058")
	if err != nil {
		log.Fatal("create tcp listener error:", err)
	}

	fmt.Printf("listening on %s\n", listener.Addr().String())

	server := NewServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("failed to accept new conn:", err)
		}
		slog.Info("connection accepted from:", conn.RemoteAddr().String())

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				request, err := scarnet.ReadExchange(conn)
				if err == io.EOF {
					slog.Info("disconnected:", conn.RemoteAddr().String())
					break
				}
				if err != nil {
					slog.Error("read request error:", err)
					break
				}

				if val, ok := request.(*scarnet.SignupRequest); ok {
					//server.mu.Lock()
					//defer server.mu.Unlock()

					if _, ok := server.users[val.Username]; !ok {
						server.users[val.Username] = val.Password
						slog.Info("created user:", "signup", val.Username)
					} else {
						slog.Info("user exists:", "signup", val.Username)
					}
				}

				if val, ok := request.(*scarnet.LoginRequest); ok {
					//server.mu.RLock()
					//defer server.mu.RUnlock()

					if _, ok := server.users[val.Username]; !ok {
						slog.Info("no user exists:", "login", val.Username)
						continue
					}

					if server.users[val.Username] == val.Password {
						slog.Info("logged in user:", "login", val.Username)
					} else {
						slog.Info("incorrect password:", "login", val.Username)
					}
				}

				if val, ok := request.(*scarnet.MessageRequest); ok {
					slog.Info("message:", "message", val.Message)
				}
			}
		}(conn)
	}
}
