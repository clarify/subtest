package subtest_test

import (
	"testing"

	"github.com/searis/subtest"
)

func TestOnFloat64(t *testing.T) {
	t.Run("given a check OnFloat64(DeepEqual(float64(v)))", func(t *testing.T) {
		const v = 42
		cf := subtest.OnFloat64(subtest.DeepEqual(float64(v)))
		t.Run("when cheking against float64(v)", func(t *testing.T) {
			vf := subtest.Value(cf(float64(v)))
			t.Run("then it should pass", vf.NoError())
		})
		t.Run("when cheking against int16(v)", func(t *testing.T) {
			vf := subtest.Value(cf(int16(v)))
			t.Run("then it should pass", vf.NoError())
		})
		t.Run("when cheking against uint(v)", func(t *testing.T) {
			vf := subtest.Value(cf(uint(v)))
			t.Run("then it should pass", vf.NoError())
		})
		t.Run(`when cheking against string("42")`, func(t *testing.T) {
			vf := subtest.Value(cf("42"))
			expect := subtest.Failure{
				Prefix: "not representable as float64",
				Got:    "string\n\t\"42\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestOnLen(t *testing.T) {
	t.Run("given a check OnLen(DeepEqual(3))", func(t *testing.T) {
		cf := subtest.OnLen(subtest.DeepEqual(3))
		t.Run("when cheking against make([]int,3)", func(t *testing.T) {
			vf := subtest.Value(cf(make([]int, 3)))
			t.Run("then it should pass", vf.NoError())
		})
		t.Run(`when cheking against int64(42)`, func(t *testing.T) {
			vf := subtest.Value(cf(int64(42)))
			expect := subtest.FailGot("can not take length of type", int64(42))
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run(`when cheking against string("42")`, func(t *testing.T) {
			vf := subtest.Value(cf("42"))
			expect := subtest.FailExpect("values are not deep equal", 2, 3)
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestOnCap(t *testing.T) {
	t.Run("given a check OnCap(DeepEqual(3))", func(t *testing.T) {
		cf := subtest.OnCap(subtest.DeepEqual(3))
		t.Run("when cheking against make([]int,3)", func(t *testing.T) {
			vf := subtest.Value(cf(make([]int, 3)))
			t.Run("then it should pass", vf.NoError())
		})
		t.Run("when cheking against make([]int,3,4)", func(t *testing.T) {
			vf := subtest.Value(cf(make([]int, 3, 4)))
			expect := subtest.FailExpect("values are not deep equal", 4, 3)
			t.Run("then it should pass", vf.ErrorIs(expect))
		})
		t.Run(`when cheking against int64(42)`, func(t *testing.T) {
			vf := subtest.Value(cf(int64(42)))
			expect := subtest.FailGot("can not take capacity of type", int64(42))
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run(`when cheking against string("42")`, func(t *testing.T) {
			vf := subtest.Value(cf("42"))
			expect := subtest.FailGot("can not take capacity of type", "42")
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}
