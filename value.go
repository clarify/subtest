package subtest

import (
	"regexp"
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

// Test returns a test function that fails fatally with the error returned by
// f.Check(vf).
func (vf ValueFunc) Test(f CheckFunc) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if err := f.Check(vf); err != nil {
			t.Fatal(err)
		}
	}
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

// ValidateMap is equivalent to vf.Test(s.ValidateMap()).
func (vf ValueFunc) ValidateMap(s Schema) func(t *testing.T) {
	return vf.Test(s.ValidateMap())
}
