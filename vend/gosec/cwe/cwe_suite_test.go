package cwe_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCwe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cwe Suite")
}
