// Package subtest provides a way of intializing small test functions suitable
// for use as with the (*testing.T).Run method. Tests using package-defined
// check functions can genrally be initalized in two ways. here is an example
// using DeepEquals:
//
//     subtest.Check(got).DeepEquals(expect) // test initialized through a C instance
//     subtest.DeepEquals(expect).Test(got)  // equivalent test initialization without a C instance
//
// In a simplar fashion, test initialization with custom check functions can be
// done through a C instance:
//
//    subtest.Check(got).Test(func(got interface{}) error {
//        if got != expect {
//            return Failure{
//                Prefix: "not plain equal",
//                Got: got,
//                Expect: expect,
//            }
//        }
//    })
//
// Or without:
//
//    subtest.Test(func() error {
//        if got != expect {
//            return Failure{
//                Prefix: "not plain equal",
//                Got: got,
//                Expect: expect,
//            }
//        }
//    })
//
package subtest

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Test returns a test that fails fatally with the error f returned by f.
func Test(f func() error) func(t *testing.T) {
	return func(t *testing.T) {
		if err := f(); err != nil {
			t.Fatal(err)
		}
	}
}

// C is a function returning a value. The main purpose of a C instance is to
// initialize tests against the result.
type C func() interface{}

// Check returns a new C for a static result value v.
func Check(v interface{}) C {
	return func() interface{} {
		return v
	}
}

// CheckJSON returns a C for the JSON unmarshaled value of b. The JSON will be
// decoded before each test, which means the result is safe to modify.
func CheckJSON(b []byte) C {
	return func() interface{} {
		var v interface{}
		err := json.Unmarshal(b, v)
		if err != nil {
			return fmt.Errorf("not valid JSON: %w", err)
		}
		return v
	}
}

// CheckJSONString is short-hand for CheckJSON([]byte(s)).
func CheckJSONString(s string) C {
	return CheckJSON([]byte(s))
}

// Value returns the value of C.
func (c C) Value() interface{} {
	return c()
}

// Test returns a test function that fails fatally with the error returned by f.
func (c C) Test(f CheckFunc) func(t *testing.T) {
	return f.Test(c())
}

// NotDeepEqual is equivalent to NotDeepEqual(v).Test(cv) where cv is the
// value wrapped by c.
func (c C) NotDeepEqual(v interface{}) func(t *testing.T) {
	return NotDeepEqual(v).Test(c())
}

// DeepEqual is equivalent to DeepEqual(v).Test(cv) where cv is the value
// wrapped by c.
func (c C) DeepEqual(v interface{}) func(t *testing.T) {
	return DeepEqual(v).Test(c())
}

// NotReflectNil is equivalent to NotReflectNil(v).Test(cv) where cv is the
// value wrapped by c.
func (c C) NotReflectNil() func(t *testing.T) {
	return NotReflectNil().Test(c())
}

// ReflectNil returns a test function that fails if the value of c is neither an
// untyped nil value nor a pointer type holding a nil value.
func (c C) ReflectNil() func(t *testing.T) {
	return ReflectNil().Test(c())
}

// NoError is equivalent to NoError(v).Test(cv) where cv is the value wrapped by
// c.
func (c C) NoError() func(t *testing.T) {
	return NoError().Test(c())
}

// Error is equivalent to Error(v).Test(cv) where cv is the value wrapped by c.
func (c C) Error() func(t *testing.T) {
	return Error().Test(c())
}

// ErrorIsNot is equivalent to ErrorIsNot(v).Test(cv) where cv is the value
// wrapped by c.
func (c C) ErrorIsNot(target error) func(t *testing.T) {
	return ErrorIsNot(target).Test(c())
}

// ErrorIs is equivalent to ErrorIs(v).Test(cv) where cv is the value wrapped by
// c.
func (c C) ErrorIs(target error) func(t *testing.T) {
	return ErrorIs(target).Test(c())
}
