package jsondiff

import (
	"encoding/json"
	"fmt"

	"github.com/clarify/subtest"
	"github.com/nsf/jsondiff"
)

// EqualJSON returns a check functions that fails if the the test value is not
// equivalent to the expected JSON, or if either value is not valid JSON.
func EqualJSON(expect string) subtest.CheckFunc {
	return func(got interface{}) error {
		d, s := compareJSON([]byte(expect), got)
		switch d {
		case jsondiff.BothArgsAreInvalidJson:
			fallthrough
		case jsondiff.SecondArgIsInvalidJson:
			return subtest.FailGot("expect value is invalid JSON", expect)
		case jsondiff.FirstArgIsInvalidJson:
			return subtest.FailGot("test value is invalid JSON", got)
		case jsondiff.FullMatch:
			return nil
		default:
			return fmt.Errorf("unequal JSON:\n %s", s)
		}
	}
}

// SupersetOfJSON returns a check functions that fails if the test value is not
// a superset of expect.
func SupersetOfJSON(expect string) subtest.CheckFunc {
	return func(got interface{}) error {
		d, s := compareJSON(got, []byte(expect))
		switch d {
		case jsondiff.BothArgsAreInvalidJson:
			fallthrough
		case jsondiff.SecondArgIsInvalidJson:
			return subtest.FailGot("expect value is invalid JSON", expect)
		case jsondiff.FirstArgIsInvalidJson:
			return subtest.FailGot("test value is invalid JSON", got)
		case jsondiff.FullMatch, jsondiff.SupersetMatch:
			return nil
		default:
			return fmt.Errorf("test value is not superset of expect:\n %s", s)
		}
	}
}

// SubsetOfJSON returns a check functions that fails if the test value is not a
// subset of expect.
func SubsetOfJSON(expect string) subtest.CheckFunc {
	return func(got interface{}) error {
		d, s := compareJSON([]byte(expect), got)
		switch d {
		case jsondiff.BothArgsAreInvalidJson:
			fallthrough
		case jsondiff.FirstArgIsInvalidJson:
			return subtest.FailGot("expect value is invalid JSON", got)
		case jsondiff.SecondArgIsInvalidJson:
			return subtest.FailGot("test value is invalid JSON", expect)
		case jsondiff.FullMatch, jsondiff.SupersetMatch:
			return nil
		default:
			return fmt.Errorf("test value is not subset of expect:\n %s", s)
		}
	}
}

func compareJSON(a, b interface{}) (jsondiff.Difference, string) {
	var ab, bb []byte

	switch at := a.(type) {
	case []byte:
		ab = at
	case json.RawMessage:
		ab = []byte(at)
	case string:
		ab = []byte(at)
	default:
		var err error
		ab, err = json.Marshal(at)
		if err != nil {
			return jsondiff.FirstArgIsInvalidJson, ""
		}
	}

	switch bt := b.(type) {
	case []byte:
		bb = bt
	case json.RawMessage:
		bb = []byte(bt)
	case string:
		bb = []byte(bt)
	default:
		var err error
		bb, err = json.Marshal(bt)
		if err != nil {
			return jsondiff.SecondArgIsInvalidJson, ""
		}
	}

	return jsondiff.Compare(ab, bb, &cfg.jsonDiff)
}
