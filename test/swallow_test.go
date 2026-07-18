package swallow_test

import (
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

var _ = Describe("process control", func() {
	It("propagates the exit code of the wrapped command", func() {
		session := run(runOptions{
			swallowDir: GinkgoT().TempDir(),
			args:       []string{"sh", "-c", "exit 5"},
		})

		wait(session, 5)
	})
})
