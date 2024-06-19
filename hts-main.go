//@hts
//do not delete previous comment

package gohtslib

import (
	"os"
	"fmt"
)

var scripts map[string]func() = map[string]func(){}
var postLogContent []string = []string{}

func Script(name string, code func()) func() {
    scripts[name] = code
    return code
}

func Main() {
	println("outer main")
	argv := ParseArgv[struct{script string; postLog bool}]()
	if !argv.present["script"] {
		panic(fmt.Sprintf("-script flag is not specified\n"))
		os.Exit(1)
	}
    code, found := scripts[argv.flags.script]
    if found {
        code()
    } else {
		fmt.Printf("Unknown script \"%s\"\n", argv.flags.script)
		os.Exit(1)
    }
    if argv.flags.postLog {
        for _, msg := range postLogContent {
            fmt.Printf("\x1b[38;5;180m%s\x1b[39m\n", msg)
        }
    }
}
