package scarnet

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/byhowe/scarnet/src/scarerror"
	"golang.org/x/exp/slog"
)

func checkIoError(err error) error {
	if err == io.EOF {
		return scarerror.ErrUserDisconnected
	}
	return scarerror.ErrIo.Wrap(err)
}

func ReadExchange(conn net.Conn) (Exchange, error) {
	// read exchange id
	buffer := make([]byte, 4)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		return nil, checkIoError(err)
	}
	exchangeId := ExchangeId(binary.BigEndian.Uint32(buffer))

	// read data length
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, checkIoError(err)
	}
	dataLen := binary.BigEndian.Uint32(buffer)

	// read data
	buffer = make([]byte, dataLen)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, checkIoError(err)
	}

	var exchange Exchange
	switch exchangeId {
	case ExchangeIdSignupRequest:
		var signupRequest SignupRequest
		err := json.Unmarshal(buffer, &signupRequest)
		if err != nil {
			return nil, scarerror.ErrSerialization.Wrap(err)
		}
		exchange = &signupRequest
	case ExchangeIdLoginRequest:
		var loginRequest LoginRequest
		err := json.Unmarshal(buffer, &loginRequest)
		if err != nil {
			return nil, scarerror.ErrSerialization.Wrap(err)
		}
		exchange = &loginRequest
	case ExchangeIdMessageRequest:
		var messageRequest MessageRequest
		err := json.Unmarshal(buffer, &messageRequest)
		if err != nil {
			return nil, scarerror.ErrSerialization.Wrap(err)
		}
		exchange = &messageRequest
	default:
		return nil, scarerror.ErrUnknown.Wrap(fmt.Errorf("unknwon exchange id: %d", exchangeId))
	}

	return exchange, nil
}

func WriteExchange(conn net.Conn, exchange Exchange) error {
	var err error
	var data []byte

	data, err = exchange.Marshal()
	if err != nil {
		return scarerror.ErrSerialization.Wrap(err)
	}

	// write exchange id
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(exchange.ExchangeId()))
	_, err = conn.Write(buffer)
	if err != nil {
		return checkIoError(err)
	}

	// write data length
	binary.BigEndian.PutUint32(buffer, uint32(len(data)))
	_, err = conn.Write(buffer)
	if err != nil {
		return checkIoError(err)
	}

	// write data
	_, err = conn.Write(data)
	if err != nil {
		slog.Error("write error in connection handler:", err)
	}

	return nil
}
