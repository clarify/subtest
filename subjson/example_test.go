package subjson_test

import (
	"encoding/json"
	"testing"

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
