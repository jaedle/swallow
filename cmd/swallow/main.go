package main

import (
	"fmt"
	"os"

	"github.com/jaedle/swallow/internal/swallow"
)

var version = "dev"

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(args) > 0 && args[0] == "--read" {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: swallow --read <log-file>")
			os.Exit(2)
		}
		os.Exit(swallow.Read(args[1]))
	}

	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: swallow [--] <command> [args...]")
		os.Exit(2)
	}

	os.Exit(swallow.Run(args))
}
