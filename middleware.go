package subtest

import "reflect"

// OnFloat64 returns a check function where the test value is converted to
// float64 before it's passed to cf. Accepted input types are all uint, int and
// float variants.
func OnFloat64(cf CheckFunc) CheckFunc {
	return func(got interface{}) error {
		f, ok := asFloat64(got)
		if !ok {
			return FailGot("not representable as float64", got)
		}
		return cf(f)
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
	}
	return f, ok
}

// OnLen returns a check function where the length of the test value is
// extracted and passed to cf. Accepted input types are arrays, slices, maps,
// channels and strings.
func OnLen(cf CheckFunc) CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)
		switch rv.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			return cf(rv.Len())
		default:
			return FailGot("can not take length of type", got)
		}
	}
}

// OnCap returns a check function where the capacity of the test value is
// extracted and passed to cf. Accepted input types are arrays, slices and
// channels.
func OnCap(cf CheckFunc) CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)
		switch rv.Kind() {
		case reflect.Array, reflect.Chan, reflect.Slice:
			return cf(rv.Cap())
		default:
			return FailGot("can not take capacity of type", got)
		}
	}
}
