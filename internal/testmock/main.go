package testmock

import (
	"flag"
	"os"
	"testing"
)

// VerboseMainTest allows override the default test runner to enforce the
// verbose settings.
func VerboseMainTest(m *testing.M) {
	var hasVerbose bool
FOR:
	for _, arg := range os.Args {
		switch arg {
		case "-v", "-test.v":
			hasVerbose = true
			break FOR
		}
	}
	if !hasVerbose {
		os.Args = append(os.Args, "-test.v")
	}

	flag.Parse()
	os.Exit(m.Run())
}
