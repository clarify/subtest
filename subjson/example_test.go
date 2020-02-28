package subjson_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/searis/subtest"
	"github.com/searis/subtest/internal/testmock"
	"github.com/searis/subtest/subjson"
)

// t is used in example tests to mimic the `t *testing.T` parameter in test
// functions.
var t = testmock.T{
	Name: "ParentTest",
}

// TestMain override the default test runner to enforce consistent verbose
// settings. This is needed because example tests compare test output. The
// override does not affect the final test output.
func TestMain(m *testing.M) {
	testmock.VerboseMainTest(m)
}

func ExampleString() {
	const v = `"foo"`

	t.Run("v match cf", subjson.String(v).DeepEqual("foo"))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnString() {
	const v = `"foo"`
	cf := subtest.DeepEqual("foo")

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnString(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleTime() {
	const v = `"1985-12-19T18:15:00.0+01:00"`
	expect := time.Date(1985, 12, 19, 17, 15, 0, 0, time.UTC)
	t.Run("v match cf", subjson.Time(v).TimeEqual(expect))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnTime() {
	const v = `"1985-12-19T18:15:00.0+01:00"`
	cf := subtest.TimeEqual(time.Date(1985, 12, 19, 17, 15, 0, 0, time.UTC))

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnTime(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleFloat64() {
	const v = `42.5`

	t.Run("v match cf", subjson.Float64(v).NumericEqual(42.5))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnFloat64() {
	const v = `42.5`
	cf := subtest.DeepEqual(float64(42.5))

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnFloat64(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleInt64() {
	const v = `42`

	t.Run("v match cf", subjson.Int64(v).NumericEqual(42))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnInt64() {
	const v = `42`
	cf := subtest.DeepEqual(int64(42))

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnInt64(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleNumber() {
	const v = `42.5`

	t.Run("v match cf", subjson.Number(v).NumericEqual(42.5))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnNumber_float() {
	const v = `42.5`
	cf := subtest.DeepEqual(json.Number("42.5"))

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnNumber(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnNumber_int() {
	const v = `42`
	cf := subtest.DeepEqual(json.Number("42"))

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnNumber(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnNumber_string() {
	const v = `"42.5"`
	cf := subtest.DeepEqual(json.Number("42.5"))

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnNumber(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleSlice() {
	const v = `["foo", 42.5]`

	t.Run("v match cf", subjson.Slice(v).DeepEqual([]json.RawMessage{
		json.RawMessage(`"foo"`),
		json.RawMessage(`42.5`),
	}))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}
func ExampleOnSlice() {
	const v = `["foo", 42.5]`
	cf := subtest.DeepEqual([]json.RawMessage{
		json.RawMessage(`"foo"`),
		json.RawMessage(`42.5`),
	})

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnSlice(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleMap() {
	const v = `{"foo": 42.5}`

	t.Run("v match cf", subjson.Map(v).DeepEqual(map[string]json.RawMessage{
		"foo": json.RawMessage(`42.5`),
	}))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnMap() {
	const v = `{"foo": 42.5}`
	cf := subtest.DeepEqual(map[string]json.RawMessage{
		"foo": json.RawMessage(`42.5`),
	})

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnMap(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleInterface_map() {
	const v = `{"foo": 42.5}`

	t.Run("v match cf", subjson.Interface(v).DeepEqual(map[string]interface{}{
		"foo": 42.5,
	}))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleInterface_slice() {
	const v = `["foo", 42.5]`

	t.Run("v match cf", subjson.Interface(v).DeepEqual([]interface{}{
		"foo",
		42.5,
	}))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnInterface_map() {
	const v = `{"foo": 42.5}`
	cf := subtest.DeepEqual(map[string]interface{}{
		"foo": 42.5,
	})

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnInterface(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}

func ExampleOnInterface_slice() {
	const v = `["foo", 42.5]`
	cf := subtest.DeepEqual([]interface{}{
		"foo",
		42.5,
	})

	t.Run("v match cf", subtest.Value(v).Test(
		subjson.OnInterface(cf),
	))
	// Output:
	// === RUN   ParentTest/v_match_cf
	// --- PASS: ParentTest/v_match_cf (0.00s)
}
