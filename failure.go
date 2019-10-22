package subtest

import (
	"bytes"
	"fmt"
)

// Failure is an error type that aid with consistent formatting of comparison
// based test failures.
type Failure struct {
	Prefix string
	Got    string
	Expect string
	Reject string
}

// FailureExpect formats a failure for content that is not matching some
// expected value. The package type formatter is used.
func FailureExpect(prefix string, got, expect interface{}) Failure {
	return Failure{
		Prefix: prefix,
		Got:    FormatType(got),
		Expect: FormatType(expect),
	}
}

// FailureReject formats a failure for content that is matching some rejected
// value. The package type formatter is used.
func FailureReject(prefix string, got, reject interface{}) Failure {
	return Failure{
		Prefix: prefix,
		Got:    FormatType(got),
		Reject: FormatType(reject),
	}
}

// FailureGot formats a failure for some unexpected content. The package type
// formatter is used.
func FailureGot(prefix string, got interface{}) Failure {
	return Failure{
		Prefix: prefix,
		Got:    FormatType(got),
	}
}

func (f Failure) Error() string {
	var fmtS string
	switch {
	case f.Reject != "":
		fmtS = "\n%10s: %s"
	case f.Expect != "":
		fmtS = "\n%4s: %s"
	default:
		fmtS = "\n%s: %s"
	}

	s := f.Prefix
	if f.Got != "" {
		s += fmt.Sprintf(fmtS, "got", f.Got)
	}
	if f.Expect != "" {
		s += fmt.Sprintf(fmtS, "want", f.Expect)
	}
	if f.Reject != "" {
		s += fmt.Sprintf(fmtS, "don't want", f.Reject)
	}

	return s
}

// Errors combine the output of multiple errors on separate lines.
type Errors []error

func (errs Errors) Error() string {
	var buf bytes.Buffer
	for _, err := range errs {
		fmt.Fprintf(&buf, "\n\t%q", err)
	}
	return buf.String()
}
