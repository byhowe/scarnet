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
	_ Exchange = &SignupRequest{} // compile time proof
	_ Exchange = &LoginRequest{}  // compile time proof
)

const (
	ExchangeIdSignupRequest ExchangeId = iota
	ExchangeIdLoginRequest
	ExchangeIdMessageRequest
)

type ExchangeId uint32

type Exchange interface {
	ExchangeId() ExchangeId
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *SignupRequest) ExchangeId() ExchangeId {
	return ExchangeIdSignupRequest
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) ExchangeId() ExchangeId {
	return ExchangeIdLoginRequest
}

type MessageRequest struct {
	Receiver string `json:"to"`
	Message  string `json:"msg"`
}

func (r *MessageRequest) ExchangeId() ExchangeId {
	return ExchangeIdMessageRequest
}

func ReadExchange(conn net.Conn) (Exchange, error) {
	// read exchange id
	buffer := make([]byte, 4)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		slog.Error("read error in connection handler:", err)
		return nil, err
	}
	exchangeId := ExchangeId(binary.BigEndian.Uint32(buffer))

	// read data length
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		slog.Error("read error in connection handler:", err)
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(buffer)

	// read data
	buffer = make([]byte, dataLen)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		slog.Error("read error in connection handler:", err)
		return nil, err
	}

	var exchange Exchange
	switch exchangeId {
	case ExchangeIdSignupRequest:
		var signupRequest SignupRequest
		err := json.Unmarshal(buffer, &signupRequest)
		if err != nil {
			slog.Error("json unmarshal error:", "signup", err)
			return nil, err
		}
		exchange = &signupRequest
	case ExchangeIdLoginRequest:
		var loginRequest LoginRequest
		err := json.Unmarshal(buffer, &loginRequest)
		if err != nil {
			slog.Error("json unmarshal error:", "login", err)
			return nil, err
		}
		exchange = &loginRequest
	case ExchangeIdMessageRequest:
		var messageRequest MessageRequest
		err := json.Unmarshal(buffer, &messageRequest)
		if err != nil {
			slog.Error("json unmarshal error:", "message", err)
			return nil, err
		}
		exchange = &messageRequest
	default:
		slog.Error("unknown action error:", exchangeId)
		return nil, fmt.Errorf("unknwon action error: %d", exchangeId)
	}

	return exchange, nil
}

func WriteExchange(conn net.Conn, exchange Exchange) error {
	var err error
	var data []byte

	switch exchange.ExchangeId() {
	case ExchangeIdSignupRequest:
		data, err = json.Marshal(exchange.(*SignupRequest))
		if err != nil {
			slog.Error("json unmarshal error:", err)
			return err
		}
	case ExchangeIdLoginRequest:
		data, err = json.Marshal(exchange.(*LoginRequest))
		if err != nil {
			slog.Error("json unmarshal error:", err)
			return err
		}
	case ExchangeIdMessageRequest:
		data, err = json.Marshal(exchange.(*MessageRequest))
		if err != nil {
			slog.Error("json unmarshal error:", err)
			return err
		}
	default:
		slog.Error("unknown exchange id:", exchange.ExchangeId())
	}

	// write exchange id
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(exchange.ExchangeId()))
	_, err = conn.Write(buffer)
	if err != nil {
		slog.Error("write error in connection handler:", err)
		return err
	}

	// write data length
	binary.BigEndian.PutUint32(buffer, uint32(len(data)))
	_, err = conn.Write(buffer)
	if err != nil {
		slog.Error("write error in connection handler:", err)
		return err
	}

	// write data
	_, err = conn.Write(data)
	if err != nil {
		slog.Error("write error in connection handler:", err)
		return err
	}

	return nil
}
