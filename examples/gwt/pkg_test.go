package gwt_test

import (
	"testing"

	"github.com/searis/subtest"
	pkg "github.com/searis/subtest/examples/gwt"
)

func TestFoo(t *testing.T) {
	t.Run("Given nothing is registered", func(t *testing.T) {
		reg := pkg.NewRegister()

		t.Run("When calling reg.Foo", func(t *testing.T) {
			result, err := reg.Foo()
			t.Run("Then it should error",
				subtest.Value(err).ErrorIs(pkg.ErrNothingRegistered),
			)
			t.Run("Then the result should hold a zero-value",
				subtest.Value(result).DeepEqual(""),
			)
		})
	})

	t.Run("Given bar is registered", func(t *testing.T) {
		reg := pkg.NewRegister()
		reg.Register("bar")

		t.Run("When calling pkg.Foo", func(t *testing.T) {
			bar, err := reg.Foo()

			// Abort further sub-tests on failure by relying on the return status from t.Run.
			if !t.Run("Then the result must not fail", subtest.Value(err).NoError()) {
				t.FailNow()
			}
			t.Run("Then the result should be as expected",
				subtest.Value(bar).DeepEqual("foobar"),
			)
		})
	})
}
