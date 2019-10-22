package subtest_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/searis/subtest"
)

func TestSchemaValidateMap(t *testing.T) {
	t.Run("given an empty schema", func(t *testing.T) {
		s := subtest.Schema{}
		cf := s.ValidateMap()

		t.Run("then additional fields should result in an error",
			subtest.Value(cf(map[string]string{"a": "", "b": "foo"})).ErrorIs(
				subtest.Failf(`got additional keys: "a", "b"`),
			),
		)
		t.Run("then a empty string should be recognized as an additional key",
			subtest.Value(cf(map[string]string{"": ""})).ErrorIs(
				subtest.Failf(`got additional keys: ""`),
			),
		)

	})
	t.Run("given a schema with required integer keys and allowed additional fields", func(t *testing.T) {
		s := subtest.Schema{
			Required:         []interface{}{1, 2, 3},
			AdditionalFields: subtest.Any(),
		}
		cf := s.ValidateMap()

		t.Run("then it should match a map[int]string value with all keys",
			subtest.Value(map[int]string{1: "", 2: "", 3: ""}).Test(cf),
		)
		t.Run("then it should match a map[int]struct{} value with all keys",
			subtest.Value(map[int]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}}).Test(cf),
		)
		t.Run("then it should match a map[interface{}]struct{} value with all keys",
			subtest.Value(map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}}).Test(cf),
		)
		t.Run("then it should match a map with additional fields",
			subtest.Value(map[int]string{1: "", 2: "", 3: "", 4: ""}).Test(cf),
		)
		t.Run("then missing fields should result in an error",
			subtest.Value(cf(map[int]string{2: ""})).ErrorIs(
				subtest.Failf("missing required keys: 1, 3"),
			),
		)
	})
}

func TestJSONSchema(t *testing.T) {
	const v = `{"foo": "bar", "bar": 42}`

	vf := subtest.ValueFunc(func() (interface{}, error) {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("not a JSON object: %w", err)
		}
		return m, nil
	})

	t.Run("v match schema", vf.ValidateMap(subtest.Schema{
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