package subtest_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/searis/subtest"
)

func TestFloat64(t *testing.T) {
	t.Run("When resolving Float64 from invalid string", func(t *testing.T) {
		vf := subtest.Float64(string("invalid"))
		v, err := vf()
		t.Run("Then the value should be nil", subtest.Value(v).DeepEqual(nil))
		t.Run("Then we should get the correct error", subtest.Value(err).ErrorIs(
			subtest.FailGot("not convertable to float64", string("invalid")),
		))
	})

	var validInputs = []interface{}{
		int(42),
		uint(42),
		int16(42),
		int32(42),
		float64(42),
		string("42"),
		json.Number("42"),
	}
	for _, v := range validInputs {
		name := fmt.Sprintf("When resolving Float64 from %T", v)
		t.Run(name, func(t *testing.T) {
			vf := subtest.Float64(v)
			t.Run("Then we should get the correct float64 value", vf.DeepEqual(float64(42)))

		})
	}
}

func TestIndex(t *testing.T) {
	t.Run("When resolving Index from an incompatible type", func(t *testing.T) {
		vf := subtest.Index(map[int]string{0: "1"}, 0)
		v, err := vf()
		t.Run("Then it should fail with type does not support index operation",
			subtest.Value(err).MatchPattern("^type does not support index operation"),
		)
		t.Run("Then the value should be nil", subtest.Value(v).DeepEqual(nil))
	})
	t.Run("When resolving Index 0 from empty slice", func(t *testing.T) {
		vf := subtest.Index([]string{}, 0)
		v, err := vf()
		t.Run("Then it should fail with index out of range",
			subtest.Value(err).MatchPattern("^index out of range$"),
		)
		t.Run("Then the value should be nil", subtest.Value(v).DeepEqual(nil))
	})

	type testData struct {
		in    interface{}
		index int
		out   interface{}
	}
	var validInputs = []testData{
		{in: []int{42}, index: 0, out: 42},
		{in: []float64{42.0, 32.5}, index: 1, out: 32.5},
		{in: "42", index: 0, out: uint8('4')},
		{in: "42", index: 1, out: uint8('2')},
		{in: []interface{}{nil}},
	}
	for _, td := range validInputs {
		name := fmt.Sprintf("When resolving Index 0 from %#v", td.in)
		t.Run(name, func(t *testing.T) {
			vf := subtest.Index(td.in, td.index)
			t.Run("Then we should get the expected value", vf.DeepEqual(td.out))
		})
	}
}
