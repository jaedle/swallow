package swallow_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// processTimeout bounds every wait on the swallow binary; generous for slow CI workers.
const processTimeout = 10 * time.Second

type runOptions struct {
	agent      bool
	env        []string
	swallowDir string
	home       string
	dir        string
	stdin      string
	args       []string
}

// run starts the swallow binary with an environment built from scratch so the
// host environment (e.g. a developer's CLAUDECODE=1) never leaks into specs.
func run(opts runOptions) *gexec.Session {
	GinkgoHelper()

	cmd := exec.Command(binary, opts.args...)
	cmd.Dir = opts.dir
	cmd.Env = []string{"PATH=" + os.Getenv("PATH")}
	if opts.home != "" {
		cmd.Env = append(cmd.Env, "HOME="+opts.home)
	}
	if opts.swallowDir != "" {
		cmd.Env = append(cmd.Env, "SWALLOW_DIR="+opts.swallowDir)
	}
	if opts.agent {
		cmd.Env = append(cmd.Env, "CLAUDECODE=1")
	}
	cmd.Env = append(cmd.Env, opts.env...)
	if opts.stdin != "" {
		cmd.Stdin = strings.NewReader(opts.stdin)
	}

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

func wait(session *gexec.Session, exitCode int) {
	GinkgoHelper()

	Eventually(session, processTimeout).Should(gexec.Exit(exitCode))
}

func findLogs(swallowDir string) []string {
	GinkgoHelper()

	logs, err := filepath.Glob(filepath.Join(swallowDir, "*", "*.log"))
	Expect(err).NotTo(HaveOccurred())
	return logs
}

func singleLog(swallowDir string) string {
	GinkgoHelper()

	logs := findLogs(swallowDir)
	Expect(logs).To(HaveLen(1))
	return logs[0]
}

func logContent(path string) string {
	GinkgoHelper()

	content, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}

// slugOf mirrors the specified origin slug: non-alphanumeric runs become a
// single dash, trimmed at both ends.
func slugOf(path string) string {
	return strings.Trim(regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(path, "-"), "-")
}

var _ = Describe("agent mode", func() {
	It("suppresses output and reports success with the log location", func() {
		swallowDir := GinkgoT().TempDir()

		session := run(runOptions{
			agent:      true,
			swallowDir: swallowDir,
			args:       []string{"sh", "-c", "echo to-stdout; echo to-stderr 1>&2"},
		})
		wait(session, 0)

		Expect(string(session.Out.Contents())).To(MatchRegexp(`^everything went fine \(log: .*\.log\)\n$`))
		Expect(session.Err.Contents()).To(BeEmpty())
		log := logContent(singleLog(swallowDir))
		Expect(log).To(ContainSubstring("out|to-stdout\n"))
		Expect(log).To(ContainSubstring("err|to-stderr\n"))
	})

	It("replays the full output split by stream on failure and propagates the exit code", func() {
		swallowDir := GinkgoT().TempDir()

		session := run(runOptions{
			agent:      true,
			swallowDir: swallowDir,
			args:       []string{"sh", "-c", "echo first-out; echo to-stderr 1>&2; echo second-out; exit 3"},
		})
		wait(session, 3)

		stdout := string(session.Out.Contents())
		Expect(stdout).To(ContainSubstring("first-out\n"))
		Expect(stdout).To(ContainSubstring("second-out\n"))
		Expect(stdout).NotTo(ContainSubstring("to-stderr"))
		stderr := string(session.Err.Contents())
		Expect(stderr).To(ContainSubstring("to-stderr\n"))
		Expect(stderr).To(MatchRegexp(`command failed with exit code 3 \(log: .*\.log\)`))
	})

	It("treats only CLAUDECODE=1 as an agentic caller", func() {
		for _, value := range []string{"CLAUDECODE=", "CLAUDECODE=true"} {
			session := run(runOptions{
				swallowDir: GinkgoT().TempDir(),
				env:        []string{value},
				args:       []string{"echo", "visible"},
			})
			wait(session, 0)

			Expect(session.Out).To(gbytes.Say("visible"))
			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("everything went fine"))
		}
	})
})

var _ = Describe("human mode", func() {
	It("tees stdout and stderr live and captures both in the log", func() {
		swallowDir := GinkgoT().TempDir()

		session := run(runOptions{
			swallowDir: swallowDir,
			args:       []string{"sh", "-c", "echo to-stdout; echo to-stderr 1>&2"},
		})
		wait(session, 0)

		Expect(session.Out).To(gbytes.Say("to-stdout"))
		Expect(string(session.Out.Contents())).NotTo(ContainSubstring("to-stderr"))
		Expect(session.Err).To(gbytes.Say("to-stderr"))
		log := logContent(singleLog(swallowDir))
		Expect(log).To(ContainSubstring("out|to-stdout\n"))
		Expect(log).To(ContainSubstring("err|to-stderr\n"))
	})
})

var _ = Describe("logging", func() {
	It("names the log after origin, timestamp, command and a unique suffix", func() {
		swallowDir := GinkgoT().TempDir()
		origin := GinkgoT().TempDir()

		session := run(runOptions{
			swallowDir: swallowDir,
			dir:        origin,
			args:       []string{"cat", "/dev/null"},
		})
		wait(session, 0)

		log := singleLog(swallowDir)
		Expect(filepath.Dir(log)).To(Equal(filepath.Join(swallowDir, slugOf(origin))))
		Expect(filepath.Base(log)).To(MatchRegexp(`^\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}-cat-[0-9a-f]{6}\.log$`))
	})

	It("creates a distinct log for every run", func() {
		swallowDir := GinkgoT().TempDir()

		wait(run(runOptions{swallowDir: swallowDir, args: []string{"true"}}), 0)
		wait(run(runOptions{swallowDir: swallowDir, args: []string{"true"}}), 0)

		Expect(findLogs(swallowDir)).To(HaveLen(2))
	})

	It("defaults to .swallow in the home directory", func() {
		home := GinkgoT().TempDir()

		session := run(runOptions{
			home: home,
			args: []string{"true"},
		})
		wait(session, 0)

		Expect(findLogs(filepath.Join(home, ".swallow"))).To(HaveLen(1))
	})
})

var _ = Describe("process control", func() {
	It("propagates the exit code of the wrapped command", func() {
		session := run(runOptions{
			swallowDir: GinkgoT().TempDir(),
			args:       []string{"sh", "-c", "exit 5"},
		})

		wait(session, 5)
	})

	It("forwards stdin to the wrapped command", func() {
		swallowDir := GinkgoT().TempDir()

		session := run(runOptions{
			swallowDir: swallowDir,
			stdin:      "hello\n",
			args:       []string{"cat"},
		})
		wait(session, 0)

		Expect(session.Out).To(gbytes.Say("hello"))
		Expect(logContent(singleLog(swallowDir))).To(ContainSubstring("out|hello\n"))
	})

	It("forwards termination signals to the wrapped command", func() {
		swallowDir := GinkgoT().TempDir()

		session := run(runOptions{
			swallowDir: swallowDir,
			args:       []string{"sh", "-c", `trap 'kill $! 2>/dev/null; exit 42' TERM; echo ready; sleep 30 & wait $!`},
		})
		Eventually(func() string {
			logs := findLogs(swallowDir)
			if len(logs) != 1 {
				return ""
			}
			return logContent(logs[0])
		}, processTimeout).Should(ContainSubstring("ready"))
		session.Signal(syscall.SIGTERM)

		wait(session, 42)
	})

	It("fails with 127 when the command does not exist", func() {
		swallowDir := GinkgoT().TempDir()

		session := run(runOptions{
			swallowDir: swallowDir,
			args:       []string{"definitely-not-here-xyz"},
		})
		wait(session, 127)

		Expect(session.Err).To(gbytes.Say("command not found: definitely-not-here-xyz"))
		Expect(findLogs(swallowDir)).To(BeEmpty())
	})

	It("accepts a -- separator before the command", func() {
		session := run(runOptions{
			swallowDir: GinkgoT().TempDir(),
			args:       []string{"--", "echo", "hi"},
		})
		wait(session, 0)

		Expect(session.Out).To(gbytes.Say("hi"))
	})

	It("prints usage and fails without a command", func() {
		session := run(runOptions{
			swallowDir: GinkgoT().TempDir(),
		})
		wait(session, 2)

		Expect(session.Err).To(gbytes.Say(`usage: swallow \[--\] <command>`))
	})
})
