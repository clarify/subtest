# subtest
[![Go Report Card](https://goreportcard.com/badge/github.com/searis/subtest)](https://goreportcard.com/report/github.com/searis/subtest)
[![Documentation](https://godoc.org/github.com/searis/subtest?status.svg)](http://godoc.org/github.com/searis/subtest)


`subtest` is a minimalist Go test-utility package used to initializing small test functions for use with the Go sub-tests feature. You can read more about Go sub-tests [here][go-sub-test].

The sub-package `subjson` defines middleware for parsing values from JSON before performing checks.

## Introduction

`subtest` was motivated by a desire to make it easier write [Given-When-Then][gwt] (GWT) style tests in Go on top of the built-in test-runner, and without a DSL.

[gwt]: https://martinfowler.com/bliki/GivenWhenThen.html

GWT is a naming schema for tests that attempts to clarify three aspects:

1. Given: How does the world look like before we do an action.
2. When: What do we do to affect the world.
3. Then: How should the world look like after the action.

These conditions can be nested to test different scenarios. Here is an example:

    Given foo is 42
        When dividing by 6
            Then the result should be 7
        When dividing by 9
            Then the result should be less than 5
            Then the result should be more than 4

A common believe is that to write GWT style tests in Go, you should use a [Behavior-Driven-Development][bdd] framework and associated [Domain-Specific-Language][dsl]. However, we will argue that this simply isn't true.

[bdd]: https://en.wikipedia.org/wiki/Behavior-driven_development
[dsl]: https://en.wikipedia.org/wiki/Domain-specific_language

One of the problems with many [Behavior-Driven-Development][bdd] frameworks in Go, is that they tend to rely on their own test-runner and sub-test logic. While there could be good reasons for using this, it also comes with a price: tooling expecting the default test runner simply does not cope. This is true either we are talking about a CI that fail to parse sub-test results into JUnit summaries, an IDE that fail to insert links for navigating to the failing code, or a command-line tool for rerunning failing sub-tests by name. This has a real impact on how well test-results can be understood when used with tooling.

As it turns out, you can actually write GWT-style tests not only without a BDD framework or DSL, but without a framework or library what so ever:

```go
func TestFoo(t *testing.T) {
    t.Run("given foo is 42", func(t *testing.T) {
        const foo = 42
        t.Run("when dividing by 6", func (t *testing.T) {
            v := float64(foo) / 6
            t.Run("then the result should be 7", func(t *testing.T) {
                if expect := float64(7); v != expect {
                    t.Fatalf("\n got: %d\nwant: %d", v, expect)
                }
            })
        })
        t.Run("when dividing by 9", func (t *testing.T) {
            v := float64(foo) / 9
            t.Run("then the result should be greater than 4" func(t *testing.T) {
                if expect := float64(4); v > expect {
                    t.Fatalf("\n got: %d\nwant > %d", v, expect)
                }
            })
            t.Run("then the result should be less than 5" func(t *testing.T) {
                if expect := float64(5); v < expect {
                    t.Fatalf("\n got: %d\nwant < %d", v, expect)
                }
            })
        })
    })
}
```

While doing this is fine, it can quickly become repetitive. It can also become challenging to maintain consistent output for failing tests over time; in particularly so for a growing team.

By now you might think that `subtest` is going to improve how to write GWT style tests without a BDD-style framework, and you are right. However, there is nothing within the design of `subtest` that restricts it to handling GWT style tests. Instead, `subtest` is a generalized test utility package for generating sub-tests, no matter your style preferences. In other words, you can use `subtest` to write GWT style tests, table-driven tests, or some other style that you prefer.

Here is a version of TestFoo with `subtest`:

```go
func TestFoo(t *testing.T) {

    t.Run("given foo is 42", func(t *testing.T) {
        const foo = 42
        t.Run("when dividing by 6", func (t *testing.T) {
            vf := subtest.Value(float64(foo) / 6)
            t.Run("then the result should be 7", vf.DeepEqual(float64(7)))
        })
        t.Run("when dividing by 9", func (t *testing.T) {
            vf := subtest.Value(float64(foo) / 6)
            t.Run("then the result should be greater than 4", vf.GreaterThan(4))
            t.Run("then the result should be less than 5", vf.LessThan(5))
        })
    })
}
```

## Short names

Some common short names used in examples:

- `v`: a value
- `vf`: a value function
- `cf`: a check function

## Examples

[go-sub-test]: https://blog.golang.org/subtests

Here is some example test code using `subtest` and `subjson`:

```go
 // TestFooExplicit shows the general syntax for using subtest written in very
 // explicit steps. All usage of subtest follows these basic steps. Tests are
 // not usually written as explicit though; see TestFooFewerLines or
 // TestFooShortHand for more realistic variants.
func TestFooExplicit(t *testing.T) {
    // 1. We declare the value we want to check; usually the result of an
    // operation or action.
    v := "foo"

    // 2. We initialize a value function for the value that we want to check.
    vf := subtest.Value(v)

    // 3. We initialize the check function we want to use.
    cf := subtest.NotDeepEqual("bar")

    // 4. We initialize a test function by passing the check function to the
    // value function's Test method.
    tf := vf.Test(cf)

    // 5. We run the test function as a sub-test.
    t.Run("v != bar", tf)
}

// TestFooFewerLines writes TestFooExplicit in fewer lines of code, but is
// exactly equivalent.
func TestFooFewerLines(t *testing.T) {
    v := "foo"
    t.Run("v != bar" subtest.Value(v).Test(subtest.NotDeepEqual("bar"))
}


// TestFooShortHand is an even shorter way of writing TestFooExplicit. It relies
// on short-hand methods defined on a value function for generating tests for
// our built-in check functions. In fact, all built-in checks has a short-hand
// format like this.
func TestFooShortHand(t *testing.T) {
    v := "foo"
    t.Run("v != bar", subtest.Value(v).NotDeepEqual("bar"))
}


// TestMultiple1 shows that a "subtest.ValueFunc" instance can be used to run
// several tests.
func TestMultiple1(t *testing.T) {
    v := "foo"
    vf := subtest.Value(v) // returns a reusable subtest.ValueFunc instance.

    t.Run("v != bar", vf.NotDeepEqual("bar"))
    t.Run("v == foo", vf.DeepEqual("foo"))
}

// TestMultiple2 shows that a "subtest.CheckFunc" can be used to test multiple
// values.
func TestMultiple2(t *testing.T) {
    notEmpty := subtest.NotDeepEqual("")

    for _, s := range []string{"foo", "bar"} {
        name := fmt.Sprintf(`%q != ""`)
        t.Run(name, subtest.Value(s).Test(notEmpty))
    }
}

// TestCustomCheckFunc shows that we can define custom check functions to
// initialize tests.
func TestCustomCheckFunc(t *testing.T) {
    v := "foo"

    t.Run("len(v) > 2", subtest.Value(len(v)).Test(intGt(2)))
}

func intGt(compare int) subtest.CheckFunc{
    return func(got interface{}) error {
        i, ok := got.(int)
        if !ok {
           return subtest.FailureGot("not an integer value", v)
        }
        if !(i > compare) {
            return subtest.FailureGot(fmt.Sprintf("value <= %v", compare), got)
        }
        return nil
    }
}

// TestSchema shows that we can validate the content of Go map types without
// requiring an exact match.
func TestSchema(t *testing.T) {
    var v = map[string]interface{}{"foo": "bar", "bar": 42}

    cf := subtest.Schema{
        Required: []interface{}{"foo", "bar"},
        Fields: subtest.Fields{
            "foo": subtest.DeepEqual("bar"),
            "bar": subtest.GreaterThan(41),
    }.ValidateMap()

    t.Run("v match schema", subtest.Value(v).Test(cf))
}

// TestJSONSchema shows that we can use schema validation and subjson middleware
// together to validate JSON fields. Package subjson middleware allows parsing
// []byte, string and json.RawMessage values to Go types before performing a
// check. See the package documentation to better understand which Go types are
// returned by each function.
func TestJSONSchema(t *testing.T) {
    const v = `{"foo": "bar", "bar": 42}`

    cf := subtest.Schema{
        Required: []interface{}{"foo", "bar"},
        Fields: subtest.Fields{
            "foo": subjson.OnString(subtest.DeepEqual("bar")),
            "bar": subjson.OnFloat64(subtest.GreaterThan(41)),
    }.ValidateMap()

    t.Run("v match schema", subtest.Value(v).Test(subjson.OnMap(cf)))
}
```

For further examples, see the `examples/` sub-directory:

- `examples/gwt`: Example of tests following the [Given-When-Then][gwt] naming convention.
- `examples/colorfmt`: Example of custom type formatting with colors via the [pp][pp] package.
- `examples/gojsonq`: Example of custom checks for JSON matching via the [gojsonq][gojsonq] package.
- `examples/jsondiff`: Example of custom checks for JSON comparison via the [jsondiff][jsondiff] package.

[pp]: https://github.com/k0kubun/pp
[gojsonq]: https://github.com/thedevsaddam/gojsonq
[jsondiff]: https://github.com/nsf/jsondiff

## Features

### Utilize the standard test runner

`subtest` initializes test functions intended for usage with the `Run` method on the `testing.T` type, and uses a plain output format by default. This means that tooling and IDE features built up around output from the standard test runner will work as expected.

### Check State-ful values

Values to check are wrapped in a producer function (`ValueFunc`). This means that multiple tests can be run against a state-ful value, such as an `io.Reader` by regenerating the value for each check.

## Allow check middleware

Generally, a sub-test performs of a single check (`CheckFunc`). These checks can be wrapped by middleware to facilitate processing or transformation of values before running the nested check. E.g. parse a byte array from JSON into a Go type.

### Zero-dependencies and un-opinionated

For now, the package is zero-dependencies. The important aspect of this to us, is that we don't _force_ potentially opinionated dependencies on the package user in order to offer features. Instead we aim to provide flexible interfaces so that it's easy to integrate your own preferred tools for everything from type formatting to advanced JSON or struct matching.

### Clear and plain output

The quicker a failed test can be understood, the quicker it can be fixed. `subtest`'s default failure formatting is inspired heavily by the short and simplistic style used for unit tests within the Go standard library.

Example output from our own tests:

```plain
--- FAIL: TestCheckNotReflectNil (0.00s)
    --- FAIL: TestCheckNotReflectNil/nil_struct_pointer (0.00s)
        /Users/smyrman/Code/subtest/check_test.go:127: error value is not matching target error
            got: subtest.Failure
                value is typed or untyped nil
                got: *subtest_test.T
                    {Foo:}
            want: subtest.Failure
                value is typed or untyped nil
                got: *subtest_test.T
                    nil
FAIL
FAIL    github.com/searis/subtest    0.006s
FAIL
Error: Tests failed.
```

Be aware that the default type formatter currently do not expand nested pointer values.

### Custom formatting

While we aim to make the default type formatting useful, it will also be somewhat limited due to our zero-dependency goal. Type formatting is also an area with potential for different opinions on what looks the most clear. For this reason we have made it easy to replace the default type formatter using libraries such as [go-spew][go-spew], [litter][litter], or [pp][pp] (with colors).

Example using `go-spew`:

```go
import (
    "github.com/davecgh/go-spew/spew"
    "github.com/searis/subtest"
)

func init() {
    subtest.SetTypeFormatter(spew.ConfigState{Indent: "\t"}.Sdump)
}
```

Example using `litter`:

```go
import (
    "github.com/sanity-io/litter"
    "github.com/searis/subtest"
)

func init() {
    subtest.SetTypeFormatter(litter.Options{}.Sdump)
}
```

Example using `pp` with conditional coloring:

```go
import (
    "golang.org/x/crypto/ssh/terminal"
    "github.com/k0kubun/pp"
    "github.com/searis/subtest"
)

func init() {
    subtest.SetTypeFormatter(pp.Sprint)

    colorEnv := strings.ToUpper(os.Getenv("GO_TEST_COLOR"))
    switch colorEnv {
    case "0", "FALSE":
      log.Println("explicitly disabling color output for test")
        pp.ColoringEnabled = false
    case "1", "TRUE":
        log.Println("explicitly enabling color output for test")
    default:
        if !terminal.IsTerminal(int(os.Stdout.Fd())) {
            log.Println("TTY not detected, disabling color output for test")
            pp.ColoringEnabled = false
        } else {
            log.Println("TTY detected, enabling color output for test")
        }
    }
}
```

 When it comes to prettifying the output of the test runner itself, there are separate tools for that. One such tool is [gotestsum][gotestsum], which wraps the Go test runner to provide alternate formatting.

[gotestsum]: https://github.com/gotestyourself/gotestsum
[litter]: https://github.com/sanity-io/litter
[go-spew]: https://github.com/davecgh/go-spew
[pp]: https://github.com/k0kubun/pp

### GWT-style sub-tests

Although `subtest` does not encourage nor limit sub-tests to follow any particular style, one of the main motivations for writing `subtest` was to allow easier [Given-When-Then][gwt] (GWT) style testing relying on pure Go sub-tests.

Below is an example of sub-tests with GWT-style naming:

[gwt]: https://martinfowler.com/bliki/GivenWhenThen.html

```go
package gwt_test

import (
    "testing"

    "github.com/searis/subtest"
    pkg "github.com/searis/subtest/examples/gwt"
)

func TestFoo(t *testing.T) {
    t.Run("given nothing is registered", func(t *testing.T) {
        reg := pkg.NewRegister()

        t.Run("when calling reg.Foo", func(t *testing.T) {
            result, err := reg.Foo()
            t.Run("then it should error",
                subtest.Value(err).ErrorIs(pkg.ErrNothingRegistered),
            )
            t.Run("then the result should hold a zero-value",
                subtest.Value(result).DeepEqual(""),
            )
        })
    })

    t.Run("given bar is registered", func(t *testing.T) {
        reg := pkg.NewRegister()
        reg.Register("bar")

        t.Run("when calling pkg.Foo", func(t *testing.T) {
            bar, err := reg.Foo()
            if !t.Run("then the result must not fail", subtest.Value(err).NoError()) {
                t.FailNow()
            }
            t.Run("then the result should be as expected",
                subtest.Value(bar).DeepEqual("foobar"),
            )
        })
    })
}

```
