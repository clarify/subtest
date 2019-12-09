package subjson

import (
	"encoding/json"

	"github.com/searis/subtest"
)

// OnString returns a check function where the test value is decoded into a
// string.
func OnString(cf subtest.CheckFunc) subtest.CheckFunc {
	var v string
	return func(got interface{}) error {
		if err := unmarshalJSON(got, &v); err != nil {
			return err
		}
		return cf(v)
	}
}

// OnNumber returns a check function where the test value is decoded into a
// json.Number before it's passed to cf.
func OnNumber(cf subtest.CheckFunc) subtest.CheckFunc {
	var v json.Number
	return func(got interface{}) error {
		if err := unmarshalJSON(got, &v); err != nil {
			return err
		}
		return cf(v)
	}
}

// OnInt64 returns a check function where the test value is decoded into a an
// int64 before it's passed to cf.
func OnInt64(cf subtest.CheckFunc) subtest.CheckFunc {
	var v int64
	return func(got interface{}) error {
		if err := unmarshalJSON(got, &v); err != nil {
			return err
		}
		return cf(v)
	}
}

// OnFloat64 returns a check function where the test value is decoded into a a
// float64 before it's passed to cf.
func OnFloat64(cf subtest.CheckFunc) subtest.CheckFunc {
	var v float64
	return func(got interface{}) error {
		if err := unmarshalJSON(got, &v); err != nil {
			return err
		}
		return cf(v)
	}
}

// OnSlice returns a check function where the test value is decoded into a
// []json.RawMessage before it's passed to cf.
func OnSlice(cf subtest.CheckFunc) subtest.CheckFunc {
	var v []json.RawMessage
	return func(got interface{}) error {
		if err := unmarshalJSON(got, &v); err != nil {
			return err
		}
		return cf(v)
	}
}

// OnMap returns a check function where the test value is decoded into a
// map[string]json.RawMessage before it's passed to cf.
func OnMap(cf subtest.CheckFunc) subtest.CheckFunc {
	var v map[string]json.RawMessage
	return func(got interface{}) error {
		if err := unmarshalJSON(got, &v); err != nil {
			return err
		}
		return cf(v)
	}
}

// OnInterface returns a check function where the test value is decoded into v
// before it's passed to cf. Will panic if v is not a pointer value.
func OnInterface(cf subtest.CheckFunc, v interface{}) subtest.CheckFunc {
	return func(got interface{}) error {
		if err := unmarshalJSON(got, v); err != nil {
			return err
		}
		return cf(v)
	}
}

func unmarshalJSON(got interface{}, target interface{}) error {
	var err error

	switch gt := got.(type) {
	case []byte:
		err = json.Unmarshal(gt, &target)
	case json.RawMessage:
		err = json.Unmarshal(gt, &target)
	case string:
		err = json.Unmarshal([]byte(gt), &target)
	default:
		return subtest.FailGot("type is not JSON decodable", got)
	}

	if err != nil {
		return subtest.FailGot(err.Error(), got)
	}
	return nil
}
