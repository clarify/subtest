package subtest

import (
	"errors"
	"reflect"
	"testing"
)

const (
	prefixNotDeepEqual  = "values are deep equal"
	prefixDeepEqual     = "values are not deep equal"
	prefixNotReflectNil = "value is typed or untyped nil"
	prefixReflectNil    = "value is neither typed nor untyped nil"
	prefixNotErrorType  = "value is not an error type"
	prefixNoError       = "error value is not nil"
	prefixError         = "error value is nil"
	prefixErrorIsNot    = "error value is matching target error"
	prefixErrorIs       = "error value is not matching target error"
)

// CheckFunc is a function that return an error on failure.
type CheckFunc func(got interface{}) error

// Test returns a test function that fail fatally if f(got) returns an error.
func (f CheckFunc) Test(got interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if err := f(got); err != nil {
			t.Fatal(err)
		}
	}
}

// NotDeepEqual returns a check function that fails when the test value deep
// equals to reject.
func NotDeepEqual(reject interface{}) CheckFunc {
	return func(got interface{}) error {
		if reflect.DeepEqual(reject, got) {
			return FailureReject(prefixNotDeepEqual, got, reject)
		}
		return nil
	}
}

// DeepEqual returns a check function that fails when the test value does not
// deep equals to expect.
func DeepEqual(expect interface{}) CheckFunc {
	return func(got interface{}) error {
		if !reflect.DeepEqual(expect, got) {
			return FailureExpect(prefixDeepEqual, got, expect)
		}
		return nil
	}
}

// NotReflectNil returns a check function that fails when the test value is
// either an untyped nil value or reflects to a pointer with a nil value.
func NotReflectNil() CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)

		if got == nil || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
			return FailureGot(prefixNotReflectNil, got)
		}
		return nil
	}
}

// ReflectNil returns a check function that fails when the test value is
// neither an untyped nil value nor reflects to a pointer with a nil value.
func ReflectNil() CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)

		if got != nil && !(rv.Kind() == reflect.Ptr && rv.IsNil()) {
			return FailureGot(prefixReflectNil, got)
		}

		return nil
	}
}

// NoError returns a check function that fails when the test value is a non-nill
// error, or if it's not an error type.
func NoError() CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil {
			return FailureGot(prefixNotErrorType, got)
		}
		if ok && err != nil {
			return FailureGot(prefixNoError, err)
		}
		return nil
	}
}

// Error returns a check function that fails if the test value is not nil, or
// not an error type.
func Error() CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil { // nil don't convert to an error(nil).
			return FailureGot(prefixNotErrorType, got)
		}
		if err == nil {
			return FailureGot(prefixError, err)
		}
		return nil
	}
}

// ErrorIsNot returns a check function that fails if the test value is an error
// matching target, or not an error type.
func ErrorIsNot(target error) CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil {
			return FailureGot(prefixNotErrorType, got)
		}
		if errors.Is(err, target) {
			return FailureReject(prefixErrorIsNot, err, target)
		}
		return nil
	}
}

// ErrorIs returns a check function that fails if the test value is not an error
// matching target, or not an error type.
func ErrorIs(target error) CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil {
			return FailureGot(prefixNotErrorType, got)
		}
		if !errors.Is(err, target) {
			return FailureExpect(prefixErrorIs, err, target)
		}
		return nil
	}
}