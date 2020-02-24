package subtest

import "fmt"

// OnFloat64 returns a check function where the test value is converted to
// float64 before it's passed to c.
func OnFloat64(c Check) CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Float64(got))
		if err != nil {
			return fmt.Errorf("on float64: %w", err)
		}
		return nil
	}
}

// OnLen returns a check function where the length of the test value is
// extracted and passed to c. Accepted input types are arrays, slices, maps,
// channels and strings.
func OnLen(c Check) CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Len(got))
		if err != nil {
			return fmt.Errorf("on len: %w", err)
		}
		return nil
	}
}

// OnCap returns a check function where the capacity of the test value is
// extracted and passed to c. Accepted input types are arrays, slices and
// channels.
func OnCap(c Check) CheckFunc {
	return func(got interface{}) error {
		err := c.Check(Cap(got))
		if err != nil {
			return fmt.Errorf("on cap: %w", err)
		}
		return nil
	}
}

// OnIndex returns a check function where the item at index i of the test value
// is passed on to c. Accepted input types are arrays, slices and strings.
func OnIndex(i int, c Check) CheckFunc {
	return func(got interface{}) error {
		vf := Index(got, i)
		err := c.Check(vf)
		if err != nil {
			return fmt.Errorf("on index %d: %w", i, err)
		}
		return nil
	}
}
