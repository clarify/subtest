package colorfmt_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/searis/subtest"
	"golang.org/x/crypto/ssh/terminal"
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
func TestFoo(t *testing.T) {
	type T struct {
		Foo *string
		Bar int
	}

	s1 := "bar"
	v := &T{Foo: &s1}

	s2 := "bar"
	t.Run("T{Foo:&s1}==T{Foo:&s2}", subtest.Value(v).DeepEqual(&T{Foo: &s2}))
}
