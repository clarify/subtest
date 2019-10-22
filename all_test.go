package subtest_test

import "github.com/searis/subtest"

// Set up formatting rules for package unit tests.
func init() {
	subtest.SetTypeFormatter(nil) // Explicitly use default formatter.
	subtest.SetIndent("\t")       // Makes it easier to validate failure output.
}
