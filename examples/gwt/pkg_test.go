package gwt_test

import (
	"testing"

	"github.com/searis/subtest"
	pkg "github.com/searis/subtest/examples/gwt"
)

func TestFoo(t *testing.T) {
	t.Run("given nothing is registered", func(t *testing.T) {
		reg := pkg.NewRegister()

		t.Run("when calling reg.Foo", func(t *testing.T) {
			result, err := reg.Foo()

			t.Run("then it should error",
				subtest.Check(err).ErrorIs(pkg.ErrNothingRegistered),
			)

			t.Run("then the result should hold a zero-value",
				subtest.Check(result).DeepEqual(""),
			)
		})
	})

	t.Run("given bar is registered", func(t *testing.T) {
		reg := pkg.NewRegister()
		reg.Register("bar")

		t.Run("when calling pkg.Foo", func(t *testing.T) {
			bar, err := reg.Foo()
			subtest.NoError().Test(err)(t) // short-hand to abort on failure.

			// An equivalent way to abort on failure.
			if !t.Run("then the result must not fail", subtest.Check(err).NoError()) {
				t.FailNow()
			}

			t.Run("then the result should be as expected",
				subtest.Check(bar).DeepEqual("foobar"),
			)
		})
	})
}
