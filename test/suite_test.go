package swallow_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var binary string

func TestSwallow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swallow Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := gexec.Build("github.com/jaedle/swallow/cmd/swallow")
	Expect(err).NotTo(HaveOccurred())
	return []byte(path)
}, func(path []byte) {
	binary = string(path)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})
