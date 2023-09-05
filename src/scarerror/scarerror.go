package scarerror

var (
	_ error     = (*Error)(nil) // compile time proof
	_ ScarError = (*Error)(nil) // compile time proof
)

var (
	ErrUserDisconnected = New("user disconnected", false)
	ErrIo               = New("io error", true)
	ErrSerialization    = New("io error", true)
	ErrUnknown          = New("io error", true)
)

type ScarError interface {
	Wrap(err error) ScarError
	Unwrap() error
	SetData(any) ScarError
	Error() string
}

type Error struct {
	Err      error
	Message  string
	Data     any
	Loggable bool
}

func (e *Error) Wrap(err error) ScarError {
	e.Err = err
	return e
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) SetData(data any) ScarError {
	e.Data = data
	return e
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ", " + e.Message
	}
	return e.Message
}

func New(msg string, loggable bool) *Error {
	return &Error{
		Message:  msg,
		Loggable: loggable,
	}
}
