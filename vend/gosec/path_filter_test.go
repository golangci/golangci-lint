package gosec_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
)

var _ = Describe("PathExclusionFilter", func() {
	Describe("NewPathExclusionFilter", func() {
		Context("with valid rules", func() {
			It("should create a filter with single rule", func() {
				rules := []gosec.PathExcludeRule{
					{Path: "cmd/.*", Rules: []string{"G204", "G304"}},
				}
				filter, err := gosec.NewPathExclusionFilter(rules)
				Expect(err).NotTo(HaveOccurred())
				Expect(filter).NotTo(BeNil())
			})

			It("should create a filter with multiple rules", func() {
				rules := []gosec.PathExcludeRule{
					{Path: "cmd/.*", Rules: []string{"G204"}},
					{Path: "test/.*", Rules: []string{"G101"}},
					{Path: "scripts/.*", Rules: []string{"*"}},
				}
				filter, err := gosec.NewPathExclusionFilter(rules)
				Expect(err).NotTo(HaveOccurred())
				Expect(filter).NotTo(BeNil())
			})

			It("should handle empty rules slice", func() {
				filter, err := gosec.NewPathExclusionFilter(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(filter).NotTo(BeNil())
			})
		})

		Context("with invalid rules", func() {
			It("should reject empty path", func() {
				rules := []gosec.PathExcludeRule{
					{Path: "", Rules: []string{"G204"}},
				}
				_, err := gosec.NewPathExclusionFilter(rules)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("path cannot be empty"))
			})

			It("should reject invalid regex", func() {
				rules := []gosec.PathExcludeRule{
					{Path: "[invalid(regex", Rules: []string{"G204"}},
				}
				_, err := gosec.NewPathExclusionFilter(rules)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid path regex"))
			})
		})
	})

	Describe("ShouldExclude", func() {
		var filter *gosec.PathExclusionFilter

		Context("with specific rule exclusions", func() {
			BeforeEach(func() {
				rules := []gosec.PathExcludeRule{
					{Path: "cmd/.*", Rules: []string{"G204", "G304"}},
					{Path: "internal/testutil/.*", Rules: []string{"G101"}},
				}
				var err error
				filter, err = gosec.NewPathExclusionFilter(rules)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should exclude matching path and rule", func() {
				Expect(filter.ShouldExclude("cmd/mytool/main.go", "G204")).To(BeTrue())
				Expect(filter.ShouldExclude("cmd/another/file.go", "G304")).To(BeTrue())
			})

			It("should not exclude matching path with non-matching rule", func() {
				Expect(filter.ShouldExclude("cmd/mytool/main.go", "G101")).To(BeFalse())
				Expect(filter.ShouldExclude("cmd/mytool/main.go", "G401")).To(BeFalse())
			})

			It("should not exclude non-matching path", func() {
				Expect(filter.ShouldExclude("pkg/server/main.go", "G204")).To(BeFalse())
				Expect(filter.ShouldExclude("internal/api/handler.go", "G304")).To(BeFalse())
			})

			It("should handle nested paths correctly", func() {
				Expect(filter.ShouldExclude("internal/testutil/helper.go", "G101")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/testutil/sub/file.go", "G101")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/other/file.go", "G101")).To(BeFalse())
			})
		})

		Context("with wildcard rule exclusion", func() {
			BeforeEach(func() {
				rules := []gosec.PathExcludeRule{
					{Path: "scripts/.*", Rules: []string{"*"}},
					{Path: "vendor/.*", Rules: []string{"*"}},
				}
				var err error
				filter, err = gosec.NewPathExclusionFilter(rules)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should exclude any rule for matching path", func() {
				Expect(filter.ShouldExclude("scripts/build.go", "G101")).To(BeTrue())
				Expect(filter.ShouldExclude("scripts/build.go", "G204")).To(BeTrue())
				Expect(filter.ShouldExclude("scripts/build.go", "G304")).To(BeTrue())
				Expect(filter.ShouldExclude("vendor/lib/file.go", "G401")).To(BeTrue())
			})

			It("should not exclude non-matching paths", func() {
				Expect(filter.ShouldExclude("cmd/main.go", "G101")).To(BeFalse())
			})
		})

		Context("with Windows-style paths", func() {
			BeforeEach(func() {
				rules := []gosec.PathExcludeRule{
					{Path: "cmd/.*", Rules: []string{"G204"}},
				}
				var err error
				filter, err = gosec.NewPathExclusionFilter(rules)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should normalize backslashes to forward slashes", func() {
				Expect(filter.ShouldExclude("cmd\\mytool\\main.go", "G204")).To(BeTrue())
				Expect(filter.ShouldExclude("cmd\\nested\\deep\\file.go", "G204")).To(BeTrue())
			})
		})

		Context("with nil or empty filter", func() {
			It("should not exclude anything with nil filter", func() {
				var nilFilter *gosec.PathExclusionFilter
				Expect(nilFilter.ShouldExclude("any/path.go", "G101")).To(BeFalse())
			})

			It("should not exclude anything with empty rules", func() {
				filter, _ := gosec.NewPathExclusionFilter(nil)
				Expect(filter.ShouldExclude("any/path.go", "G101")).To(BeFalse())
			})
		})

		Context("with complex regex patterns", func() {
			BeforeEach(func() {
				rules := []gosec.PathExcludeRule{
					{Path: `.*_test\.go$`, Rules: []string{"G101"}},
					{Path: `^(cmd|tools)/`, Rules: []string{"G204"}},
					{Path: `internal/(mock|fake|stub)s?/`, Rules: []string{"*"}},
				}
				var err error
				filter, err = gosec.NewPathExclusionFilter(rules)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should match test files", func() {
				Expect(filter.ShouldExclude("pkg/auth/auth_test.go", "G101")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/handler_test.go", "G101")).To(BeTrue())
				Expect(filter.ShouldExclude("pkg/auth/auth.go", "G101")).To(BeFalse())
			})

			It("should match cmd or tools prefix", func() {
				Expect(filter.ShouldExclude("cmd/server/main.go", "G204")).To(BeTrue())
				Expect(filter.ShouldExclude("tools/generator/gen.go", "G204")).To(BeTrue())
				Expect(filter.ShouldExclude("pkg/cmd/helper.go", "G204")).To(BeFalse())
			})

			It("should match mock/fake/stub directories", func() {
				Expect(filter.ShouldExclude("internal/mocks/service.go", "G401")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/mock/client.go", "G304")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/fakes/repo.go", "G101")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/stub/handler.go", "G204")).To(BeTrue())
				Expect(filter.ShouldExclude("internal/real/service.go", "G401")).To(BeFalse())
			})
		})
	})

	Describe("FilterIssues", func() {
		var filter *gosec.PathExclusionFilter

		BeforeEach(func() {
			rules := []gosec.PathExcludeRule{
				{Path: "cmd/.*", Rules: []string{"G204", "G304"}},
				{Path: "test/.*", Rules: []string{"*"}},
			}
			var err error
			filter, err = gosec.NewPathExclusionFilter(rules)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter matching issues", func() {
			issues := []*issue.Issue{
				{File: "cmd/main.go", RuleID: "G204"},
				{File: "cmd/config.go", RuleID: "G304"},
				{File: "pkg/server.go", RuleID: "G204"},
				{File: "test/helper.go", RuleID: "G101"},
			}

			filtered, excluded := filter.FilterIssues(issues)
			Expect(excluded).To(Equal(3))
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].File).To(Equal("pkg/server.go"))
		})

		It("should handle empty issues slice", func() {
			filtered, excluded := filter.FilterIssues(nil)
			Expect(excluded).To(Equal(0))
			Expect(filtered).To(BeNil())
		})

		It("should preserve issue order", func() {
			issues := []*issue.Issue{
				{File: "a.go", RuleID: "G101"},
				{File: "b.go", RuleID: "G102"},
				{File: "c.go", RuleID: "G103"},
			}

			filtered, excluded := filter.FilterIssues(issues)
			Expect(excluded).To(Equal(0))
			Expect(filtered).To(HaveLen(3))
			Expect(filtered[0].File).To(Equal("a.go"))
			Expect(filtered[1].File).To(Equal("b.go"))
			Expect(filtered[2].File).To(Equal("c.go"))
		})
	})

	Describe("ParseCLIExcludeRules", func() {
		Context("with valid input", func() {
			It("should parse single rule", func() {
				rules, err := gosec.ParseCLIExcludeRules("cmd/.*:G204,G304")
				Expect(err).NotTo(HaveOccurred())
				Expect(rules).To(HaveLen(1))
				Expect(rules[0].Path).To(Equal("cmd/.*"))
				Expect(rules[0].Rules).To(ConsistOf("G204", "G304"))
			})

			It("should parse multiple rules separated by semicolon", func() {
				rules, err := gosec.ParseCLIExcludeRules("cmd/.*:G204;test/.*:G101,G102")
				Expect(err).NotTo(HaveOccurred())
				Expect(rules).To(HaveLen(2))
				Expect(rules[0].Path).To(Equal("cmd/.*"))
				Expect(rules[0].Rules).To(ConsistOf("G204"))
				Expect(rules[1].Path).To(Equal("test/.*"))
				Expect(rules[1].Rules).To(ConsistOf("G101", "G102"))
			})

			It("should handle wildcard rule", func() {
				rules, err := gosec.ParseCLIExcludeRules("scripts/.*:*")
				Expect(err).NotTo(HaveOccurred())
				Expect(rules).To(HaveLen(1))
				Expect(rules[0].Rules).To(ConsistOf("*"))
			})

			It("should handle empty input", func() {
				rules, err := gosec.ParseCLIExcludeRules("")
				Expect(err).NotTo(HaveOccurred())
				Expect(rules).To(BeNil())
			})

			It("should trim whitespace", func() {
				rules, err := gosec.ParseCLIExcludeRules("  cmd/.* : G204 , G304  ;  test/.* : G101  ")
				Expect(err).NotTo(HaveOccurred())
				Expect(rules).To(HaveLen(2))
				Expect(rules[0].Path).To(Equal("cmd/.*"))
				Expect(rules[0].Rules).To(ConsistOf("G204", "G304"))
			})
		})

		Context("with invalid input", func() {
			It("should reject missing colon separator", func() {
				_, err := gosec.ParseCLIExcludeRules("cmd/.*G204")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("missing ':'"))
			})

			It("should reject empty path", func() {
				_, err := gosec.ParseCLIExcludeRules(":G204")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("path pattern cannot be empty"))
			})

			It("should reject empty rules", func() {
				_, err := gosec.ParseCLIExcludeRules("cmd/.*:")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("rules list cannot be empty"))
			})
		})
	})

	Describe("MergeExcludeRules", func() {
		It("should merge CLI and config rules with CLI first", func() {
			cliRules := []gosec.PathExcludeRule{
				{Path: "cli/.*", Rules: []string{"G204"}},
			}
			configRules := []gosec.PathExcludeRule{
				{Path: "config/.*", Rules: []string{"G304"}},
			}

			merged := gosec.MergeExcludeRules(configRules, cliRules)
			Expect(merged).To(HaveLen(2))
			Expect(merged[0].Path).To(Equal("cli/.*")) // CLI first
			Expect(merged[1].Path).To(Equal("config/.*"))
		})

		It("should handle empty CLI rules", func() {
			configRules := []gosec.PathExcludeRule{
				{Path: "config/.*", Rules: []string{"G304"}},
			}

			merged := gosec.MergeExcludeRules(configRules, nil)
			Expect(merged).To(Equal(configRules))
		})

		It("should handle empty config rules", func() {
			cliRules := []gosec.PathExcludeRule{
				{Path: "cli/.*", Rules: []string{"G204"}},
			}

			merged := gosec.MergeExcludeRules(nil, cliRules)
			Expect(merged).To(Equal(cliRules))
		})
	})
})

// Standard Go tests for those who prefer table-driven tests
func TestShouldExclude(t *testing.T) {
	tests := []struct {
		name     string
		rules    []gosec.PathExcludeRule
		filePath string
		ruleID   string
		want     bool
	}{
		{
			name: "exact match",
			rules: []gosec.PathExcludeRule{
				{Path: "cmd/.*", Rules: []string{"G204"}},
			},
			filePath: "cmd/main.go",
			ruleID:   "G204",
			want:     true,
		},
		{
			name: "no match - wrong rule",
			rules: []gosec.PathExcludeRule{
				{Path: "cmd/.*", Rules: []string{"G204"}},
			},
			filePath: "cmd/main.go",
			ruleID:   "G304",
			want:     false,
		},
		{
			name: "no match - wrong path",
			rules: []gosec.PathExcludeRule{
				{Path: "cmd/.*", Rules: []string{"G204"}},
			},
			filePath: "pkg/main.go",
			ruleID:   "G204",
			want:     false,
		},
		{
			name: "wildcard excludes all rules",
			rules: []gosec.PathExcludeRule{
				{Path: "scripts/.*", Rules: []string{"*"}},
			},
			filePath: "scripts/build.go",
			ruleID:   "G999",
			want:     true,
		},
		{
			name: "multiple rules in single exclusion",
			rules: []gosec.PathExcludeRule{
				{Path: "cmd/.*", Rules: []string{"G204", "G304", "G404"}},
			},
			filePath: "cmd/tool/main.go",
			ruleID:   "G304",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := gosec.NewPathExclusionFilter(tt.rules)
			if err != nil {
				t.Fatalf("NewPathExclusionFilter() error = %v", err)
			}

			got := filter.ShouldExclude(tt.filePath, tt.ruleID)
			if got != tt.want {
				t.Errorf("ShouldExclude(%q, %q) = %v, want %v",
					tt.filePath, tt.ruleID, got, tt.want)
			}
		})
	}
}
