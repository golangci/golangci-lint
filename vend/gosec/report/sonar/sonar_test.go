package sonar_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/report/sonar"
)

var _ = Describe("Sonar Formatter", func() {
	BeforeEach(func() {
	})
	Context("when converting to Sonarqube issues", func() {
		It("it should parse the report info", func() {
			data := &gosec.ReportInfo{
				Errors: map[string][]gosec.Error{},
				Issues: []*issue.Issue{
					{
						Severity:   2,
						Confidence: 0,
						RuleID:     "test",
						What:       "test",
						File:       "/home/src/project/test.go",
						Code:       "",
						Line:       "1-2",
					},
				},
				Stats: &gosec.Metrics{
					NumFiles: 0,
					NumLines: 0,
					NumNosec: 0,
					NumFound: 0,
				},
			}
			want := &sonar.Report{
				Rules: []*sonar.Rule{
					{
						ID:                 "test",
						Name:               "test",
						Description:        "test",
						EngineID:           "gosec",
						CleanCodeAttribute: "TRUSTWORTHY",
						Impacts: []*sonar.Impact{
							{
								SoftwareQuality: "SECURITY",
								Severity:        "HIGH",
							},
						},
					},
				},
				Issues: []*sonar.Issue{
					{
						RuleID: "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "test.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						EffortMinutes: sonar.EffortMinutes,
					},
				},
			}

			rootPath := "/home/src/project"

			issues, err := sonar.GenerateReport([]string{rootPath}, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*issues).To(Equal(*want))
		})

		It("it should enrich rules with metadata when available", func() {
			data := &gosec.ReportInfo{
				Errors: map[string][]gosec.Error{},
				Issues: []*issue.Issue{
					{
						Severity:   issue.Medium,
						Confidence: 0,
						RuleID:     "G101",
						What:       "Potential hardcoded credentials",
						File:       "/home/src/project/test.go",
						Code:       "",
						Line:       "5",
					},
				},
				Stats: &gosec.Metrics{},
			}

			want := &sonar.Report{
				Rules: []*sonar.Rule{
					{
						ID:                 "G101",
						Name:               "G101",
						Description:        "Look for hardcoded credentials",
						EngineID:           "gosec",
						CleanCodeAttribute: "TRUSTWORTHY",
						Impacts: []*sonar.Impact{
							{
								SoftwareQuality: "SECURITY",
								Severity:        "MEDIUM",
							},
						},
					},
				},
				Issues: []*sonar.Issue{
					{
						RuleID: "G101",
						PrimaryLocation: &sonar.Location{
							Message:  "Potential hardcoded credentials",
							FilePath: "test.go",
							TextRange: &sonar.TextRange{
								StartLine: 5,
								EndLine:   5,
							},
						},
						EffortMinutes: sonar.EffortMinutes,
					},
				},
			}

			rootPath := "/home/src/project"

			issues, err := sonar.GenerateReport([]string{rootPath}, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*issues).To(Equal(*want))
		})

		It("it should parse the report info with files in subfolders", func() {
			data := &gosec.ReportInfo{
				Errors: map[string][]gosec.Error{},
				Issues: []*issue.Issue{
					{
						Severity:   2,
						Confidence: 0,
						RuleID:     "test",
						What:       "test",
						File:       "/home/src/project/subfolder/test.go",
						Code:       "",
						Line:       "1-2",
					},
				},
				Stats: &gosec.Metrics{
					NumFiles: 0,
					NumLines: 0,
					NumNosec: 0,
					NumFound: 0,
				},
			}
			want := &sonar.Report{
				Rules: []*sonar.Rule{
					{
						ID:                 "test",
						Name:               "test",
						Description:        "test",
						EngineID:           "gosec",
						CleanCodeAttribute: "TRUSTWORTHY",
						Impacts: []*sonar.Impact{
							{
								SoftwareQuality: "SECURITY",
								Severity:        "HIGH",
							},
						},
					},
				},
				Issues: []*sonar.Issue{
					{
						RuleID: "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "subfolder/test.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						EffortMinutes: sonar.EffortMinutes,
					},
				},
			}

			rootPath := "/home/src/project"

			issues, err := sonar.GenerateReport([]string{rootPath}, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*issues).To(Equal(*want))
		})
		It("it should not parse the report info for files from other projects", func() {
			data := &gosec.ReportInfo{
				Errors: map[string][]gosec.Error{},
				Issues: []*issue.Issue{
					{
						Severity:   2,
						Confidence: 0,
						RuleID:     "test",
						What:       "test",
						File:       "/home/src/project1/test.go",
						Code:       "",
						Line:       "1-2",
					},
				},
				Stats: &gosec.Metrics{
					NumFiles: 0,
					NumLines: 0,
					NumNosec: 0,
					NumFound: 0,
				},
			}
			want := &sonar.Report{
				Rules:  []*sonar.Rule{},
				Issues: []*sonar.Issue{},
			}

			rootPath := "/home/src/project2"

			issues, err := sonar.GenerateReport([]string{rootPath}, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*issues).To(Equal(*want))
		})

		It("it should parse the report info for multiple projects", func() {
			data := &gosec.ReportInfo{
				Errors: map[string][]gosec.Error{},
				Issues: []*issue.Issue{
					{
						Severity:   2,
						Confidence: 0,
						RuleID:     "test",
						What:       "test",
						File:       "/home/src/project1/test-project1.go",
						Code:       "",
						Line:       "1-2",
					},
					{
						Severity:   2,
						Confidence: 0,
						RuleID:     "test",
						What:       "test",
						File:       "/home/src/project2/test-project2.go",
						Code:       "",
						Line:       "1-2",
					},
				},
				Stats: &gosec.Metrics{
					NumFiles: 0,
					NumLines: 0,
					NumNosec: 0,
					NumFound: 0,
				},
			}
			want := &sonar.Report{
				Rules: []*sonar.Rule{
					{
						ID:                 "test",
						Name:               "test",
						Description:        "test",
						EngineID:           "gosec",
						CleanCodeAttribute: "TRUSTWORTHY",
						Impacts: []*sonar.Impact{
							{
								SoftwareQuality: "SECURITY",
								Severity:        "HIGH",
							},
						},
					},
				},
				Issues: []*sonar.Issue{
					{
						RuleID: "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "test-project1.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						EffortMinutes: sonar.EffortMinutes,
					},
					{
						RuleID: "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "test-project2.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						EffortMinutes: sonar.EffortMinutes,
					},
				},
			}

			rootPaths := []string{"/home/src/project1", "/home/src/project2"}

			issues, err := sonar.GenerateReport(rootPaths, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*issues).To(Equal(*want))
		})
	})
})
