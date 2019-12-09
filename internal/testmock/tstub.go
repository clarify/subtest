package testmock

import (
	"strings"
	"testing"
)

// T mimics the *testing.T type in example tests.
type T struct {
	Name string
}

// Run mimics the (*testing.T).Run method.
func (t T) Run(name string, f func(t *testing.T)) {
	testing.RunTests(matchAll, []testing.InternalTest{{
		F:    f,
		Name: t.Name + "/" + rewrite(name),
	}})
}

func matchAll(pat, str string) (bool, error) {
	return true, nil
}

// rewrite is only an approximate replacement of the name testing pacakage's
// name sanitizing. https://golang.org/src/testing/match.go?h=rewrite#L135
func rewrite(name string) string {
	r := strings.NewReplacer(
		" ", "_",
		"\n", "_",
		"\r", "_",
		"\t", "_",
		"\n", "_",
		"\v", "_",
		"\f", "_",
		"\r", "_",
		" ", "_",
		"\u0085", "_",
		"\u00A0", "_",
		"\u1680", "_",
	)
	return r.Replace(name)
}
