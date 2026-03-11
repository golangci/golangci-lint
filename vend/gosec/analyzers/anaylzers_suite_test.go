package analyzers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAnalyzers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Analyzers Suite")
}
