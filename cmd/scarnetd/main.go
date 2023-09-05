package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/byhowe/scarnet/src/scarerror"
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

func (s *Server) CheckUserCredentials(req *scarnet.LoginRequest) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.users[req.Username]; !ok {
		slog.Info("no user exists:", "login", req.Username)
		return false
	}

	if s.users[req.Username] == req.Password {
		slog.Info("logged in user:", "login", req.Username)
		return true
	} else {
		slog.Info("incorrect password:", "login", req.Username)
	}

	return false
}

func (s *Server) CreateUser(req *scarnet.SignupRequest) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[req.Username]; !ok {
		s.users[req.Username] = req.Password
		slog.Info("created user:", "signup", req.Username)
		return true
	} else {
		slog.Info("user exists:", "signup", req.Username)
	}

	return false
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

				if err != nil {
					if errors.Is(err, scarerror.ErrUserDisconnected) {
						slog.Info("user disconnected:", "loop", conn.RemoteAddr().String())
						break
					} else {
						slog.Error("read request error:", err)
						break
					}
				}

				if req, ok := request.(*scarnet.SignupRequest); ok {
					server.CreateUser(req)
				}

				if req, ok := request.(*scarnet.LoginRequest); ok {
					server.CheckUserCredentials(req)
				}

				if req, ok := request.(*scarnet.MessageRequest); ok {
					slog.Info("message:", "message", req.Message)
				}
			}
		}(conn)
	}
}
