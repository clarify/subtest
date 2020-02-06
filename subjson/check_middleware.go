package subjson

import (
	"github.com/searis/subtest"
)

// OnString returns a check function where the test value is decoded into a
// string.
func OnString(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(String(got))
	}
}

// OnNumber returns a check function where the test value is decoded into a
// json.Number before it's passed to cf.
func OnNumber(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(Number(got))
	}
}

// OnInt64 returns a check function where the test value is decoded into a an
// int64 before it's passed to cf.
func OnInt64(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(Int64(got))
	}
}

// OnFloat64 returns a check function where the test value is decoded into a a
// float64 before it's passed to cf.
func OnFloat64(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(Float64(got))
	}
}

// OnSlice returns a check function where the test value is decoded into a
// []json.RawMessage before it's passed to cf.
func OnSlice(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(Slice(got))
	}
}

// OnMap returns a check function where the test value is decoded into a
// map[string]json.RawMessage before it's passed to cf.
func OnMap(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(Map(got))
	}
}

// OnInterface returns a check function where the test value is decoded into a
// interface{} value.
func OnInterface(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		return c.Check(Interface(got))
	}
}
