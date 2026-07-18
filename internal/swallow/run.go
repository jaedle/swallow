package swallow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

const (
	tagStdout = "out|"
	tagStderr = "err|"
)

func Run(argv []string) int {
	agent := os.Getenv("CLAUDECODE") == "1"

	logPath, logFile, err := createLog(argv[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 1
	}
	defer func() { _ = logFile.Close() }()

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin = os.Stdin
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 1
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 1
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 126
	}

	var tee, teeErr io.Writer
	if !agent {
		tee, teeErr = os.Stdout, os.Stderr
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go capture(&wg, stdout, logFile, tagStdout, tee)
	go capture(&wg, stderr, logFile, tagStderr, teeErr)
	wg.Wait()

	code := exitCode(cmd.Wait())
	_ = logFile.Close()

	if agent {
		if code == 0 {
			fmt.Printf("everything went fine (log: %s)\n", logPath)
		} else {
			replay(logPath)
			fmt.Fprintf(os.Stderr, "swallow: command failed with exit code %d (log: %s)\n", code, logPath)
		}
	}

	return code
}

// replay streams the log back, restoring every line to its original stream.
func replay(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return
	}
	defer func() { _ = file.Close() }()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			target := os.Stdout
			if rest, ok := strings.CutPrefix(line, tagStderr); ok {
				target, line = os.Stderr, rest
			} else if rest, ok := strings.CutPrefix(line, tagStdout); ok {
				line = rest
			}
			_, _ = target.Write([]byte(line))
		}
		if err != nil {
			return
		}
	}
}

// capture streams one child stream line-wise into the shared log, prefixing
// every line with its stream tag. Memory usage is bounded by the longest
// single line; the output as a whole is never held in memory. Each log write
// is one complete tagged line so concurrent appends of both streams never
// break the line framing.
func capture(wg *sync.WaitGroup, stream io.Reader, log io.Writer, tag string, tee io.Writer) {
	defer wg.Done()

	reader := bufio.NewReader(stream)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if tee != nil {
				_, _ = tee.Write([]byte(line))
			}
			tagged := tag + line
			if line[len(line)-1] != '\n' {
				tagged += "\n"
			}
			_, _ = log.Write([]byte(tagged))
		}
		if err != nil {
			return
		}
	}
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}

	var exit *exec.ExitError
	if errors.As(err, &exit) {
		status := exit.Sys().(syscall.WaitStatus)
		if status.Signaled() {
			return 128 + int(status.Signal())
		}
		return status.ExitStatus()
	}

	fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
	return 1
}
