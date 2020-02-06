package subtest

import (
	"reflect"
	"regexp"
	"strconv"
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

// ValueFunc is a function returning a value. The main purpose of a ValueFunc
// instance is to initialize tests against the the returned value.
type ValueFunc func() (interface{}, error)

// Value returns a new ValueFunc for a static value v.
func Value(v interface{}) ValueFunc {
	return func() (interface{}, error) {
		return v, nil
	}
}

// Len returns a new ValueFunc for the length of v.
func Len(v interface{}) ValueFunc {
	return func() (interface{}, error) {
		l, ok := asLen(v)
		if !ok {
			return nil, FailGot(msgNoLen, v)
		}
		return l, nil
	}
}

func asLen(v interface{}) (l int, ok bool) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l = rv.Len()
		ok = true
	}
	return
}

// Cap returns a new ValueFunc for the capacity of v.
func Cap(v interface{}) ValueFunc {
	return func() (interface{}, error) {
		l, ok := asCap(v)
		if !ok {
			return nil, FailGot(msgNoCap, v)
		}
		return l, nil
	}
}

func asCap(v interface{}) (c int, ok bool) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		c = rv.Cap()
		ok = true
	}
	return
}

// Float64 returns a new ValueFunc that parses v into a float64. Valid input
// types are any numeric kinds  or string kinds.
func Float64(v interface{}) ValueFunc {
	return func() (interface{}, error) {
		f, ok := asFloat64(v)
		if !ok {
			return nil, FailGot(msgNoFloat64, v)
		}
		return f, nil
	}
}

func asFloat64(v interface{}) (f float64, ok bool) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f = float64(rv.Uint())
		ok = true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f = float64(rv.Int())
		ok = true
	case reflect.Float32, reflect.Float64:
		f = rv.Float()
		ok = true
	case reflect.String:
		// E.g. json.Number
		var err error
		f, err = strconv.ParseFloat(rv.String(), 64)
		ok = err == nil
	}
	return f, ok
}

// Test returns a test function that fails fatally with the error returned by
// f.Check(vf).
func (vf ValueFunc) Test(c Check) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if err := c.Check(vf); err != nil {
			t.Fatal(err)
		}
	}
}

// LessThan is equivalent to vf.Test(LessThan(v)).
func (vf ValueFunc) LessThan(v float64) func(t *testing.T) {
	return vf.Test(LessThan(v))
}

// LessThanOrEqual is equivalent to vf.Test(LessThanOrEqual(v)).
func (vf ValueFunc) LessThanOrEqual(v float64) func(t *testing.T) {
	return vf.Test(LessThanOrEqual(v))
}

// GreaterThan is equivalent to vf.Test(GreaterThan(v)).
func (vf ValueFunc) GreaterThan(v float64) func(t *testing.T) {
	return vf.Test(GreaterThan(v))
}

// GreaterThanOrEqual is equivalent to vf.Test(GreaterThanOrEqual(v)).
func (vf ValueFunc) GreaterThanOrEqual(v float64) func(t *testing.T) {
	return vf.Test(GreaterThanOrEqual(v))
}

// NumericNotEqual is equivalent to vf.Test(NumericNotEqual(v)).
func (vf ValueFunc) NumericNotEqual(v float64) func(t *testing.T) {
	return vf.Test(NumericNotEqual(v))
}

// NumericEqual is equivalent to vf.Test(NumericEqual(v)).
func (vf ValueFunc) NumericEqual(v float64) func(t *testing.T) {
	return vf.Test(NumericEqual(v))
}

// NotDeepEqual is equivalent to vf.Test(NotDeepEqual(v)).
func (vf ValueFunc) NotDeepEqual(v interface{}) func(t *testing.T) {
	return vf.Test(NotDeepEqual(v))
}

// DeepEqual is equivalent to vf.Test(DeepEqual(v)).
func (vf ValueFunc) DeepEqual(v interface{}) func(t *testing.T) {
	return vf.Test(DeepEqual(v))
}

// NotReflectNil is equivalent to vf.Test(NotReflectNil(v)).
func (vf ValueFunc) NotReflectNil() func(t *testing.T) {
	return vf.Test(NotReflectNil())
}

// ReflectNil is equivalent to vf.Test(ReflectNil(v)).
func (vf ValueFunc) ReflectNil() func(t *testing.T) {
	return vf.Test(ReflectNil())
}

// MatchRegexp is equivalent to vf.Test(s.MatchRegexp(r)).
func (vf ValueFunc) MatchRegexp(r *regexp.Regexp) func(t *testing.T) {
	return vf.Test(MatchRegexp(r))
}

// MatchRegexpPattern is equivalent to vf.Test(s.MatchRegexpPattern(pattern)).
func (vf ValueFunc) MatchRegexpPattern(pattern string) func(t *testing.T) {
	return vf.Test(MatchRegexpPattern(pattern))
}

// NoError is equivalent to vf.Test(NoError(v)).
func (vf ValueFunc) NoError() func(t *testing.T) {
	return vf.Test(NoError())
}

// Error is equivalent to vf.Test(Error(v)).
func (vf ValueFunc) Error() func(t *testing.T) {
	return vf.Test(Error())
}

// ErrorIsNot is equivalent to vf.Test(ErrorIsNot(v)).
// wrapped by vf.
func (vf ValueFunc) ErrorIsNot(target error) func(t *testing.T) {
	return vf.Test(ErrorIsNot(target))
}

// ErrorIs is equivalent to vf.Test(ErrorIs(v)).
func (vf ValueFunc) ErrorIs(target error) func(t *testing.T) {
	return vf.Test(ErrorIs(target))
}

// TestField is equivalent to vf.Test(Schema{Fields: fields}).
func (vf ValueFunc) TestField(fields Fields) func(t *testing.T) {
	return vf.Test(Schema{Fields: fields})
}
