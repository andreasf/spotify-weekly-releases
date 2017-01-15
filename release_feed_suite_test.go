package release_feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestReleaseFeed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReleaseFeed Suite")
}
