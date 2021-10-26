package jsondiff_test

import (
	"encoding/json"
	"testing"

	"github.com/clarify/subtest"
	pkg "github.com/clarify/subtest/examples/jsondiff"
)

func TestEqualJSON(t *testing.T) {
	// This test validates that the EqualJSON check functions work; it is not a
	// direct example of usage.

	const (
		a = `{"foo": "A", "bar": "B"}`
		b = `{"foo": "A", "bar": "C"}`
		c = `{"foo": "A"}`
	)
	t.Log("given a is not equal to b")
	t.Log("given c is a JSON sub-set of a")
	t.Log("given c is a JSON sub-set of b")

	t.Run("when a is cheked against EqualJSON(a)", func(t *testing.T) {
		err := pkg.EqualJSON(a)(a)
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when []byte(a) is checked against EqualJSON(a)", func(t *testing.T) {
		err := pkg.EqualJSON(a)([]byte(a))
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when json.RawMessage(a) is checked against EqualJSON(a)", func(t *testing.T) {
		err := pkg.EqualJSON(a)(json.RawMessage(a))
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when b is cheked against EqualJSON(a)", func(t *testing.T) {
		err := pkg.EqualJSON(a)(b)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
	t.Run("when c is cheked against EqualJSON(a)", func(t *testing.T) {
		err := pkg.EqualJSON(a)(c)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
	t.Run("when a is cheked against EqualJSON(c)", func(t *testing.T) {
		err := pkg.EqualJSON(c)(a)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
}

func TestSubsetOfJSON(t *testing.T) {
	// This test validates that the SubsetOfJSON check functions work; it is not
	// a direct example of usage.

	const (
		a = `{"foo": "A", "bar": "B"}`
		b = `{"foo": "A", "bar": "C"}`
		c = `{"foo": "A"}`
	)
	t.Log("given a is not equal to b")
	t.Log("given c is a sub-set of a")
	t.Log("given c is a sub-set of b")

	t.Run("when a is cheked against SubsetOfJSON(a)", func(t *testing.T) {
		err := pkg.SubsetOfJSON(a)(a)
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when []byte(a) is checked against SubsetOfJSON(a)", func(t *testing.T) {
		err := pkg.SubsetOfJSON(a)([]byte(a))
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when json.RawMessage(a) is checked against SubsetOfJSON(a)", func(t *testing.T) {
		err := pkg.SubsetOfJSON(a)(json.RawMessage(a))
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when b is checked against SubsetOfJSON(a)", func(t *testing.T) {
		err := pkg.SubsetOfJSON(a)(b)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
	t.Run("when c is checked against SubsetOfJSON(a)", func(t *testing.T) {
		err := pkg.SubsetOfJSON(a)(c)
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when a is cheked against SubsetOfJSON(c)", func(t *testing.T) {
		err := pkg.SubsetOfJSON(c)(a)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
}

func TestSupersetOfJSON(t *testing.T) {
	// This test validates that the SupersetOfJSON check functions work; it is
	// not a direct example of usage.

	const (
		a = `{"foo": "A", "bar": "B"}`
		b = `{"foo": "A", "bar": "C"}`
		c = `{"foo": "A"}`
	)
	t.Log("given a is not equal to b")
	t.Log("given c is a JSON sub-set of a")
	t.Log("given c is a JSON sub-set of b")

	t.Run("when a is cheked against SupersetOfJSON(a)", func(t *testing.T) {
		err := pkg.SupersetOfJSON(a)(a)
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when []byte(a) is checked against SupersetOfJSON(a)", func(t *testing.T) {
		err := pkg.SupersetOfJSON(a)([]byte(a))
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when json.RawMessage(a) is checked against SupersetOfJSON(a)", func(t *testing.T) {
		err := pkg.SupersetOfJSON(a)(json.RawMessage(a))
		t.Run("then the check should pass", subtest.Value(err).NoError())
	})
	t.Run("when b is checked against SupersetOfJSON(a)", func(t *testing.T) {
		err := pkg.SupersetOfJSON(a)(b)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
	t.Run("when c is checked against SupersetOfJSON(a)", func(t *testing.T) {
		err := pkg.SupersetOfJSON(a)(c)
		t.Run("then the check should fail", subtest.Value(err).Error())
	})
	t.Run("when a is checked against SupersetOfJSON(c)", func(t *testing.T) {
		err := pkg.SupersetOfJSON(c)(a)
		t.Run("then the check should fail", subtest.Value(err).NoError())
	})
}
