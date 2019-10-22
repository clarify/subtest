# subtest

**Project status:** WIP

`subtest` is a minimalist Go test-utility package used to initializing small test functions. These micro-tests are perfect for use with the Go sub-tests feature [introduced in Go 1.7][go-sub-test]. Go sub-tests allow you to nest your tests to provide a logical hierarchy, and _"focus"_ on specific sub-test by use of the test `-run` parameter.

For now, the package is zero-dependencies. The important aspect of this to us, is that we don't _force_ potentially opinionated dependencies on the package user in order to offer features. Instead we aim to provide flexible interfaces so that it's easy to integrate your own preferred tools for everything from type formatting to advanced JSON or struct matching.

[go-sub-test]: https://blog.golang.org/subtests

## Formatting the output

The quicker a failed test can be understood, the quicker it can be fixed. Therefore an important but easy to neglect detail of a test utility package, is how clearly it formats test failure. `subtest`'s failure formatting is inspired heavily by the short and simplistic style of Go standard library, extended only to provide clarity.

Unlike some other test utility library, this package does _not_ meddle with the output provided by the Go test runner. This sort of meddling could impact tools abilities to parse Go tests and treat them in a standardized way. E.g. VS Code has the ability to click on line numbers in the test output to jump to the failing code, and CIs may want to parse the test results into a JUnit format for improved representation and stats. That doesn't mean that the standard Go test output is good for all cases. Luckily tools such as [gotestsum][gotestsum] exists to handle these problems separate from the test implementations.

For formatting complex types, we rely on `fmt.Printf("%#v", v)` by default. For many cases this is good enough, but it also have some limitations such as replacing pointer references. However, the field of pretty-printing Go types involves preferences and opinions on everything from syntax, coloring to indentation. We therefore allow you decide for yourself by calling `subtest.SetTestFormatter`. Good alternatives for type formatting include [go-spew][go-spew], [litter][litter], and [pp][pp] (with colors).

Using `pp` and conditional coloring:

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

Using `litter`:

```go
import (
    "github.com/sanity-io/litter"
    "github.com/searis/subtest"
)

func init() {
    subtest.SetTypeFormatter(litter.Options{Separator: " "}.Sdump)
}
```

Using `go-spew`:

```go
import (
    "github.com/davecgh/go-spew/spew"
    "github.com/searis/subtest"
)

func init() {
    subtest.SetTypeFormatter(spew.ConfigState{Indent: "\t"}.Sdump)
}
```

[gotestsum]: https://github.com/gotestyourself/gotestsum
[litter]: https://github.com/sanity-io/litter
[go-spew]: https://github.com/davecgh/go-spew
[pp]: https://github.com/k0kubun/pp

## GWT-style sub-tests

Although `subtest` does not encourage nor limit sub-tests to follow any particular style, one of the main motivations for writing `subtest` was to allow easier [Given-When-Then]][gwt] style testing relying on pure Go sub-tests rather then a [Behavior-Driven-Development][bdd] framework with a custom test-runner or DSL.

[gwt]: https://martinfowler.com/bliki/GivenWhenThen.html
[bdd]: https://en.wikipedia.org/wiki/Behavior-driven_development

Example usage with GWT-style sub-test naming:

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
                subtest.Check(err).ErrorIs(pkg.ErrNothingRegistered),
            )

            t.Run("then the result should hold a zero-value",
                subtest.Check(result).DeepEqual(""),
            )
        })
    })

    t.Run("given bar is registered", func(t *testing.T) {
        reg := pkg.NewRegister()
        reg.Register("bar")

        t.Run("when calling pkg.Foo", func(t *testing.T) {
            bar, err := reg.Foo()
            subtest.NoError().Test(err)(t) // short-hand to abort on failure.

            // An equivalent way to abort on failure.
            if !t.Run("then the result must not fail", subtest.Check(err).NoError()) {
                t.FailNow()
            }

            t.Run("then the result should be as expected",
                subtest.Check(bar).DeepEqual("foobar"),
            )
        })
    })
}
```

## JSON matching via third-party

Zero-dependencies, means leaving some features out, but allowing users to plug them in. When it comes to validating or comparing JSON, various packages exists to allow easy parsing and querying of JSON data. One such tool is [gojsonq][gojsonq].

Below is an example JSON validation via gojson:

```go
package gojsonq_test

import (
  "errors"
  "testing"

  "github.com/searis/subtest"
  "gopkg.in/thedevsaddam/gojsonq.v2"
)

func TestJSON1(t *testing.T) {
    // j allows queries against v, but there is a catch; for each query j must
    // be reset to allow for the next query. This test shows how to handle this
    // manually prior to each test.

    const v = `{"name": "foo", "type": "bar"}`

    j := gojsonq.New().JSONString(v)
    t.Run(".", subtest.Check(j).Test(validSchema))
    j.Reset()
    t.Run(".name", subtest.Check(j.Find("name")).DeepEqual("foo"))
    j.Reset()
    t.Run(".type", subtest.Check(j.Find("type")).DeepEqual("bar"))
}

func validSchema(got interface{}) error {
    j := got.(*gojsonq.JSONQ)

    var errs subtest.Errors

    if cnt, expect := j.Count(), 2; cnt != expect {
      errs = append(errs, subtest.FailureExpect("incorrect number of keys", cnt, expect))

    }

    j.Reset()
    if _, ok := j.Find("name").(string); !ok {
       err s = append(errs, errors.New(".name not a string"))
    }

    j.Reset()
    if _, ok := j.Find("type").(string); !ok {
       errs = append(errs, errors.New(".type not a string"))
    }

    if len(errs) > 0 {
        return errs
    }
    return nil
}

func TestJSON2(t *testing.T) {
    // This test demonstrates that a check don't have to be defined against a
    // static value. When it's useful, a new instance can be generated for each
    // test.

    const v = `{"name": "foo", "type": "bar"}`

    c := subtest.C(func() interface{} {
      return gojsonq.New().JSONString(v)
    })

    t.Run(".", c.Test(validSchema))
    t.Run(".name", c.Test(jsonPathEqual("name", "foo")))
    t.Run(".type", c.Test(jsonPathEqual("type", "bar")))
}

func jsonPathEqual(path string, expect interface{}) subtest.CheckFunc {
    check := subtest.DeepEqual(expect)

    return func(got interface{}) error {
        j := got.(*gojsonq.JSONQ)
        return check(j.Find(path))
    }
}

```

[gojsonq]: https://github.com/thedevsaddam/gojsonq
