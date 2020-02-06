package subjson

import (
	"encoding/json"

	"github.com/searis/subtest"
)

// Interface returns a ValueFunc that decodes v into an interface{} value.
func Interface(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t interface{}
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// String returns a ValueFunc that decodes v into an string value.
func String(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t string
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// Int64 returns a ValueFunc that decodes v into an float64 value.
func Int64(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t int64
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// Float64 returns a ValueFunc that decodes v into an float64 value.
func Float64(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t float64
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// Number returns a ValueFunc that decodes v into an json.Number value.
func Number(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t json.Number
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// Slice returns a ValueFunc that decodes v into a []json.RawMessage value.
func Slice(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t []json.RawMessage
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// Map returns a ValueFunc that decodes v into a map[string]json.RawMessage value.
func Map(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		var t map[string]json.RawMessage
		err := unmarshalJSON(v, &t)
		return t, err
	}
}

// Len returns a ValueFunc that returns the length of the decoded value.
func Len(v interface{}) subtest.ValueFunc {
	return func() (interface{}, error) {
		v, err := Interface(v)()
		if err != nil {
			return v, err
		}
		return subtest.Len(v)()
	}
}

func unmarshalJSON(got interface{}, target interface{}) error {
	var err error

	switch gt := got.(type) {
	case []byte:
		err = json.Unmarshal(gt, target)
	case json.RawMessage:
		err = json.Unmarshal(gt, target)
	case string:
		err = json.Unmarshal([]byte(gt), target)
	default:
		err = subtest.FailGot("type is not JSON decodable", got)
	}

	if err != nil {
		return subtest.FailGot(err.Error(), got)
	}
	return nil
}
