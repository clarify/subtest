package subjson

import (
	"encoding/json"

	"github.com/clarify/subtest"
)

// Fields is a short-hand for OnMap(subtest.Fields{...})
type Fields subtest.Fields

// Check validates the json decoded vf against m, expecting vf to return a json
// map.
func (m Fields) Check(vf subtest.ValueFunc) error {
	return OnMap(subtest.Fields(m)).Check(vf)
}

// RawEqual is a short-hand for subtest.DeepEqual(json.RawMessage(s))
type RawEqual []byte

// Check validates the raw json equals m.
func (m RawEqual) Check(vf subtest.ValueFunc) error {
	return subtest.DeepEqual(json.RawMessage(m)).Check(vf)
}

// IterateSlice is a short-hand for OnSlice(subtest.Iterate(cs...))
func IterateSlice(cs ...subtest.Check) subtest.Check {
	return OnSlice(subtest.Iterate(cs...))
}

// LessThan is a short-hand for OnNumber(subtest.LessThan(expect)).
func LessThan(expect float64) subtest.CheckFunc {
	return OnNumber(subtest.LessThan(expect))
}

// LessThanOrEqual is a short-hand for OnNumber(subtest.LessThanOrEqual(expect)).
func LessThanOrEqual(expect float64) subtest.CheckFunc {
	return OnNumber(subtest.LessThanOrEqual(expect))
}

// GreaterThan is a short-hand for OnNumber(subtest.GreaterThan(expect)).
func GreaterThan(expect float64) subtest.CheckFunc {
	return OnNumber(subtest.GreaterThan(expect))
}

// GreaterThanOrEqual is a short-hand for OnNumber(subtest.GreaterThanOrEqual(expect)).
func GreaterThanOrEqual(expect float64) subtest.CheckFunc {
	return OnNumber(subtest.GreaterThanOrEqual(expect))
}

// NumericEqual is a short-hand for OnNumber(subtest.NumericEqual(expect)).
func NumericEqual(expect float64) subtest.CheckFunc {
	return OnNumber(subtest.NumericEqual(expect))
}

// NotDecodesTo is a short-hand for OnInterface(subtest.NotDeepEqual(reject)).
func NotDecodesTo(reject interface{}) subtest.CheckFunc {
	return OnInterface(subtest.NotDeepEqual(reject))
}

// DecodesTo is a short-hand for OnInterface(subtest.DeepEqual(expect)).
func DecodesTo(expect interface{}) subtest.CheckFunc {
	return OnInterface(subtest.DeepEqual(expect))
}

// NotNil is a short-hand for OnInterface(subtest.NotDeepEqual(nil)).
func NotNil() subtest.CheckFunc {
	return OnInterface(subtest.NotDeepEqual(nil))
}

// Nil is a short-hand for OnInterface(subtest.DeepEqual(nil)).
func Nil() subtest.CheckFunc {
	return OnInterface(subtest.DeepEqual(nil))
}
