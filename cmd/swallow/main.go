package main

import (
	"fmt"
	"os"

	"github.com/jaedle/swallow/internal/swallow"
)

var version = "dev"

const usage = `usage:
  swallow [--] <command> [args...]  run a command, swallowing its output into a log
  swallow --read <log-file>         print a captured log of the current origin
  swallow --version                 print the version
  swallow --help                    print this help
`

func main() {
	args := os.Args[1:]

	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Print(usage)
		os.Exit(0)
	}

	if len(args) > 0 && args[0] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(args) > 0 && args[0] == "--read" {
		if len(args) != 2 {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(2)
		}
		os.Exit(swallow.Read(args[1]))
	}

	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	os.Exit(swallow.Run(args))
}
