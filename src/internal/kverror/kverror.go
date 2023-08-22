package kverror

var (
	_ error   = (*Error)(nil) // compile time proof
	_ KVError = (*Error)(nil) // compile time proof
)

// sentinel errors.
var (
	ErrKeyExists   = New("key exist", true)
	ErrKeyNotFound = New("key not found", false)
	ErrUnknown     = New("unknown error", true)
)

// KVError defines custom error behaviours.
type KVError interface {
	Wrap(err error) KVError
	Unwrap() error
	AddData(any) KVError
	DestoryData() KVError
	Error() string
}

// Error is a type alias, custom error.
type Error struct {
	Err      error
	Message  string
	Data     any `json:"-"`
	Loggable bool
}

// AddData adds extra data to error.
func (e *Error) AddData(data any) KVError {
	e.Data = data
	return e
}

// Unwrap unwraps error.
func (e *Error) Unwrap() error {
	return e.Err
}

// DestoryData removes added data from error.
func (e *Error) DestoryData() KVError {
	e.Data = nil
	return e
}

// Wrap wraps given error.
func (e *Error) Wrap(err error) KVError {
	e.Err = err
	return e
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ", " + e.Message
	}
	return e.Message
}

// New instantiates new Error instance.
func New(m string, l bool) KVError {
	return &Error{
		Message:  m,
		Loggable: l,
	}
}
