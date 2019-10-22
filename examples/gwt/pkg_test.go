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
				subtest.Value(err).ErrorIs(pkg.ErrNothingRegistered),
			)
			t.Run("then the result should hold a zero-value",
				subtest.Value(result).DeepEqual(""),
			)
		})
	})

	t.Run("given bar is registered", func(t *testing.T) {
		reg := pkg.NewRegister()
		reg.Register("bar")

		t.Run("when calling pkg.Foo", func(t *testing.T) {
			bar, err := reg.Foo()
			if !t.Run("then the result must not fail", subtest.Value(err).NoError()) {
				t.FailNow()
			}
			t.Run("then the result should be as expected",
				subtest.Value(bar).DeepEqual("foobar"),
			)
		})
	})
}
