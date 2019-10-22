package subtest

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

// Failure is an error type that aid with consistent formatting of test
// failures. In error matching, two Failure instances are considered equal when
// their formattet content is the same.
type Failure struct {
	Prefix string
	Got    string
	Expect string
	Reject string
}

// Failf formats a plain text failure.
func Failf(format string, v ...interface{}) Failure {
	return Failure{Prefix: fmt.Sprintf(format, v...)}
}

// FailExpect formats a failure for content that is not matching some expected
// value. The package type formatter is used.
func FailExpect(prefix string, got, expect interface{}) Failure {
	return Failure{
		Prefix: prefix,
		Got:    formatIndentedType(got),
		Expect: formatIndentedType(expect),
	}
}

// FailReject formats a failure for content that is matching some rejected
// value. The package type formatter is used.
func FailReject(prefix string, got, reject interface{}) Failure {
	return Failure{
		Prefix: prefix,
		Got:    formatIndentedType(got),
		Reject: formatIndentedType(reject),
	}
}

// FailGot formats a failure for some unexpected content. The package type
// formatter is used.
func FailGot(prefix string, got interface{}) Failure {
	return Failure{
		Prefix: prefix,
		Got:    formatIndentedType(got),
	}
}

func (f Failure) Error() string {
	const fmtS = "\n%s: %s"
	s := f.Prefix
	if f.Got != "" {
		s += fmt.Sprintf(fmtS, "got", f.Got)
	}
	if f.Expect != "" {
		s += fmt.Sprintf(fmtS, "want", f.Expect)
	}
	if f.Reject != "" {
		s += fmt.Sprintf(fmtS, "don't want", f.Reject)
	}
	return s
}

// PrefixError wraps a (type-formatted) error with a prefix string.
type PrefixError struct {
	Key     string
	Err     error
	Newline bool
}

// KeyError returns an error prefixed by a key.
func KeyError(key interface{}, err error) PrefixError {
	return PrefixError{
		Key: fmt.Sprintf("key %#v", key),
		Err: err,
	}
}

func (err PrefixError) Error() string {
	if err.Newline {
		return fmt.Sprintf("%s:\n%s", err.Key, err.Err)
	}
	return fmt.Sprintf("%s: %s", err.Key, err.Err)
}

// Unwrap returns the wrapped error.
func (err PrefixError) Unwrap() error {
	return err.Err
}

// Errors combine the output of multiple errors on separate lines.
type Errors []error

func (errs Errors) Error() string {
	var buf bytes.Buffer

	for i, err := range errs {
		fmt.Fprintf(&buf, "#%d: %s\n", i, err)
	}

	return buf.String()
}

// Is returns true if target is found within errs or if target deep equals
// errs.
func (errs Errors) Is(target error) bool {
	for _, err := range errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return reflect.DeepEqual(errs, target)
}
