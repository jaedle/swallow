package swallow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

const (
	tagStdout = "out|"
	tagStderr = "err|"

	// replayLimit caps the lines replayed after a failure: a huge failing
	// output would flood the caller's context, defeating swallow's purpose.
	// The last lines win — that is where the error usually is; the full log
	// stays available via the read hint.
	replayLimit = 100
)

func Run(argv []string) int {
	agent := os.Getenv("CLAUDECODE") == "1"

	if dir, err := swallowDir(); err == nil {
		prune(dir)
	}

	if _, err := exec.LookPath(argv[0]); err != nil {
		fmt.Fprintf(os.Stderr, "swallow: command not found: %s\n", argv[0])
		return 127
	}

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

	if agent {
		// The command name only — echoed arguments could leak shell-expanded
		// secrets into the caller's context, see ADR 0009. Printed only once
		// the command has started, so every start line gets a done line.
		fmt.Printf("swallow: running %s, swallowing output\n", filepath.Base(argv[0]))
	}

	var tee, teeErr io.Writer
	if !agent {
		tee, teeErr = os.Stdout, os.Stderr
	}

	signals := make(chan os.Signal, 8)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go forward(signals, cmd)

	var wg sync.WaitGroup
	wg.Add(2)
	go capture(&wg, stdout, logFile, tagStdout, tee)
	go capture(&wg, stderr, logFile, tagStderr, teeErr)
	wg.Wait()

	code := exitCode(cmd.Wait())
	signal.Stop(signals)
	close(signals)
	_ = logFile.Close()

	if agent {
		hint := fmt.Sprintf("read logs: `swallow --read %s`", filepath.Base(logPath))
		lines := countLogLines(logPath)
		if code == 0 {
			fmt.Printf("swallow: done, exit code 0, %d log lines, %s\n", lines, hint)
		} else {
			if lines > replayLimit {
				fmt.Fprintf(os.Stderr, "swallow: done, exit code %d, last %d of %d lines:\n", code, replayLimit, lines)
				replay(logPath, lines-replayLimit)
			} else {
				fmt.Fprintf(os.Stderr, "swallow: done, exit code %d, full output (%d lines):\n", code, lines)
				replay(logPath, 0)
			}
			fmt.Fprintf(os.Stderr, "swallow: end of output, exit code %d, %s\n", code, hint)
		}
	}

	return code
}

// forward relays received signals to the child. The child stays in swallow's
// process group, so an interactive Ctrl-C may reach it twice — an accepted
// trade-off, see ADR 0005.
func forward(signals <-chan os.Signal, cmd *exec.Cmd) {
	for s := range signals {
		_ = cmd.Process.Signal(s)
	}
}

// countLogLines counts the tagged lines of a log by streaming over it;
// memory usage stays bounded by the longest single line.
func countLogLines(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer func() { _ = file.Close() }()

	count := 0
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			count++
		}
		if err != nil {
			return count
		}
	}
}

// replay streams the log back, restoring every line to its original stream.
// The first skip lines are dropped, implementing the replay cap.
func replay(path string, skip int) {
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
			if skip > 0 {
				skip--
			} else {
				target := os.Stdout
				if rest, ok := strings.CutPrefix(line, tagStderr); ok {
					target, line = os.Stderr, rest
				} else if rest, ok := strings.CutPrefix(line, tagStdout); ok {
					line = rest
				}
				_, _ = target.Write([]byte(line))
			}
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
