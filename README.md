# subtest

[![Go Report Card](https://goreportcard.com/badge/github.com/searis/subtest)](https://goreportcard.com/report/github.com/searis/subtest)
[![GoDev](https://img.shields.io/static/v1?label=go.dev&message=reference&color=blue)](https://pkg.go.dev/github.com/searis/subtest)

**subtest** is a minimalist Go test-utility package used to initializing small test functions for use with the Go sub-tests feature. You can read more about Go sub-tests [here][go-sub-test].

The sub-package `subjson` defines middleware for parsing values from JSON before performing checks.

[go-sub-test]: https://blog.golang.org/subtests

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
    t.Run("Given foo is 42", func(t *testing.T) {
        const foo = 42
        t.Run("When dividing by 6", func (t *testing.T) {
            v := float64(foo) / 6
            t.Run("Then the result should be 7", func(t *testing.T) {
                if expect := float64(7); v != expect {
                    t.Fatalf("\n got: %d\nwant: %d", v, expect)
                }
            })
        })
        t.Run("When dividing by 9", func (t *testing.T) {
            v := float64(foo) / 9
            t.Run("Then the result should be greater than 4" func(t *testing.T) {
                if expect := float64(4); v > expect {
                    t.Fatalf("\n got: %d\nwant > %d", v, expect)
                }
            })
            t.Run("Then the result should be less than 5" func(t *testing.T) {
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
    t.Run("Given foo is 42", func(t *testing.T) {
        const foo = 42
        t.Run("When dividing by 6", func (t *testing.T) {
            vf := subtest.Value(float64(foo) / 6)
            t.Run("Then the result should be 7", vf.NumericEqual(7))
        })
        t.Run("When dividing by 9", func (t *testing.T) {
            vf := subtest.Value(float64(foo) / 6)
            t.Run("Then the result should be greater than 4", vf.GreaterThan(4))
            t.Run("Then the result should be less than 5", vf.LessThan(5))
        })
    })
}
```

## Usage

### Illustrative example

Building and running a subtest is generally composed of six steps. Normally you do not do each steps as explicit as described below, but to illustrate the general flow, we have spelled this out.

```go
 // TestFooExplicit shows the different steps of building a subtest.
func TestFooExplicit(t *testing.T) {
    // 1. We declare the value we want to check; usually the result of an
    // operation or action.
    v := "foo"

    // 2. We initialize a value function for the value that we want to check.
    // There are several different initializers we can call to get a value
    // function; this is the simplest one.
    vf := subtest.Value(v)

    // 3. We initialize the check we want to use. A check is anything that
    // implements the Check interface.
    c := subtest.NumericEquals(3)

    // 4. We can optionally wrap our check with middleware.
    c = subtest.OnLen(c)

    // 5. We initialize a test function by passing the check to the value
    // function's Test method.
    tf := vf.Test(c)

    // 6. We run the test function as a sub-test.
    t.Run("len(v) == 3", tf)
}
```

If we where going to do this every time, we would grow weary. Therefore there is several short-hand methods defined on the Check and ValueFunc instances that makes things easier. The least verbose variant we can write of the test above is as follows:

```go
func TestFoo(t *testing.T) {
    v := "foo"

    t.Run("len(v) == 3", subtest.Len(v).NumericEquals(3))
}

```

### JSON Schema validation

It is possible to validate more than just equality with subtest. The `subtest.Schema` type allows advanced validation of any Go map type, and in the future, perhaps also for structs. From the `subjson` package we can use `ValueFunc` initializers, `Check` implementations and check middleware to decode JSON from `string`, `[]byte` and `json.RawMessage` values. Combining these two mechanisms we can do advanced validation of JSON content.

```go
func TestJSONMap(t *testing.T) {
    v := `{"foo": "bar", "bar": "foobar", "baz": ["foo", "bar", "baz"]}`

    expect := subtest.Schema{
        Fields: subtest.Fields{
            "foo": subjson.DecodesTo("bar")
            "bar": subjson.OnLen(subtest.AllOf{
                subtest.GreaterThan(3),
                subtest.LessThan(8),
            }),
            "baz": subjson.OnSlice(subtest.AllOf{
                subtest.OnLen(subtest.DeepEqual(3)),
                subtest.OnIndex(0, subjson.DecodesTo("foo"),
                subtest.OnIndex(1, subtest.MatchPattern(`"^b??$"`), // regex match against raw JSON
                subtest.OnIndex(2, subtest.DeepEqual(json.RawMessage(`"baz"`)), // raw JSON equals
            }),
        },
    }

    t.Run("match expectations", subjson.Map(v).Test(expect))
}
```

### Required checks

This is perhaps not commonly known, but the `t.Run` function actually return `false` if there is a failure. Or to be more accurate:

> Run reports whether f succeeded (or at least did not fail before calling t.Parallel).

Because subtest checks do not call `t.Parallel`, this can be utilized to stop test-execution if a "required" sub-test fails.

```go
func TestFoo(t *testing.T) {
    v, err := foo()

    if !t.Run("err == nil", subtest.Value(err).NoError()) {
        // Abort further tests if failed.
        t.FailNow()
    }
    // Never run when err != nil.
    t.Run("v == foo", subtest.Value(v).DeepEqual("foo"))
}

func foo() (string, error) {
    return "", errors.New("failed")
}
```

### Extendability

The subtest library itself is currently zero-dependencies. The important aspect of this is that we do not force opinionated dependencies on the user. However, it's also written to be relatively easy to extend.

For specialized use cases and customization, see the `examples/` sub-directory:

- `examples/gwt`: Example of tests following the [Given-When-Then][gwt] naming convention.
- `examples/colorfmt`: Example of custom type formatting with colors via the [pp][pp] package.
- `examples/gojsonq`: Example of custom checks for JSON matching via the [gojsonq][gojsonq] package.
- `examples/jsondiff`: Example of custom checks for JSON comparison via the [jsondiff][jsondiff] package.

[pp]: https://github.com/k0kubun/pp
[gojsonq]: https://github.com/thedevsaddam/gojsonq
[jsondiff]: https://github.com/nsf/jsondiff

## Features

Some key features of **subtest** is described below.

### Utilize the standard test runner

**subtest** initializes test functions intended for usage with the `Run` method on the `testing.T` type, and uses a plain output format by default. This means that tooling and IDE features built up around output from the standard test runner will work as expected.

### Check State-ful values

Values to check are wrapped in a value function (`ValueFunc`). By setting up your own value function, you can easily run several tests against state-ful types, such as an io.Reader, where each check starts with
a clean slate.

### Check middleware

Generally, a sub-test performs of a single check (`CheckFunc`). These checks can be wrapped by middleware to facilitate processing or transformation of values before running nested checks. E.g. parse a byte array from JSON into a Go type, or extract the length of an array.

### Plain output

The quicker a failed test can be understood, the quicker it can be fixed. `subtest`'s default failure formatting is inspired by the short and simplistic style used for unit tests within the Go standard library. We have extended this syntax only so that we can more easily format the expected type and value.

Example output from an `exaples/gwt`:

```plain
--- FAIL: TestFoo (0.00s)
    --- FAIL: TestFoo/Given_nothing_is_registered (0.00s)
        --- FAIL: TestFoo/Given_nothing_is_registered/When_calling_reg.Foo (0.00s)
            --- FAIL: TestFoo/Given_nothing_is_registered/When_calling_reg.Foo/Then_the_result_should_hold_a_zero-value (0.00s)
                pkg_test.go:19: not deep equal
                    got: string
                        "oops"
                    want: string
                        ""
FAIL
FAIL	github.com/searis/subtest/examples/gwt	0.057s
FAIL
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
