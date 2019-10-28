package subtest

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"regexp"
)

const (
	prefixNotDeepEqual  = "values are deep equal"
	prefixDeepEqual     = "values are not deep equal"
	prefixNotReflectNil = "value is typed or untyped nil"
	prefixReflectNil    = "value is neither typed nor untyped nil"
	prefixNotRegexpType = "value type is not applicable to regular expression matching"
	prefixMatchRegexp   = "value does not match regular expression"
	prefixNotErrorType  = "value is not an error type"
	prefixNoError       = "error value is not nil"
	prefixError         = "error value is nil"
	prefixErrorIsNot    = "error value is matching target error"
	prefixErrorIs       = "error value is not matching target error"
)

// CheckFunc is a function that return an error on failure.
type CheckFunc func(got interface{}) error

// Check runs the check function against a value function.
func (f CheckFunc) Check(vf ValueFunc) error {
	if vf == nil {
		return FailGot("missing value function", vf)
	}
	got, err := vf()
	if err != nil {
		return FailGot("value function returns an error", vf)
	}
	return f(got)
}

// Any returns a no-operation check function that never fails.
func Any() CheckFunc {
	return func(got interface{}) error { return nil }
}

// NotDeepEqual returns a check function that fails when the test value deep
// equals to reject.
func NotDeepEqual(reject interface{}) CheckFunc {
	return func(got interface{}) error {
		if reflect.DeepEqual(reject, got) {
			return FailReject(prefixNotDeepEqual, got, reject)
		}
		return nil
	}
}

// DeepEqual returns a check function that fails when the test value does not
// deep equals to expect.
func DeepEqual(expect interface{}) CheckFunc {
	return func(got interface{}) error {
		if !reflect.DeepEqual(expect, got) {
			return FailExpect(prefixDeepEqual, got, expect)
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
			return FailGot(prefixNotReflectNil, got)
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
			return FailGot(prefixReflectNil, got)
		}

		return nil
	}
}

// MatchRegexp returns a check function that fails if the test value does not
// match r. Allowed test value types are string, []byte, json.RawMessage and
// io.RuneReader.
func MatchRegexp(r *regexp.Regexp) CheckFunc {
	return func(got interface{}) error {
		var match bool
		switch gt := got.(type) {
		case string:
			match = r.MatchString(gt)
		case []byte:
			match = r.Match(gt)
		case json.RawMessage:
			match = r.Match([]byte(gt))
		case io.RuneReader:
			match = r.MatchReader(gt)
		default:
			return FailGot(prefixNotRegexpType, got)
		}
		if !match {
			return FailExpect(prefixMatchRegexp, got, r)
		}
		return nil
	}
}

// MatchRegexpPattern is a short-hand for
// MatchRegexp(regexp.MustCompile(pattern)).
func MatchRegexpPattern(pattern string) CheckFunc {
	return MatchRegexp(regexp.MustCompile(pattern))
}

// NoError returns a check function that fails when the test value is a non-nill
// error, or if it's not an error type.
func NoError() CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil { // nil never converts to an error interface.
			return FailGot(prefixNotErrorType, got)
		}
		if ok && err != nil {
			return FailGot(prefixNoError, err)
		}
		return nil
	}
}

// Error returns a check function that fails if the test value is nil or not an
// error type.
func Error() CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil { // nil never converts to an error interface.
			return FailGot(prefixNotErrorType, got)
		}
		if err == nil {
			return FailGot(prefixError, err)
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
			return FailGot(prefixNotErrorType, got)
		}
		if errors.Is(err, target) {
			return FailReject(prefixErrorIsNot, err, target)
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
			return FailGot(prefixNotErrorType, got)
		}
		if !errors.Is(err, target) {
			return FailExpect(prefixErrorIs, err, target)
		}
		return nil
	}
}
