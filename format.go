package subtest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var fmtCfg = struct {
	f      func(v ...interface{}) string
	indent string
}{indent: "    "}

// SetTypeFormatter replaces the type formatter used by the package. This
// function is not thread-safe, and should be called as part of initialization
// only. E.g. in a test package init function.
func SetTypeFormatter(f func(...interface{}) string) {
	fmtCfg.f = f
}

// SetIndent sets a string to use in package error and type formatting. The
// default is four spaces, as that's what's used by the Go test-runner.  This
// function is not thread-safe, and should be called as part of initialization
// only. E.g. in a test package init function.
func SetIndent(s string) {
	fmtCfg.indent = s
}

func formatIndentedType(v interface{}) string {
	switch vt := v.(type) {
	case nil:
		return "untyped nil"
	case error:
		return fmt.Sprintf("%T\n%s", v, indentString(vt.Error()))
	default:
		return fmt.Sprintf("%T\n%s", v, indentString(FormatType(v)))
	}
}

// FormatType formats a type using the configured type formatter for the
// package.
func FormatType(v interface{}) string {
	var f func(...interface{}) string
	f = fmtCfg.f

	if f == nil {
		return defaultTypeFormatter(v)
	}

	return f(v)
}

func defaultTypeFormatter(v interface{}) string {

	switch vt := v.(type) {
	case nil:
		return "nil"
	case error:
		return vt.Error()
	case []byte, json.RawMessage:
		return fmt.Sprintf("`%s`", v)
	case string, fmt.Stringer:
		return quoteString(fmt.Sprintf("%s", v))
	}

	rv := reflect.ValueOf(v)
	switch {
	case rv.Kind() == reflect.Ptr && rv.IsNil():
		return "nil"
	case rv.Kind() == reflect.Ptr:
		return fmt.Sprintf("%+v", rv.Elem().Interface())
	}
	return fmt.Sprintf("%+v", v)
}

func quoteString(s string) string {
	if strings.Contains(s, "\n") {
		return "`" + s + "`"
	}
	return fmt.Sprintf("%q", s)
}

func indentString(s string) string {
	sIndent := fmtCfg.indent
	return sIndent + strings.ReplaceAll(s, "\n", "\n"+sIndent)
}
