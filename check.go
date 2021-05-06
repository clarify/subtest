package subtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"time"
)

// Check describes the interface for a check.
type Check interface {
	Check(vf ValueFunc) error
}

// CheckFunc is a function that return an error on failure.
type CheckFunc func(got interface{}) error

// Check runs the check function against a value function.
func (f CheckFunc) Check(vf ValueFunc) error {
	return check(vf, f)
}

func check(vf ValueFunc, cf CheckFunc) error {
	if vf == nil {
		return FailGot("missing value function", vf)
	}
	got, err := vf()
	if err != nil {
		return fmt.Errorf("value function: %w", err)
	}
	return cf(got)
}

// Any returns a no-operation check function that never fails.
func Any() CheckFunc {
	return func(got interface{}) error { return nil }
}

// AllOff is a Check type that fails if any of it's members fails.
type AllOff []Check

// Check runs all member checks and returns an aggregated error of at least one
// check fails.
func (cs AllOff) Check(vf ValueFunc) error {
	var errs Errors

	for _, c := range cs {
		err := c.Check(vf)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// LessThan returns a check function that fails when the test value is not a
// numeric value less than expect.
func LessThan(expect float64) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot(msgNotFloat64, got)
		}

		if !(f < expect) {
			msg := fmt.Sprintf("%s %f", msgLessThan, expect)
			return FailGot(msg, got)
		}

		return nil
	}
}

// LessThanOrEqual returns a check function that fails when the test value is
// not a numeric value less than or equal to expect.
func LessThanOrEqual(expect float64) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot(msgNotFloat64, got)
		}

		if !(f <= expect) {
			msg := fmt.Sprintf("%s %f", msgLessThanOrEqual, expect)
			return FailGot(msg, got)
		}

		return nil
	}
}

// GreaterThan returns a check function that fails when the test value is not a
// numeric value greater than expect.
func GreaterThan(expect float64) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot(msgNotFloat64, got)
		}

		if !(f > expect) {
			msg := fmt.Sprintf("%s %f", msgGreaterThan, expect)
			return FailGot(msg, got)
		}

		return nil
	}
}

// GreaterThanOrEqual returns a check function that fails when the test value is
// not a numeric value greater than or equal to expect.
func GreaterThanOrEqual(expect float64) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot(msgNotFloat64, got)
		}

		if !(f >= expect) {
			msg := fmt.Sprintf("%s %f", msgGreaterThanOrEqual, expect)
			return FailGot(msg, got)
		}

		return nil
	}
}

// NotNumericEqual returns a check function that fails when the test value is
// a numeric value equal to expect.
func NotNumericEqual(expect float64) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot(msgNotFloat64, got)
		}

		if f == expect {
			return FailReject(msgNotNumericEqual, got, expect)
		}

		return nil
	}
}

// NumericEqual returns a check function that fails when the test value is
// not a numeric value equal to expect.
func NumericEqual(expect float64) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot(msgNotFloat64, got)
		}

		if f != expect {
			return FailExpect(msgNumericEqual, got, expect)
		}

		return nil
	}
}

// NotBefore returns a check function that fails when the test value is before
// expect. Accepts type time.Time and *time.Time.
func NotBefore(expect time.Time) CheckFunc {
	return func(got interface{}) error {
		t, ok := asTime(got)
		if !ok {
			return FailGot(msgNotTimeType, got)
		}
		if t.Before(expect) {
			msg := fmt.Sprintf("%s %v", msgNotBefore, expect)
			return FailGot(msg, got)
		}
		return nil
	}
}

// Before returns a check function that fails when the test value is not before
// expect. Accepts type time.Time and *time.Time.
func Before(expect time.Time) CheckFunc {
	return func(got interface{}) error {
		t, ok := asTime(got)
		if !ok {
			return FailGot(msgNotTimeType, got)
		}
		if !t.Before(expect) {
			msg := fmt.Sprintf("%s %v", msgBefore, expect)
			return FailGot(msg, got)
		}
		return nil
	}
}

// NotTimeEqual returns a check function that fails when the test value is a
// time semantically equal to expect. Accepts type time.Time and *time.Time.
func NotTimeEqual(expect time.Time) CheckFunc {
	return func(got interface{}) error {
		t, ok := asTime(got)
		if !ok {
			return FailGot(msgNotTimeType, got)
		}
		if t.Equal(expect) {
			return FailReject(msgNotTimeEqual, got, expect)
		}
		return nil
	}
}

// TimeEqual returns a check function that fails when the test value is not a
// time semantically equal to expect. Accepts type time.Time and *time.Time.
func TimeEqual(expect time.Time) CheckFunc {
	return func(got interface{}) error {
		t, ok := asTime(got)
		if !ok {
			return FailGot(msgNotTimeType, got)
		}
		if !t.Equal(expect) {
			return FailExpect(msgTimeEqual, got, expect)
		}
		return nil
	}
}

func asTime(got interface{}) (time.Time, bool) {
	switch gt := got.(type) {
	case time.Time:
		return gt, true
	case *time.Time:
		if gt != nil {
			return *gt, true
		}
	}
	return time.Time{}, false
}

// NotDeepEqual returns a check function that fails when the test value deep
// equals to reject.
func NotDeepEqual(reject interface{}) CheckFunc {
	return func(got interface{}) error {
		if reflect.DeepEqual(reject, got) {
			return FailReject(msgNotDeepEqual, got, reject)
		}
		return nil
	}
}

// DeepEqual returns a check function that fails when the test value does not
// deep equals to expect.
func DeepEqual(expect interface{}) CheckFunc {
	return func(got interface{}) error {
		if !reflect.DeepEqual(expect, got) {
			return FailExpect(msgDeepEqual, got, expect)
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
			return FailGot(msgNotReflectNil, got)
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
			return FailGot(msgReflectNil, got)
		}

		return nil
	}
}

// MatchRegexp returns a check function that fails if the test value does not
// match r. Allowed test value types are string, []byte, json.RawMessage,
// io.RuneReader and error.
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
		case error:
			match = gt != nil && r.MatchString(gt.Error())
		default:
			return FailGot(msgNotRegexpType, got)
		}
		if !match {
			return FailExpect(msgMatchRegexp, got, r)
		}
		return nil
	}
}

// MatchPattern is a short-hand for MatchRegexp(regexp.MustCompile(pattern)).
func MatchPattern(pattern string) CheckFunc {
	return MatchRegexp(regexp.MustCompile(pattern))
}

// NoError returns a check function that fails when the test value is a non-nill
// error, or if it's not an error type.
func NoError() CheckFunc {
	return func(got interface{}) error {
		err, ok := got.(error)
		if !ok && got != nil { // nil never converts to an error interface.
			return FailGot(msgNotErrorType, got)
		}
		if ok && err != nil {
			return FailGot(msgNoError, err)
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
			return FailGot(msgNotErrorType, got)
		}
		if err == nil {
			return FailGot(msgError, err)
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
			return FailGot(msgNotErrorType, got)
		}
		if errors.Is(err, target) {
			return FailReject(msgErrorIsNot, err, target)
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
			return FailGot(msgNotErrorType, got)
		}
		if !errors.Is(err, target) {
			return FailExpect(msgErrorIs, err, target)
		}
		return nil
	}
}

// ContainsMatch returns a check function that fails if the test value does not
// contain the check. Accepts input of type array and slice.
func ContainsMatch(c Check) CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)
		switch rv.Kind() {
		case reflect.Array, reflect.Slice:
		default:
			return FailGot(msgNotSliceArrType, got)
		}

		var subfail Failure
		for j := 0; j < rv.Len(); j++ {
			err := c.Check(Index(got, j))
			switch {
			case err == nil:
				return nil
			case errors.As(err, &subfail):
			}
		}
		fail := FailGot(msgContainsMatch, got)
		fail.Expect = subfail.Expect // may be empty
		return fail
	}
}

// Contains returns a check function that fails if the test value does not
// contain the input. Accepts input of type array and slice.
func Contains(v interface{}) CheckFunc {
	return ContainsMatch(DeepEqual(v))
}
