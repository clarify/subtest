package subtest

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Fields provides a map of check functions.
type Fields map[interface{}]CheckFunc

// OrderedKeys returns all keys in m in alphanumerical order.
func (m Fields) OrderedKeys() []interface{} {
	keys := make([]interface{}, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})

	return keys
}

// Schema allows simple validation of fields. Currently support only maps.
type Schema struct {
	// Fields map
	Fields Fields
	// Required lists required fields. Defaults to all keys in the Fields map
	// if nil. To not require any fields, set an empty slice.
	Required []interface{}
	// AdditionalFields, if set, is run against all fields
	AdditionalFields CheckFunc
}

// ValidateMap returns a check functions that fails if the test value is either
// not a map type or does not validate against s.
func (s Schema) ValidateMap() CheckFunc {
	return func(got interface{}) error {
		rv := reflect.ValueOf(got)
		if rv.Kind() != reflect.Map {
			return FailGot("not a map", got)
		}

		rKeys := rv.MapKeys()
		sort.Slice(rKeys, func(i, j int) bool {
			return fmt.Sprint(rKeys[i].Interface()) < fmt.Sprint(rKeys[j].Interface())
		})

		var errs Errors
		var extraKeys []string

		var check CheckFunc
		var k, v interface{}
		for _, rk := range rKeys {
			check = nil
			k = rk.Interface()
			v = rv.MapIndex(rk).Interface()

			if s.Fields != nil {
				check, _ = s.Fields[k]
			}
			if check == nil {
				check = s.AdditionalFields
			}
			if check == nil {
				extraKeys = append(extraKeys, fmt.Sprintf("%#v", k))
				continue
			}
			if err := check(v); err != nil {
				errs = append(errs, KeyError(k, err))
			}
		}
		if len(extraKeys) > 0 {
			errs = append(errs, Failf("got additional keys: %v", strings.Join(extraKeys, ", ")))
		}

		keySet := make(map[interface{}]struct{}, rv.Len())
		for _, rk := range rKeys {
			keySet[rk.Interface()] = struct{}{}
		}

		var required []interface{}
		if s.Required == nil {
			required = s.Fields.OrderedKeys()
		} else {
			required = s.Required
		}

		var missingKeys []string
		for _, k := range required {
			_, ok := keySet[k]
			if !ok {
				missingKeys = append(missingKeys, fmt.Sprintf("%#v", k))
			}
		}
		if len(missingKeys) > 0 {
			errs = append(errs, Failf("missing required keys: %v", strings.Join(missingKeys, ", ")))
		}

		if len(errs) > 0 {
			return PrefixError{
				Key:     "value not matching schema",
				Err:     errs,
				Newline: true,
			}
		}
		return nil
	}
}
