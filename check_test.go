package subtest_test

import (
	"testing"

	"github.com/searis/subtest"
)

func TestDeepEqual(t *testing.T) {
	// Demonstrates GWT-style sub-tests.

	t.Run("when checking bool(true)", func(t *testing.T) {
		v := true
		cv := subtest.Check(v)

		t.Run("then it should match true", cv.DeepEqual(true))
		t.Run("then it should match true", subtest.DeepEqual(true).Test(v)) // equivalent

		expectErr := subtest.Failure{Prefix: "values are not deep equal", Got: "bool\n\ttrue", Expect: "bool\n\tfalse"}
		t.Run("then a test against false should fail",
			subtest.Check(subtest.DeepEqual(false)(v)).ErrorIs(expectErr),
		)

		w := new(bool)
		*w = true

		expectErr = subtest.Failure{Prefix: "values are not deep equal", Got: "bool\n\ttrue", Expect: "*bool\n\ttrue"}
		t.Run("then a test against &true should fail",
			subtest.Check(subtest.DeepEqual(w)(v)).ErrorIs(expectErr),
		)
	})

	t.Run("given a nested struct type T", func(t *testing.T) {
		type T struct {
			A string
			B map[string]string
		}

		t.Run("when checking a non-nil *T value", func(t *testing.T) {
			v := &T{A: "a", B: map[string]string{"C": "D"}}

			t.Run("then it should match an equivalent *T value",
				subtest.Check(v).DeepEqual(&T{A: "a", B: map[string]string{"C": "D"}}),
			)
			t.Run("then it should not match a different *T value",
				subtest.Check(subtest.DeepEqual(&T{A: "a", B: map[string]string{"C": "E"}})(v)).Error(),
			)
			t.Run("then it should not match an equivalent T value",
				subtest.Check(subtest.DeepEqual(T{A: "a", B: map[string]string{"C": "D"}})(v)).Error(),
			)

		})
	})

}

func TestNotDeepEqual(t *testing.T) {
	// Demonstrates GWT-style sub-tests.

	t.Run("when checking bool(true)", func(t *testing.T) {
		v := true
		cv := subtest.Check(v)

		t.Run("then it should accept false", cv.NotDeepEqual(false))
		t.Run("then it should accept false", subtest.NotDeepEqual(false).Test(v)) // equivalent

		w := new(bool)
		*w = true
		t.Run("then a test against &true should pass", cv.NotDeepEqual(w))

		expectErr := subtest.Failure{Prefix: "values are deep equal", Got: "bool\n\ttrue", Reject: "bool\n\ttrue"}
		t.Run("then a test against false should fail",
			subtest.ErrorIs(expectErr).Test(subtest.NotDeepEqual(true)(v)),
		)

	})

	t.Run("given a nested struct type T", func(t *testing.T) {
		type T struct {
			A string
			B map[string]string
		}

		t.Run("when checking a non-nil *T value", func(t *testing.T) {
			v := &T{A: "a", B: map[string]string{"C": "D"}}

			t.Run("then it should not accept an equivalent *T value",
				subtest.Error().Test(subtest.NotDeepEqual(&T{A: "a", B: map[string]string{"C": "D"}})(v)),
			)
			t.Run("then it should accept a different *T value",
				subtest.Check(v).NotDeepEqual(&T{A: "a", B: map[string]string{"C": "E"}}),
			)
			t.Run("then it should not match an equivalent T value",
				subtest.Check(v).NotDeepEqual(T{A: "a", B: map[string]string{"C": "D"}}),
			)

		})
	})
}

func TestCheckReflectNil(t *testing.T) {
	type T struct{ Foo string }

	f := subtest.ReflectNil()

	t.Run("nil", f.Test(nil))
	t.Run("nil struct pointer", f.Test((*T)(nil)))
	t.Run("non-nill struct pointer",
		subtest.ErrorIs(subtest.Failure{
			Prefix: "value is neither typed nor untyped nil",
			Got:    "*subtest_test.T\n\t{Foo:}",
		}).Test(f(&T{})),
	)
}

func TestCheckNotReflectNil(t *testing.T) {
	type T struct{ Foo string }

	f := subtest.NotReflectNil()

	t.Run("nil",
		subtest.ErrorIs(subtest.Failure{
			Prefix: "value is typed or untyped nil",
			Got:    "untyped nil",
		}).Test(f(nil)),
	)
	t.Run("nil struct pointer",
		subtest.ErrorIs(subtest.Failure{
			Prefix: "value is typed or untyped nil",
			Got:    "*subtest_test.T\n\tnil",
		}).Test(f((*T)(nil))),
	)
	t.Run("non-nill struct pointer", f.Test(&T{}))
}
