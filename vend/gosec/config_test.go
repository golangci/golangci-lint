package gosec_test

import (
	"bytes"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
)

var _ = Describe("Configuration", func() {
	var configuration gosec.Config
	BeforeEach(func() {
		configuration = gosec.NewConfig()
	})

	Context("when loading from disk", func() {
		It("should be possible to load configuration from a file", func() {
			json := `{"G101": {}}`
			buffer := bytes.NewBufferString(json)
			nread, err := configuration.ReadFrom(buffer)
			Expect(nread).Should(Equal(int64(len(json))))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error if configuration file is invalid", func() {
			var err error
			invalidBuffer := bytes.NewBuffer([]byte{0xc0, 0xff, 0xee})
			_, err = configuration.ReadFrom(invalidBuffer)
			Expect(err).Should(HaveOccurred())

			emptyBuffer := bytes.NewBuffer([]byte{})
			_, err = configuration.ReadFrom(emptyBuffer)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when saving to disk", func() {
		It("should be possible to save an empty configuration to file", func() {
			expected := `{"global":{}}`
			buffer := bytes.NewBuffer([]byte{})
			nbytes, err := configuration.WriteTo(buffer)
			Expect(int(nbytes)).Should(Equal(len(expected)))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).Should(Equal(expected))
		})

		It("should be possible to save configuration to file", func() {
			configuration.Set("G101", map[string]string{
				"mode": "strict",
			})

			buffer := bytes.NewBuffer([]byte{})
			nbytes, err := configuration.WriteTo(buffer)
			Expect(int(nbytes)).ShouldNot(BeZero())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).Should(Equal(`{"G101":{"mode":"strict"},"global":{}}`))
		})
	})

	Context("when configuring rules", func() {
		It("should be possible to get configuration for a rule", func() {
			settings := map[string]string{
				"ciphers": "AES256-GCM",
			}
			configuration.Set("G101", settings)

			retrieved, err := configuration.Get("G101")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(retrieved).Should(HaveKeyWithValue("ciphers", "AES256-GCM"))
			Expect(retrieved).ShouldNot(HaveKey("foobar"))
		})
	})

	Context("when using global configuration options", func() {
		It("should have a default global section", func() {
			settings, err := configuration.Get("global")
			Expect(err).ShouldNot(HaveOccurred())
			expectedType := make(map[gosec.GlobalOption]string)
			Expect(settings).Should(BeAssignableToTypeOf(expectedType))
		})

		It("should save global settings to correct section", func() {
			configuration.SetGlobal(gosec.Nosec, "enabled")
			settings, err := configuration.Get("global")
			Expect(err).ShouldNot(HaveOccurred())
			if globals, ok := settings.(map[gosec.GlobalOption]string); ok {
				Expect(globals["nosec"]).Should(MatchRegexp("enabled"))
			} else {
				Fail("globals are not defined as map")
			}

			setValue, err := configuration.GetGlobal(gosec.Nosec)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(setValue).Should(MatchRegexp("enabled"))
		})

		It("should find global settings which are enabled", func() {
			configuration.SetGlobal(gosec.Nosec, "enabled")
			enabled, err := configuration.IsGlobalEnabled(gosec.Nosec)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enabled).Should(BeTrue())
		})

		It("should parse the global settings of type string from file", func() {
			config := `
			{
				"global": {
					"nosec": "enabled"
				}
			}`
			cfg := gosec.NewConfig()
			_, err := cfg.ReadFrom(strings.NewReader(config))
			Expect(err).ShouldNot(HaveOccurred())

			value, err := cfg.GetGlobal(gosec.Nosec)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal("enabled"))
		})
		It("should parse the global settings of other types from file", func() {
			config := `
			{
				"global": {
					"nosec": true
				}
			}`
			cfg := gosec.NewConfig()
			_, err := cfg.ReadFrom(strings.NewReader(config))
			Expect(err).ShouldNot(HaveOccurred())

			value, err := cfg.GetGlobal(gosec.Nosec)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal("true"))
		})
	})
})
