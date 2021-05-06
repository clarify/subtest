package subtest_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/searis/subtest"
)

func TestSchema_map(t *testing.T) {
	t.Run("given an empty schema", func(t *testing.T) {
		c := subtest.Schema{}

		vf1 := subtest.Value(map[string]string{"a": "", "b": "foo"})
		t.Run("then additional fields should result in an error",
			subtest.Value(c.Check(vf1)).ErrorIs(
				subtest.Failf(`got additional keys: "a", "b"`),
			),
		)

		vf2 := subtest.Value(map[string]string{"": ""})
		t.Run("then an empty string should be recognized as an additional key",
			subtest.Value(c.Check(vf2)).ErrorIs(
				subtest.Failf(`got additional keys: ""`),
			),
		)

	})
	t.Run("given a schema with required integer keys and allowed additional fields", func(t *testing.T) {
		c := subtest.Schema{
			Required:         []interface{}{1, 2, 3},
			AdditionalFields: subtest.Any(),
		}

		t.Run("then it should match a map[int]string value with all keys",
			subtest.Value(map[int]string{1: "", 2: "", 3: ""}).Test(c),
		)
		t.Run("then it should match a map[int]struct{} value with all keys",
			subtest.Value(map[int]struct{}{1: {}, 2: {}, 3: {}}).Test(c),
		)
		t.Run("then it should match a map[interface{}]struct{} value with all keys",
			subtest.Value(map[interface{}]struct{}{1: {}, 2: {}, 3: {}}).Test(c),
		)
		t.Run("then it should match a map with additional fields",
			subtest.Value(map[int]string{1: "", 2: "", 3: "", 4: ""}).Test(c),
		)
		vf := subtest.Value(map[int]string{2: ""})
		t.Run("then missing fields should result in an error",
			subtest.Value(c.Check(vf)).ErrorIs(
				subtest.Failf("missing required keys: 1, 3"),
			),
		)
	})
}

func TestSchema_JSON(t *testing.T) {
	const v = `{"foo": "bar", "bar": 42}`

	vf := subtest.ValueFunc(func() (interface{}, error) {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("not a JSON object: %w", err)
		}
		return m, nil
	})

	t.Run("v match schema", vf.Test(subtest.Schema{
		Required: []interface{}{"foo", "bar"},
		Fields: subtest.Fields{
			"foo": subtest.DeepEqual("bar"),
			"bar": floatGt(41),
		},
	}))
}

func floatGt(compare float64) subtest.CheckFunc {
	return func(got interface{}) error {
		v, ok := got.(float64)
		if !ok {
			return subtest.FailGot("not a float64 value", v)
		}
		if !(v > compare) {
			return subtest.FailGot(fmt.Sprintf("value <= %v", compare), got)
		}
		return nil
	}
}
