package kverror

var _ error = (*Error)(nil) // compile time proof

// sentinel errors.
var (
	ErrKeyExists   = New("key exist", true)
	ErrKeyNotFound = New("key not found", false)
)

// Error is a type alias, custom error.
type Error struct {
	Err      error
	Message  string
	Loggable bool
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ", " + e.Message
	}
	return e.Message
}

// New instantiates new Error instance.
func New(m string, l bool) error {
	return &Error{
		Message:  m,
		Loggable: l,
	}
}
