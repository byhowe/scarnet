package scarnet

import "encoding/json"

var (
	_ Exchange = &SignupRequest{}  // compile time proof
	_ Exchange = &LoginRequest{}   // compile time proof
	_ Exchange = &MessageRequest{} // compile time proof
)

type ExchangeId uint32

type Exchange interface {
	ExchangeId() ExchangeId
	Marshal() ([]byte, error)
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *SignupRequest) ExchangeId() ExchangeId {
	return ExchangeIdSignupRequest
}

func (r *SignupRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) ExchangeId() ExchangeId {
	return ExchangeIdLoginRequest
}

func (r *LoginRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MessageRequest struct {
	Receiver string `json:"to"`
	Message  string `json:"msg"`
}

func (r *MessageRequest) ExchangeId() ExchangeId {
	return ExchangeIdMessageRequest
}

func (r *MessageRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
