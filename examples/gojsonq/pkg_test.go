package gojsonq_test

import (
	"errors"
	"testing"

	"github.com/searis/subtest"
	"gopkg.in/thedevsaddam/gojsonq.v2"
)

func TestJSON1(t *testing.T) {
	// j allows queries against v, but there is a catch; for each query j must
	// be reset to allow for the next query. This test shows how to handle this
	// manually prior to each test.

	const v = `{"name": "foo", "type": "bar"}`

	j := gojsonq.New().JSONString(v)
	t.Run(".", subtest.Value(j).Test(subtest.CheckFunc(validSchema)))
	j.Reset()
	t.Run(".name", subtest.Value(j.Find("name")).DeepEqual("foo"))
	j.Reset()
	t.Run(".type", subtest.Value(j.Find("type")).DeepEqual("bar"))
}

func validSchema(got interface{}) error {
	j := got.(*gojsonq.JSONQ)

	var errs subtest.Errors

	if cnt, expect := j.Count(), 2; cnt != expect {
		errs = append(errs, subtest.FailExpect("incorrect number of keys", cnt, expect))

	}

	j.Reset()
	if _, ok := j.Find("name").(string); !ok {
		errs = append(errs, errors.New(".name not a string"))
	}

	j.Reset()
	if _, ok := j.Find("type").(string); !ok {
		errs = append(errs, errors.New(".type not a string"))
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func TestJSON2(t *testing.T) {
	// This test demonstrates that a check don't have to be defined against a
	// static value. When it's useful, a new instance can be generated for each
	// test.

	const v = `{"name": "foo", "type": "bar"}`

	vf := subtest.ValueFunc(func() (interface{}, error) {
		return gojsonq.New().JSONString(v), nil
	})

	t.Run(".", vf.Test(subtest.CheckFunc(validSchema)))
	t.Run(".name", vf.Test(jsonPathEqual("name", "foo")))
	t.Run(".type", vf.Test(jsonPathEqual("type", "bar")))
}

func jsonPathEqual(path string, expect interface{}) subtest.CheckFunc {
	check := subtest.DeepEqual(expect)

	return func(got interface{}) error {
		j := got.(*gojsonq.JSONQ)
		return check(j.Find(path))
	}
}
