# subtest

**Project status:** WIP

`subtest` is a minimalist Go test-utility package used to initializing small test functions for use with the Go sub-tests feature. You can read more about Go sub-tests [here][go-sub-test]. Go sub-tests allow you to nest your tests to provide a logical hierarchy, and _"focus"_ on specific sub-test by use of the test `-run` parameter.

For now, the package is zero-dependencies. The important aspect of this to us, is that we don't _force_ potentially opinionated dependencies on the package user in order to offer features. Instead we aim to provide flexible interfaces so that it's easy to integrate your own preferred tools for everything from type formatting to advanced JSON or struct matching.

[go-sub-test]: https://blog.golang.org/subtests

Here are some example tests:

```go
// Short-hand syntax in for built-in check functions.
func TestFoo1(t *testing.T) {
    v := "foo"
    t.Run("v != bar", subtest.Value(v).NotDeepEqual("bar"))
}

 // Equivalents to TestFoo1 using a longer syntax.
func TestFoo2(t *testing.T) {
    v := "foo"
    t.Run("v != bar" subtest.Value(v).Test(subtest.NotDeepEqual("bar"))
}

// A "subtest.ValueFunc" instance can be used to run several tests.
func TestMultiple1(t *testing.T) {
    v := "foo"
    vf := subtest.Value(v) // returns a subtest.ValueFunc instance.

    t.Run("v != bar", vf.NotDeepEqual("bar"))
    t.Run("v == foo", vf.DeepEqual("foo"))
}

// Similarly, a "subtest.CheckFunc" can be used to test multiple values.
func TestMultiple2(t *testing.T) {
    notEmpty := subtest.NotDeepEqual("")

    for _, s := range []string{"foo", "bar"} {
        name := fmt.Sprintf(`%q != ""`)
        t.Run(name, subtest.Value(s).Test(notEmpty))
    }
}

// Custom check functions can be used to initialize sub-tests.
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

// For validating JSON, you are perhaps not always interested in an exact match.
// If we parse it into a map, we can use subtest's reflect based schema
// validation feature to validate each field. We can also define our own
// subtest.ValueFunc function in order to delay the JSON Unmarshalling to happen
// inside each test.
func TestJSONSchema(t *testing.T) {
    const v = `{"foo": "bar", "bar": 42}`

    c := subtest.ValueFunc(func() (interface{}, error) {
        // Decode and return a new m instance for each test.
        var m map[string]interface{}
        err := json.Unmarshal([]byte(v), &m); err != nil {
            // Abort the test if v is not a JSON object.
            return nil, subtest.FailGot("not a valid JSON object:" + err.Error(), v)
        }
        return m, nil
    })

    t.Run("v match schema", c.ValidateMap(subtest.Schema{
        Required: []interface{}{"foo", "bar"},
        Fields: subtest.Fields{
            "foo": subtest.DeepEqual("bar"),
            "bar": floatGt(41),
        },
    }))
}

func floatGt(compare float64) subtest.CheckFunc {
    return func(got interface{}) error {
        v, ok := got.(float64)
        if !ok {
            return subtest.FailureGot("not a float64 value", v)
        }
        if !(v > compare) {
            return subtest.FailureGot(fmt.Sprintf("value <= %v", compare), got)
        }
        return nil
    }
}
```

## Clear and plain output

The quicker a failed test can be understood, the quicker it can be fixed. Therefore an important but easy to neglect detail of a test utility package, is how clearly it formats test failure. `subtest`'s failure formatting is inspired heavily by the short and simplistic style of Go standard library, extended only to provide clarity.

Unlike some other test utility library, this package does _not_ meddle with the output provided by the Go test runner. This sort of meddling could impact tools abilities to parse Go tests and treat them in a standardized way. E.g. VS Code has the ability to click on line numbers in the test output of failing tests to jump to the failing code, and CIs may want to parse the test results into a JUnit format for improved representation and stats. Enabling color  by default could also be problematic as some terminals, certainly IDEs and editors, just prints them as text. Therefore our _default_ error output is clear and colorless plain text.

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

## Custom formatting

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

 When it comes to prettifying the output of the test runner itself, there are separate tools for that as well. One such tool is [gotestsum][gotestsum], which wraps the Go test runner to provide alternate formatting.

[gotestsum]: https://github.com/gotestyourself/gotestsum
[litter]: https://github.com/sanity-io/litter
[go-spew]: https://github.com/davecgh/go-spew
[pp]: https://github.com/k0kubun/pp

## GWT-style sub-tests

Although `subtest` does not encourage nor limit sub-tests to follow any particular style, one of the main motivations for writing `subtest` was to allow easier [Given-When-Then][gwt] (GWT) style testing relying on pure Go sub-tests rather then a [Behavior-Driven-Development][bdd] framework with a custom test-runner or [Domain-Specific-Language][dsl]. Below is an example of sub-tests with GWT-style naming:

[gwt]: https://martinfowler.com/bliki/GivenWhenThen.html
[bdd]: https://en.wikipedia.org/wiki/Behavior-driven_development
[dsl]: https://en.wikipedia.org/wiki/Domain-specific_language

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

This example is accessible at `examples/gwt`.

## JSON matching via third-party

Zero-dependencies, means leaving some features out, but allowing users to plug them in. This example shows how you could integrate a custom tool, such as [gojsonq][gojsonq], into your tests. This example can be found in `examples/gojsonq`.

[gojsonq]: https://github.com/thedevsaddam/gojsonq
