package gwt

import "errors"

// ErrNothingRegistered is returned if a Register is empty.
var ErrNothingRegistered = errors.New("not registered")

// Register is a dummy type.
type Register struct {
	value string
}

// NewRegister returns an empty register
func NewRegister() *Register {
	return &Register{}
}

// Register sets the register.
func (r *Register) Register(s string) {
	r.value = s
}

// Foo returns the registered value with extra foo or errors if r is empty.
func (r *Register) Foo() (string, error) {
	if r.value == "" {
		return "", ErrNothingRegistered
	}
	return "foo" + r.value, nil
}
