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
			subtest.FailGot("value not convertable to float64", string("invalid")),
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
		name := fmt.Sprintf("When calling asFloat64 on %T", v)
		t.Run(name, func(t *testing.T) {
			vf := subtest.Float64(v)
			t.Run("Then we should get the correct float64 value", vf.DeepEqual(float64(42)))

		})
	}
}
