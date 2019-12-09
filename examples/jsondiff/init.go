package jsondiff

import (
	"log"
	"os"
	"strings"

	"github.com/nsf/jsondiff"
	"golang.org/x/crypto/ssh/terminal"
)

var cfg struct {
	jsonDiff jsondiff.Options
}

func init() {
	if checkTerminalColor() {
		cfg.jsonDiff = jsondiff.DefaultConsoleOptions()
	} else {
		cfg.jsonDiff = jsondiff.Options{
			Normal:  jsondiff.Tag{Begin: "", End: ""},
			Added:   jsondiff.Tag{Begin: "+[", End: "]"},
			Removed: jsondiff.Tag{Begin: "-[", End: "]"},
			Changed: jsondiff.Tag{Begin: "~[", End: "]"},
			Indent:  "    ",
		}
	}

}

func checkTerminalColor() bool {
	colorEnv := strings.ToUpper(os.Getenv("GO_TEST_COLOR"))
	switch colorEnv {
	case "1", "TRUE":
		log.Println("explicitly enabling color output for test")
		return true
	case "0", "FALSE":
		log.Println("explicitly disabling color output for test")
		return false
	}

	// Environment variable not set; attempt to detect color support
	// automatically.
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		log.Println("TTY detected, enabling color output for test")
		return true
	}
	log.Println("TTY not detected, disabling color output for test")
	return false
}
