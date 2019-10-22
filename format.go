package subtest

import (
	"fmt"
	"reflect"
	"sync"
)

var formatCfg struct {
	sync.RWMutex
	f func(v ...interface{}) string
}

// FormatType formats a type using the configured type formatter for the
// package. Note that the default formatter is limited to use a compact format
// without syntax highlighting, and does not expand nested pointer values.
func FormatType(v interface{}) string {
	formatCfg.RLock()
	defer formatCfg.RUnlock()

	if formatCfg.f == nil {
		return defaultTypeFormatter(v)
	}

	return formatCfg.f(v)
}

func defaultTypeFormatter(v interface{}) string {
	switch v.(type) {
	case Errors:
		return fmt.Sprintf("%[1]T%[1]s", v)
	case error, string, []byte, fmt.Stringer:
		return fmt.Sprintf("%[1]T\n\t%[1]q", v)
	case nil:
		return "untyped nil"
	}

	rv := reflect.ValueOf(v)
	switch {
	case rv.Kind() == reflect.Ptr && rv.IsNil():
		return fmt.Sprintf("%T\n\tnil", v)
	case rv.Kind() == reflect.Ptr:
		return fmt.Sprintf("%T\n\t%+v", v, rv.Elem().Interface())
	}

	return fmt.Sprintf("%[1]T\n\t%+[1]v", v)
}

// SetTypeFormatter replaces the type formatter used by the package.
func SetTypeFormatter(f func(...interface{}) string) {
	formatCfg.Lock()
	formatCfg.f = f
	formatCfg.Unlock()
}
