package subtest

// OnFloat64 returns a check function where the test value is converted to
// float64 before it's passed to c.
func OnFloat64(c Check) CheckFunc {
	return func(got interface{}) error {
		return c.Check(Float64(got))
	}
}

// OnLen returns a check function where the length of the test value is
// extracted and passed to c. Accepted input types are arrays, slices, maps,
// channels and strings.
func OnLen(c Check) CheckFunc {
	return func(got interface{}) error {
		vf := Len(got)
		return c.Check(vf)
	}
}

// OnCap returns a check function where the capacity of the test value is
// extracted and passed to c. Accepted input types are arrays, slices and
// channels.
func OnCap(c Check) CheckFunc {
	return func(got interface{}) error {
		vf := Cap(got)
		return cf.Check(vf)
	}
}

// OnIndex returns a check function where the item at index i of the test value
// is passed on to c. Accepted input types are arrays, slices and strings.
func OnIndex(i int, c Check) CheckFunc {
	return func(got interface{}) error {
		vf := Index(got, i)
		return c.Check(vf)
	}
}
