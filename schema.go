package subtest

import (
	"fmt"
	"reflect"
)

// Fields provides a map of check functions.
type Fields map[interface{}]CheckFunc

// Schema allows simple validation of fields in map instances.
type Schema struct {
	// Fields map
	Fields Fields
	// Required lists required fields.
	Required map[interface{}]struct{}
	// AdditionalFields, if set, is run against all fields
	AdditionalFields CheckFunc
}

// CheckFunc returns a check functions that fails if Schema don't validate
// against the test value. Accepted test values include map[string]interface{},
// a JSON string, or JSON []byte array.
func (s Schema) CheckFunc() CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)
		if rv.Kind() != reflect.Map {
			return FailureGot("not a map", got)
		}

		var missingKeys []interface{}

		for k := range s.Required {
			if !rv.MapIndex(reflect.ValueOf(k)).IsValid() {
				missingKeys = append(missingKeys, k)
			}
		}

		var errs Errors
		var extraKeys []interface{}

		itr := rv.MapRange()
		var k, v interface{}
		var check CheckFunc
		for itr.Next() {
			k = itr.Key().Interface()
			v = itr.Value()
			check, _ = s.Fields[k]
			if check == nil {
				check = s.AdditionalFields
			}
			if check == nil {
				extraKeys = append(extraKeys, k)
				continue
			}
			if err := check(v); err != nil {
				errs = append(errs, fieldError(k, err))
			}
		}

		if len(missingKeys) > 0 {
			errs = append(errs, FailureGot("got missing required fields", missingKeys))
		}
		if len(extraKeys) > 0 {
			errs = append(errs, FailureGot("got additional fields fields", extraKeys))
		}

		if len(errs) > 0 {
			return errs
		}
		return nil
	}
}

// TODO: this error initializer may be simplified.
func fieldError(key interface{}, err error) error {
	switch et := err.(type) {
	case Errors:
		errs := make(Errors, 0, len(et))
		for _, err := range et {
			errs = append(errs, fieldError(key, err))
		}
		return errs
	case Failure:
		ret := et
		ret.Prefix = fmt.Sprintf(".%#v: %s", key, et.Prefix)
		return ret
	default:
		return fmt.Errorf(".%#v: %w", key, err)
	}
}
