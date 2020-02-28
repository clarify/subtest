package subjson

import (
	"fmt"

	"github.com/searis/subtest"
)

// OnString returns a check function where the test value is decoded into a
// string.
func OnString(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(String(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded string: %w", err)
		}
		return nil
	}
}

// OnNumber returns a check function where the test value is decoded into a
// json.Number before it's passed to cf.
func OnNumber(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Number(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded number: %w", err)
		}
		return nil
	}
}

// OnInt64 returns a check function where the test value is decoded into a an
// int64 before it's passed to cf.
func OnInt64(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Int64(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded int64: %w", err)
		}
		return nil
	}
}

// OnFloat64 returns a check function where the test value is decoded into a a
// float64 before it's passed to cf.
func OnFloat64(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Float64(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded float64: %w", err)
		}
		return nil
	}
}

// OnSlice returns a check function where the test value is decoded into a
// []json.RawMessage before it's passed to cf.
func OnSlice(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Slice(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded slice: %w", err)
		}
		return nil
	}
}

// OnMap returns a check function where the test value is decoded into a
// map[string]json.RawMessage before it's passed to cf.
func OnMap(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Map(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded map: %w", err)
		}
		return nil
	}
}

// OnTime returns a check function where the test value is decoded into a
// time.Time value.
func OnTime(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Time(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded time: %w", err)
		}
		return nil
	}
}

// OnInterface returns a check function where the test value is decoded into a
// interface{} value.
func OnInterface(c subtest.Check) subtest.CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Interface(got))
		if err != nil {
			return fmt.Errorf("on JSON decoded value: %w", err)
		}
		return nil
	}
}
