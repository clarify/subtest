package subjson_test

import (
	"encoding/json"

	"github.com/searis/subtest"
	"github.com/searis/subtest/subjson"
)

func ExampleOnMap_failingTest() {
	const v = `{"foo":"bar", "bar":"baz"}`
	cf := subtest.DeepEqual(map[string]json.RawMessage{
		"foo": json.RawMessage(`"bar"`),
		"bar": json.RawMessage(`"foobar"`),
	})

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnMap(cf),
	))
	// FIXME: t.Helper issue causes return of value.go:112 instead of this file.
	// Could be related to https://github.com/golang/go/issues/23249.

	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- FAIL: ParentTest/v_match_cf (0.00s)
	//     value.go:112: values are not deep equal
	//         got: map[string]json.RawMessage
	//             map[bar:[34 98 97 122 34] foo:[34 98 97 114 34]]
	//         want: map[string]json.RawMessage
	//             map[bar:[34 102 111 111 98 97 114 34] foo:[34 98 97 114 34]]
}

func ExampleOnMap_failingSchemaTest() {
	const v = `{"foo":"bar", "bar":"baz"}`
	c := subtest.Schema{
		Fields: subtest.Fields{
			"foo": subjson.DecodesTo("bar"),                       // check decoded content.
			"bar": subtest.DeepEqual(json.RawMessage(`"foobar"`)), // check raw JSON.
		},
	}

	t.Run("v match cf", subtest.Value(v).Test(subjson.OnMap(c)))
	// FIXME: t.Helper issue causes return of value.go:112 instead of this file.
	// Could be related to https://github.com/golang/go/issues/23249.

	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- FAIL: ParentTest/v_match_cf (0.00s)
	//     value.go:112: value not matching schema:
	//         #0: key "bar": values are not deep equal
	//         got: json.RawMessage
	//             `"baz"`
	//         want: json.RawMessage
	//             `"foobar"`
}
func ExampleOnMap_passingTest() {
	const v = `{"foo":"bar", "bar":"baz"}`
	cf := subtest.DeepEqual(map[string]json.RawMessage{
		"foo": json.RawMessage(`"bar"`),
		"bar": json.RawMessage(`"baz"`),
	})

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnMap(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnMap_passingSchemaTest() {
	const v = `{"foo":"bar", "bar":"baz"}`
	c := subtest.Schema{
		Fields: subtest.Fields{
			"foo": subjson.DecodesTo("bar"),
			"bar": subtest.DeepEqual(json.RawMessage(`"baz"`)),
		},
	}

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnMap(c),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}
