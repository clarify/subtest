package subtest_test

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/searis/subtest"
)

func TestDeepEqual(t *testing.T) {
	t.Run("given check DeepEqual(true)", func(t *testing.T) {
		cf := subtest.DeepEqual(true)
		t.Run("when cheking against true", func(t *testing.T) {
			vf := subtest.Value(cf(true))
			t.Run("then there should be no failure", vf.NoError())
		})
		t.Run("when cheking against false", func(t *testing.T) {
			vf := subtest.Value(cf(false))
			t.Run("then an appropriate error should be returned", vf.ErrorIs(subtest.Failure{
				Prefix: "not deep equal",
				Got:    "bool\n\tfalse",
				Expect: "bool\n\ttrue",
			}))
		})
		t.Run("when testing against true with different syntax", func(t *testing.T) {
			v := true
			vf := subtest.Value(v)
			t.Run("then <ValueFunc>.<CheckFuncName> should not fail", vf.DeepEqual(true))
			t.Run("the <ValueFunc>.Test(<CheckFunc>) should not fail", vf.Test(subtest.DeepEqual(true))) // equivalent
		})
	})

	t.Run("given a nested struct type T", func(t *testing.T) {
		type T struct {
			A string
			B map[string]string
		}

		t.Run("when checking a non-nil *T value", func(t *testing.T) {
			v := &T{A: "a", B: map[string]string{"C": "D"}}

			t.Run("then it should match an equivalent *T value", // equal value.
				subtest.Value(v).DeepEqual(&T{A: "a", B: map[string]string{"C": "D"}}),
			)
			t.Run("then it should not match a different *T value", // different value.
				subtest.Value(subtest.DeepEqual(&T{A: "a", B: map[string]string{"C": "E"}})(v)).Error(),
			)
			t.Run("then it should not match an equivalent T value", // equal value, different type
				subtest.Value(subtest.DeepEqual(T{A: "a", B: map[string]string{"C": "D"}})(v)).Error(),
			)

		})
	})

}

func TestNotDeepEqual(t *testing.T) {
	t.Run("given check NotDeepEqual(false)", func(t *testing.T) {
		cf := subtest.NotDeepEqual(false)
		t.Run("when cheking against true", func(t *testing.T) {
			vf := subtest.Value(cf(true))
			t.Run("then there should be no failure", vf.NoError())
		})
		t.Run("when cheking against false", func(t *testing.T) {
			vf := subtest.Value(cf(false))
			t.Run("then an appropriate error should be returned", vf.ErrorIs(subtest.Failure{
				Prefix: "deep equal",
				Got:    "bool\n\tfalse",
				Reject: "bool\n\tfalse",
			}))
		})
		t.Run("when testing against false with different syntax", func(t *testing.T) {
			v := false
			vf := subtest.Value(v)
			t.Run("then <ValueFunc>.<CheckFuncName> should not fail", vf.NotDeepEqual(true))
			t.Run("the <ValueFunc>.Test(<CheckFunc>) should not fail", vf.Test(subtest.NotDeepEqual(true))) // equivalent

		})
	})

	t.Run("given a nested struct type T", func(t *testing.T) {
		type T struct {
			A string
			B map[string]string
		}

		t.Run("when checking a non-nil *T value", func(t *testing.T) {
			v := &T{A: "a", B: map[string]string{"C": "D"}}

			t.Run("then it should not accept an equivalent *T value",
				subtest.Value(subtest.NotDeepEqual(&T{A: "a", B: map[string]string{"C": "D"}})(v)).Error(),
			)
			t.Run("then it should accept a different *T value",
				subtest.Value(v).NotDeepEqual(&T{A: "a", B: map[string]string{"C": "E"}}),
			)
			t.Run("then it should not match an equivalent T value",
				subtest.Value(v).NotDeepEqual(T{A: "a", B: map[string]string{"C": "D"}}),
			)

		})
	})
}

func TestCheckReflectNil(t *testing.T) {
	type T struct{ Foo string }

	cf := subtest.ReflectNil()

	t.Run("when cheking against untyped nil", func(t *testing.T) {
		vf := subtest.Value(cf(nil))
		t.Run("then it should not fail", vf.NoError())
	})
	t.Run("when cheking against a nil struct pointer", func(t *testing.T) {
		vf := subtest.Value(cf((*T)(nil)))
		t.Run("then it should not fail", vf.NoError())
	})

	t.Run("when cheking against a non-nil struct pointer", func(t *testing.T) {
		vf := subtest.Value(cf(&T{}))
		t.Run("then an appropriate error should be returned", vf.ErrorIs(subtest.Failure{
			Prefix: "neither typed nor untyped nil",
			Got:    "*subtest_test.T\n\t{Foo:}",
		}))
	})

}

func TestCheckNotReflectNil(t *testing.T) {
	type T struct{ Foo string }

	cf := subtest.NotReflectNil()

	t.Run("when cheking against untyped nil", func(t *testing.T) {
		vf := subtest.Value(cf(nil))
		t.Run("then an appropriate error should be returned", vf.ErrorIs(subtest.Failure{
			Prefix: "typed or untyped nil",
			Got:    "untyped nil",
		}))
	})
	t.Run("when cheking against a nil struct pointer", func(t *testing.T) {
		vf := subtest.Value(cf((*T)(nil)))
		t.Run("then an appropriate error should be returned", vf.ErrorIs(subtest.Failure{
			Prefix: "typed or untyped nil",
			Got:    "*subtest_test.T\n\tnil",
		}))
	})

	t.Run("when cheking against a non-nil struct pointer", func(t *testing.T) {
		vf := subtest.Value(cf(&T{}))
		t.Run("then it should not fail", vf.NoError())
	})

}

func TestLessThan(t *testing.T) {
	t.Run("given a value float64(42)", func(t *testing.T) {
		v := float64(42)
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.LessThan(43)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.LessThan(42)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not less than 42.000000",
				Got:    "float64\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a value int16(42)", func(t *testing.T) {
		v := int16(42)
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.LessThan(43)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.LessThan(42)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not less than 42.000000",
				Got:    "int16\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestLessThanOrEqual(t *testing.T) {
	t.Run("given a value float64(42)", func(t *testing.T) {
		v := float64(42)
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.LessThanOrEqual(43)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.LessThanOrEqual(42)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.LessThanOrEqual(41)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not less than or equal to 41.000000",
				Got:    "float64\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a value int16(42)", func(t *testing.T) {
		v := int16(42)
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.LessThanOrEqual(42)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.LessThanOrEqual(41)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not less than or equal to 41.000000",
				Got:    "int16\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestGreaterThan(t *testing.T) {
	t.Run("given a value float64(42)", func(t *testing.T) {
		v := float64(42)
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.GreaterThan(41)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.GreaterThan(42)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not greater than 42.000000",
				Got:    "float64\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a value int16(42)", func(t *testing.T) {
		v := int16(42)
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.GreaterThan(41)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.GreaterThan(42)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not greater than 42.000000",
				Got:    "int16\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestGreaterThanOrEqual(t *testing.T) {
	t.Run("given a value float64(42)", func(t *testing.T) {
		v := float64(42)
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.GreaterThanOrEqual(41)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.GreaterThanOrEqual(42)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.GreaterThanOrEqual(43)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not greater than or equal to 43.000000",
				Got:    "float64\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a value int16(42)", func(t *testing.T) {
		v := int16(42)
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.GreaterThanOrEqual(42)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.GreaterThanOrEqual(43)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not greater than or equal to 43.000000",
				Got:    "int16\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestNumericEqual(t *testing.T) {
	t.Run("given a value float64(42)", func(t *testing.T) {
		v := float64(42)
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.NumericEqual(41)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not numeric equal",
				Got:    "float64\n\t42",
				Expect: "float64\n\t41",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.NumericEqual(42)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.NumericEqual(43)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "not numeric equal",
				Got:    "float64\n\t42",
				Expect: "float64\n\t43",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestNotNumericEqual(t *testing.T) {
	t.Run("given a value float64(42)", func(t *testing.T) {
		v := float64(42)
		t.Run("when cheking against 41", func(t *testing.T) {
			cf := subtest.NotNumericEqual(41)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.NotNumericEqual(42)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "numeric equal",
				Got:    "float64\n\t42",
				Reject: "float64\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.NotNumericEqual(43)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
	})
	t.Run("given a value int16(42)", func(t *testing.T) {
		v := int16(42)
		t.Run("when cheking against 42", func(t *testing.T) {
			cf := subtest.NotNumericEqual(42)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "numeric equal",
				Got:    "int16\n\t42",
				Reject: "float64\n\t42",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against 43", func(t *testing.T) {
			cf := subtest.NotNumericEqual(43)
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
	})
}
func TestNotBefore(t *testing.T) {
	tz := time.FixedZone("Europe/Oslo", 2*3600)
	t1 := time.Date(1985, 12, 19, 18, 15, 0, 0, tz)

	t.Run("given a string value", func(t *testing.T) {
		v := "1985-12-19T18:15:00.0+02:00"
		t.Run("when cheking against any time", func(t *testing.T) {
			cf := subtest.NotBefore(t1)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "type is not time.Time or *time.Time",
				Got:    "string\n\t\"1985-12-19T18:15:00.0+02:00\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a time value", func(t *testing.T) {
		v := t1
		t.Run("when cheking against an earlier time", func(t *testing.T) {
			cf := subtest.NotBefore(t1.Add(-time.Second))
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against a later time", func(t *testing.T) {
			cf := subtest.NotBefore(t1.Add(time.Second))
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "time before 1985-12-19 18:15:01 +0200 Europe/Oslo",
				Got:    "time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.NotBefore(t1.UTC())
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
	})
	t.Run("given a time pointer value", func(t *testing.T) {
		v := &t1
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.NotBefore(t1.UTC())
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
	})
}

func TestBefore(t *testing.T) {
	tz := time.FixedZone("Europe/Oslo", 2*3600)
	t1 := time.Date(1985, 12, 19, 18, 15, 0, 0, tz)

	t.Run("given a string value", func(t *testing.T) {
		v := "1985-12-19T18:15:00.0+02:00"
		t.Run("when cheking against any time", func(t *testing.T) {
			cf := subtest.Before(t1)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "type is not time.Time or *time.Time",
				Got:    "string\n\t\"1985-12-19T18:15:00.0+02:00\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a time value", func(t *testing.T) {
		v := t1
		t.Run("when cheking against an earlier time", func(t *testing.T) {
			cf := subtest.Before(t1.Add(-time.Second))
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "time not before 1985-12-19 18:14:59 +0200 Europe/Oslo",
				Got:    "time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against a later time", func(t *testing.T) {
			cf := subtest.Before(t1.Add(time.Second))
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.Before(t1.UTC())
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "time not before 1985-12-19 16:15:00 +0000 UTC",
				Got:    "time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a time pointer value", func(t *testing.T) {
		v := &t1
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.Before(t1.UTC())
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "time not before 1985-12-19 16:15:00 +0000 UTC",
				Got:    "*time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestNotTimeEqual(t *testing.T) {
	tz := time.FixedZone("Europe/Oslo", 2*3600)
	t1 := time.Date(1985, 12, 19, 18, 15, 0, 0, tz)

	t.Run("given a string value", func(t *testing.T) {
		v := "1985-12-19T18:15:00.0+01:00"
		t.Run("when cheking against any time", func(t *testing.T) {
			cf := subtest.TimeEqual(t1)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "type is not time.Time or *time.Time",
				Got:    "string\n\t\"1985-12-19T18:15:00.0+01:00\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a time value", func(t *testing.T) {
		v := t1
		t.Run("when cheking against an earlier time", func(t *testing.T) {
			cf := subtest.NotTimeEqual(t1.Add(-time.Second))
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against a later time", func(t *testing.T) {
			cf := subtest.NotTimeEqual(t1.Add(time.Second))
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.NotTimeEqual(t1.UTC())
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "times not equal",
				Got:    "time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
				Reject: "time.Time\n\t\"1985-12-19 16:15:00 +0000 UTC\"",
			}
			t.Run("then it should not fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a time pointer value", func(t *testing.T) {
		v := &t1
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.NotTimeEqual(t1.UTC())
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "times not equal",
				Got:    "*time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
				Reject: "time.Time\n\t\"1985-12-19 16:15:00 +0000 UTC\"",
			}
			t.Run("then it should not fail", vf.ErrorIs(expect))
		})
	})
}

func TestTimeEqual(t *testing.T) {
	tz := time.FixedZone("Europe/Oslo", 2*3600)
	t1 := time.Date(1985, 12, 19, 18, 15, 0, 0, tz)

	t.Run("given a string value", func(t *testing.T) {
		v := "1985-12-19T18:15:00.0+01:00"
		t.Run("when cheking against any time", func(t *testing.T) {
			cf := subtest.TimeEqual(t1)
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "type is not time.Time or *time.Time",
				Got:    "string\n\t\"1985-12-19T18:15:00.0+01:00\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
	t.Run("given a time value", func(t *testing.T) {
		v := t1
		t.Run("when cheking against an earlier time", func(t *testing.T) {
			cf := subtest.TimeEqual(t1.Add(-time.Microsecond))
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "times not equal",
				Got:    "time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
				Expect: "time.Time\n\t\"1985-12-19 18:14:59.999999 +0200 Europe/Oslo\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against a later time", func(t *testing.T) {
			cf := subtest.TimeEqual(t1.Add(time.Microsecond))
			vf := subtest.Value(cf(v))
			expect := subtest.Failure{
				Prefix: "times not equal",
				Got:    "time.Time\n\t\"1985-12-19 18:15:00 +0200 Europe/Oslo\"",
				Expect: "time.Time\n\t\"1985-12-19 18:15:00.000001 +0200 Europe/Oslo\"",
			}
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.TimeEqual(t1.UTC())
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
	})
	t.Run("given a time pointer value", func(t *testing.T) {
		v := &t1
		t.Run("when cheking against a semantically equivalent time", func(t *testing.T) {
			cf := subtest.TimeEqual(t1.UTC())
			vf := subtest.Value(cf(v))
			t.Run("then it should not fail", vf.NoError())
		})
	})
}

func TestRegexp(t *testing.T) {
	t.Run("given a regular expression check function", func(t *testing.T) {
		cf := subtest.MatchRegexp(regexp.MustCompile(`^"f.*a.?r"$`))
		t.Run("when cheking against a non matching string", func(t *testing.T) {
			vf := subtest.Value(cf(`"foo"`))
			t.Run("then it should fail", vf.Error())
		})
		t.Run("when cheking against a matching string", func(t *testing.T) {
			vf := subtest.Value(cf(`"foobar"`))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against a matching []byte", func(t *testing.T) {
			vf := subtest.Value(cf([]byte(`"foobar"`)))
			t.Run("then it should not fail", vf.NoError())
		})
		t.Run("when cheking against a matching json.RawMessage", func(t *testing.T) {
			vf := subtest.Value(cf(json.RawMessage(`"foobar"`)))
			t.Run("then it should not fail", vf.NoError())
		})
	})
}

func TestContainsMatch(t *testing.T) {
	t.Run("given a slice []string{a, b}", func(t *testing.T) {
		v := []string{"a", "b"}
		t.Run("when checking against a", func(t *testing.T) {
			t.Run("then it should match", // equal value.
				subtest.Value(v).ContainsMatch(subtest.DeepEqual("a")),
			)
		})
		t.Run("when checking against b", func(t *testing.T) {
			t.Run("then it should match", // equal value.
				subtest.Value(v).ContainsMatch(subtest.DeepEqual("b")),
			)
		})
		t.Run("when checking against c", func(t *testing.T) {
			cf := subtest.ContainsMatch(
				subtest.DeepEqual("c"),
			)
			vf := subtest.Value(cf(v))
			expect := subtest.FailExpect("does not match any elements", v, "c")
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when checking against a and 42", func(t *testing.T) {
			cf := subtest.AllOff{
				subtest.ContainsMatch(subtest.DeepEqual("a")),
				subtest.ContainsMatch(subtest.DeepEqual(42)),
			}
			err := cf.Check(subtest.Value(v))
			expect := subtest.Errors{
				subtest.FailExpect("does not match any elements", v, 42),
			}
			t.Run("then it should fail", subtest.Value(err).ErrorIs(expect))
		})
	})

	t.Run("given a slice []interface{}{true, 1, b}", func(t *testing.T) {
		v := []interface{}{true, 1, "b"}
		t.Run("when checking against 1", func(t *testing.T) {
			t.Run("then it should match",
				subtest.Value(v).ContainsMatch(subtest.DeepEqual(1)),
			)
		})
		t.Run("when checking against true", func(t *testing.T) {
			t.Run("then it should match",
				subtest.Value(v).ContainsMatch(subtest.DeepEqual(true)),
			)
		})
		t.Run("when checking against > 0", func(t *testing.T) {
			t.Run("then it should match",
				subtest.Value(v).ContainsMatch(subtest.GreaterThan(0)),
			)
		})
		t.Run("when checking against true and >= 42", func(t *testing.T) {
			cf := subtest.AllOff{
				subtest.ContainsMatch(subtest.DeepEqual(true)),
				subtest.ContainsMatch(subtest.GreaterThanOrEqual(42)),
			}
			err := cf.Check(subtest.Value(v))
			expect := subtest.Errors{
				subtest.FailGot("does not match any elements", v),
			}
			t.Run("then it should fail", subtest.Value(err).ErrorIs(expect))
		})
	})

	t.Run("given a slice of complex types", func(t *testing.T) {
		type T struct {
			A string
			B map[string]string
		}
		v := []interface{}{
			&T{A: "a", B: map[string]string{"C": "D"}},
			&T{A: "b", B: map[string]string{"E": "F"}},
		}

		t.Run("when checking against an instance in the slice", func(t *testing.T) {
			t.Run("then it should match",
				subtest.Value(v).ContainsMatch(subtest.DeepEqual(&T{A: "a", B: map[string]string{"C": "D"}})),
			)
		})
		t.Run("when checking against an instance not in the slice", func(t *testing.T) {
			cf := subtest.ContainsMatch(
				subtest.DeepEqual(&T{A: "c", B: map[string]string{"C": "D"}}),
			)
			vf := subtest.Value(cf(v))
			expect := subtest.FailExpect("does not match any elements", v, &T{A: "c", B: map[string]string{"C": "D"}})
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
	})
}

func TestContains(t *testing.T) {
	t.Run("given a slice []string{a, b}", func(t *testing.T) {
		v := []string{"a", "b"}
		t.Run("when checking against a", func(t *testing.T) {
			t.Run("then it should match", // equal value.
				subtest.Value(v).Contains("a"),
			)
		})
		t.Run("when checking against b", func(t *testing.T) {
			t.Run("then it should match", // equal value.
				subtest.Value(v).Contains("b"),
			)
		})
		t.Run("when checking against c", func(t *testing.T) {
			cf := subtest.Contains("c")
			vf := subtest.Value(cf(v))
			expect := subtest.FailExpect("does not match any elements", v, "c")
			t.Run("then it should fail", vf.ErrorIs(expect))
		})
		t.Run("when checking against a and 42", func(t *testing.T) {
			cf := subtest.AllOff{
				subtest.Contains("a"),
				subtest.Contains(42),
			}
			err := cf.Check(subtest.Value(v))
			expect := subtest.Errors{
				subtest.FailExpect("does not match any elements", v, 42),
			}
			t.Run("then it should fail", subtest.Value(err).ErrorIs(expect))
		})
	})

	t.Run("given a slice []interface{}{true, 1, b}", func(t *testing.T) {
		v := []interface{}{true, 1, "b"}
		t.Run("when checking against 1", func(t *testing.T) {
			t.Run("then it should match",
				subtest.Value(v).Contains(1),
			)
		})
		t.Run("when checking against true", func(t *testing.T) {
			t.Run("then it should match",
				subtest.Value(v).Contains(true),
			)
		})
		t.Run("when checking against true, 42 and 43", func(t *testing.T) {
			cf := subtest.AllOff{
				subtest.Contains(true),
				subtest.Contains(42),
				subtest.Contains(43),
			}
			err := cf.Check(subtest.Value(v))
			expect := subtest.Errors{
				subtest.FailExpect("does not match any elements", v, 42),
				subtest.FailExpect("does not match any elements", v, 43),
			}
			t.Run("then it should fail", subtest.Value(err).ErrorIs(expect))
		})
	})
}
