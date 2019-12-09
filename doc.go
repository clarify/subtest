// Package subtest provides a way of intializing small test functions suitable
// for use as with the (*testing.T).Run method. Tests using package-defined
// check functions can generally be initialized in two main ways. here is an
// example using DeepEquals:
//
//     // Short-hand syntax for built-in check functions.
//     t.Run("got==expect", subtest.Value(got).DeepEquals(expect))
//
//     // Long syntax.
//      t.Run("got==expect", subtest.Value(got).Test(subtest.DeepEquals(expect)))
//
// Custom CheckFunc implementations can also be turned into tests:
//
//    t.Run("got==expect", subtest.Value(got).Test(func(got interface{}) error {
//        if got != expect {
//            return subtest.FailExpect("not equal", got, expect)
//        }
//    }))
//
// Experimentally, any function that takes no parameter and returns an error can
// also be converted to a test:
//
//    t.Run("got==expect", subtest.Test(func() error {
//        if got != expect {
//            return subtest.FailExpect("not plain equal", got, expect)
//        }
//    }))
//
// When necessary, either CheckFunc middleware or a custom ValueFunc instances
// can be used to prepare or transform the test value before evaluation.
//
// PS! Note that the all experimental syntax may be removed in a later version.
package subtest
