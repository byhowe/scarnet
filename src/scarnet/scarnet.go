package scarnet

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"golang.org/x/exp/slog"
)

var (
	_ Request = &SignupRequest{} // compile time proof
	_ Request = &LoginRequest{}  // compile time proof
)

const (
	ActionSignup = iota
	ActionLogin
	ActionMessage
)

type Request interface {
	ActionCode() uint32
}

type SignupRequest struct {
	Creds AccountCredentials
}

func (r *SignupRequest) ActionCode() uint32 {
	return ActionSignup
}

type LoginRequest struct {
	Creds AccountCredentials
}

func (r *LoginRequest) ActionCode() uint32 {
	return ActionLogin
}

type MessageRequest struct {
	Receiver string
	Message  string
}

func (r *MessageRequest) ActionCode() uint32 {
	return ActionMessage
}

type AccountCredentials struct {
	Username string
	Password string
}

func ReadRequest(conn net.Conn) (Request, error) {
	buffer := make([]byte, 4)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		slog.Error("read error in connection handler:", err)
		return nil, err
	}
	msgLen := binary.BigEndian.Uint32(buffer)

	buffer = make([]byte, msgLen + 4)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		slog.Error("read error in connection handler:", err)
		return nil, err
	}

	var request Request
	action := binary.BigEndian.Uint32(buffer[:4])
	switch action {
	case ActionSignup:
		var signupRequest SignupRequest
		err := json.Unmarshal(buffer[4:], &signupRequest)
		if err != nil {
			slog.Error("json unmarshal error:", err)
		}
		request = &signupRequest
	case ActionLogin:
		var loginRequest LoginRequest
		err := json.Unmarshal(buffer[4:], &loginRequest)
		if err != nil {
			slog.Error("json unmarshal error:", err)
		}
		request = &loginRequest
	case ActionMessage:
		var messageRequest MessageRequest
		err := json.Unmarshal(buffer[4:], &messageRequest)
		if err != nil {
			slog.Error("json unmarshal error:", err)
		}
		request = &messageRequest
	default:
		slog.Error("unknown action error:", action)
		return nil, fmt.Errorf("unknwon action error: %d", action)
	}

	return request, nil
}

func WriteRequest(conn net.Conn, request Request) error {
	var err error
	var data []byte

	switch request.ActionCode() {
	case ActionSignup:
		data, err = json.Marshal(request.(*SignupRequest))
		if err != nil {
			slog.Error("json unmarshal error:", err)
			return err
		}
	case ActionLogin:
		data, err = json.Marshal(request.(*LoginRequest))
		if err != nil {
			slog.Error("json unmarshal error:", err)
			return err
		}
	case ActionMessage:
		data, err = json.Marshal(request.(*MessageRequest))
		if err != nil {
			slog.Error("json unmarshal error:", err)
			return err
		}
	}

	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(len(data)))
	_, err = conn.Write(buffer)
	if err != nil {
		slog.Error("write error in connection handler:", err)
		return err
	}

	binary.BigEndian.PutUint32(buffer, request.ActionCode())
	_, err = conn.Write(buffer)
	if err != nil {
		slog.Error("write error in connection handler:", err)
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		slog.Error("write error in connection handler:", err)
		return err
	}

	return nil
}
