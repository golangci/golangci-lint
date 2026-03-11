package gosec_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGosec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gosec Suite")
}
