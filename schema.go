package subtest

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Fields provides a map of check functions.
type Fields map[interface{}]Check

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
	// Fields map keys to checks.
	Fields Fields
	// Required, if set, contain a list of required keys. When the list is
	// explicitly defined as an empty list, no keys will be considered required.
	// When the field holds a nil value, all keys present in Fields will be
	// considered required.
	Required []interface{}
	// AdditionalFields if set, contain a check used to validate all fields
	// where the key is not present in Fields. When the field holds a nil
	// value, no additional keys are allowed. To skip validation of additional
	// keys, the Any() check can be used.
	AdditionalFields Check
}

// Check validates vf against s. For now, vf must return a map.
func (s Schema) Check(vf ValueFunc) error {
	if vf == nil {
		return FailGot("missing value function", vf)
	}
	got, err := vf()
	if err != nil {
		return FailGot("value function returns an error", err)
	}

	rv := reflect.ValueOf(got)
	switch rv.Kind() {
	case reflect.Map:
		return s.checkMap(got)
		// TODO: handle reflect.Struct
	default:
		return FailGot("not a map", got)
	}
}

func (s Schema) checkMap(got interface{}) error {
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

	var check Check
	var vf ValueFunc
	var k interface{}

	for _, rk := range rKeys {
		check = nil
		k = rk.Interface()
		vf = Value(rv.MapIndex(rk).Interface())

		if s.Fields != nil {
			check = s.Fields[k]
		}
		if check == nil {
			check = s.AdditionalFields
		}
		if check == nil {
			extraKeys = append(extraKeys, fmt.Sprintf("%#v", k))
			continue
		}
		if err := check.Check(vf); err != nil {
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
		return fmt.Errorf("%s: %w", msgSchemaMatch, errs)
	}
	return nil
}
