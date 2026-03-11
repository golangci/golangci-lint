package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.yaml.in/yaml/v3"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/cwe"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/report/junit"
	"github.com/securego/gosec/v2/report/sonar"
)

func createIssueWithFileWhat(file, what string) *issue.Issue {
	issue := createIssue("i1", issue.GetCweByRule("G101"))
	issue.File = file
	issue.What = what
	return &issue
}

func createIssue(ruleID string, weakness *cwe.Weakness) issue.Issue {
	return issue.Issue{
		File:       "/home/src/project/test.go",
		Line:       "1",
		Col:        "1",
		RuleID:     ruleID,
		What:       "test",
		Confidence: issue.High,
		Severity:   issue.High,
		Code:       "1: testcode",
		Cwe:        weakness,
	}
}

func createReportInfo(rule string, weakness *cwe.Weakness) gosec.ReportInfo {
	newissue := createIssue(rule, weakness)
	metrics := gosec.Metrics{}
	return gosec.ReportInfo{
		Errors: map[string][]gosec.Error{},
		Issues: []*issue.Issue{
			&newissue,
		},
		Stats: &metrics,
	}
}

func stripString(str string) string {
	ret := strings.ReplaceAll(str, "\n", "")
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.ReplaceAll(ret, "\t", "")
	return ret
}

var _ = Describe("Formatter", func() {
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
				Issues: []*sonar.Issue{
					{
						EngineID: "gosec",
						RuleID:   "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "test.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						Type:          "VULNERABILITY",
						Severity:      "BLOCKER",
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
				Issues: []*sonar.Issue{
					{
						EngineID: "gosec",
						RuleID:   "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "subfolder/test.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						Type:          "VULNERABILITY",
						Severity:      "BLOCKER",
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
				Issues: []*sonar.Issue{
					{
						EngineID: "gosec",
						RuleID:   "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "test-project1.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						Type:          "VULNERABILITY",
						Severity:      "BLOCKER",
						EffortMinutes: sonar.EffortMinutes,
					},
					{
						EngineID: "gosec",
						RuleID:   "test",
						PrimaryLocation: &sonar.Location{
							Message:  "test",
							FilePath: "test-project2.go",
							TextRange: &sonar.TextRange{
								StartLine: 1,
								EndLine:   2,
							},
						},
						Type:          "VULNERABILITY",
						Severity:      "BLOCKER",
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

	Context("When using junit", func() {
		It("preserves order of issues", func() {
			issues := []*issue.Issue{createIssueWithFileWhat("i1", "1"), createIssueWithFileWhat("i2", "2"), createIssueWithFileWhat("i3", "1")}

			junitReport := junit.GenerateReport(&gosec.ReportInfo{Issues: issues})

			testSuite := junitReport.Testsuites[0]

			Expect(testSuite.Testcases[0].Name).To(Equal(issues[0].File))
			Expect(testSuite.Testcases[1].Name).To(Equal(issues[2].File))

			testSuite = junitReport.Testsuites[1]
			Expect(testSuite.Testcases[0].Name).To(Equal(issues[1].File))
		})
	})
	Context("When using different report formats", func() {
		grules := []string{
			"G101",
			"G102",
			"G103",
			"G104",
			"G106",
			"G107",
			"G109",
			"G110",
			"G111",
			"G112",
			"G201",
			"G202",
			"G203",
			"G204",
			"G301",
			"G302",
			"G303",
			"G304",
			"G305",
			"G401",
			"G402",
			"G403",
			"G404",
			"G405",
			"G406",
			"G407",
			"G501",
			"G502",
			"G503",
			"G504",
			"G505",
			"G506",
			"G507",
			"G601",
		}

		It("csv formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err := CreateReport(buf, "csv", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				pattern := "/home/src/project/test.go,1,test,HIGH,HIGH,1: testcode,CWE-%s\n"
				expect := fmt.Sprintf(pattern, cwe.ID)
				Expect(buf.String()).To(Equal(expect))
			}
		})
		It("xml formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{NumFiles: 0, NumLines: 0, NumNosec: 0, NumFound: 0}, errors).WithVersion("v2.7.0")
				err := CreateReport(buf, "xml", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				pattern := "Results:\n\n\n[/home/src/project/test.go:1] - %s (CWE-%s): test (Confidence: HIGH, Severity: HIGH)\n  > 1: testcode\n\nAutofix: \n\nSummary:\n  Gosec  : v2.7.0\n  Files  : 0\n  Lines  : 0\n  Nosec  : 0\n  Issues : 0\n\n"
				expect := fmt.Sprintf(pattern, rule, cwe.ID)
				Expect(buf.String()).To(Equal(expect))
			}
		})
		It("json formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				data := createReportInfo(rule, cwe)

				expect := new(bytes.Buffer)
				enc := json.NewEncoder(expect)
				err := enc.Encode(data)
				Expect(err).ShouldNot(HaveOccurred())
				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err = CreateReport(buf, "json", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				result := stripString(buf.String())
				expectation := stripString(expect.String())
				Expect(result).To(Equal(expectation))
			}
		})
		It("html formatted report should  contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				data := createReportInfo(rule, cwe)

				expect := new(bytes.Buffer)
				enc := json.NewEncoder(expect)
				err := enc.Encode(data)
				Expect(err).ShouldNot(HaveOccurred())
				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err = CreateReport(buf, "html", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				result := stripString(buf.String())
				expectation := stripString(expect.String())
				Expect(result).To(ContainSubstring(expectation))
			}
		})
		It("yaml formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				data := createReportInfo(rule, cwe)

				expect := new(bytes.Buffer)
				enc := yaml.NewEncoder(expect)
				err := enc.Encode(data)
				Expect(err).ShouldNot(HaveOccurred())
				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err = CreateReport(buf, "yaml", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				result := stripString(buf.String())
				expectation := stripString(expect.String())
				Expect(result).To(ContainSubstring(expectation))
			}
		})
		It("junit-xml formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				data := createReportInfo(rule, cwe)

				expect := new(bytes.Buffer)
				enc := yaml.NewEncoder(expect)
				err := enc.Encode(data)
				Expect(err).ShouldNot(HaveOccurred())
				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err = CreateReport(buf, "junit-xml", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				expectation := stripString(fmt.Sprintf("[/home/src/project/test.go:1] - test (Confidence: 2, Severity: 2, CWE: %s)", cwe.ID))
				result := stripString(buf.String())
				Expect(result).To(ContainSubstring(expectation))
			}
		})
		It("text formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				data := createReportInfo(rule, cwe)

				expect := new(bytes.Buffer)
				enc := yaml.NewEncoder(expect)
				err := enc.Encode(data)
				Expect(err).ShouldNot(HaveOccurred())
				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err = CreateReport(buf, "text", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				expectation := stripString(fmt.Sprintf("[/home/src/project/test.go:1] - %s (CWE-%s): test (Confidence: HIGH, Severity: HIGH)", rule, cwe.ID))
				result := stripString(buf.String())
				Expect(result).To(ContainSubstring(expectation))
			}
		})
		It("sonarqube formatted report shouldn't contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}
				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err := CreateReport(buf, "sonarqube", false, []string{"/home/src/project"}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())

				result := stripString(buf.String())

				expect := new(bytes.Buffer)
				enc := json.NewEncoder(expect)
				err = enc.Encode(cwe)
				Expect(err).ShouldNot(HaveOccurred())

				expectation := stripString(expect.String())
				Expect(result).ShouldNot(ContainSubstring(expectation))
			}
		})
		It("golint formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors)
				err := CreateReport(buf, "golint", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())
				pattern := "/home/src/project/test.go:1:1: [CWE-%s] test (Rule:%s, Severity:HIGH, Confidence:HIGH)\n"
				expect := fmt.Sprintf(pattern, cwe.ID, rule)
				Expect(buf.String()).To(Equal(expect))
			}
		})
		It("sarif formatted report should contain the CWE mapping", func() {
			for _, rule := range grules {
				cwe := issue.GetCweByRule(rule)
				newissue := createIssue(rule, cwe)
				errors := map[string][]gosec.Error{}

				buf := new(bytes.Buffer)
				reportInfo := gosec.NewReportInfo([]*issue.Issue{&newissue}, &gosec.Metrics{}, errors).WithVersion("v2.7.0")
				err := CreateReport(buf, "sarif", false, []string{}, reportInfo)
				Expect(err).ShouldNot(HaveOccurred())

				result := stripString(buf.String())

				ruleIDPattern := "\"id\":\"%s\""
				expectedRule := fmt.Sprintf(ruleIDPattern, rule)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(result).To(ContainSubstring(expectedRule))

				cweURIPattern := "\"helpUri\":\"https://cwe.mitre.org/data/definitions/%s.html\""
				expectedCweURI := fmt.Sprintf(cweURIPattern, cwe.ID)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(result).To(ContainSubstring(expectedCweURI))

				cweIDPattern := "\"id\":\"%s\""
				expectedCweID := fmt.Sprintf(cweIDPattern, cwe.ID)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(result).To(ContainSubstring(expectedCweID))
			}
		})
	})

	Context("When converting suppressed issues", func() {
		ruleID := "G101"
		cwe := issue.GetCweByRule(ruleID)
		suppressions := []issue.SuppressionInfo{
			{
				Kind:          "kind",
				Justification: "justification",
			},
		}
		suppressedIssue := createIssue(ruleID, cwe)
		suppressedIssue.WithSuppressions(suppressions)

		It("text formatted report should contain the suppressed issues", func() {
			errors := map[string][]gosec.Error{}
			reportInfo := gosec.NewReportInfo([]*issue.Issue{&suppressedIssue}, &gosec.Metrics{}, errors)

			buf := new(bytes.Buffer)
			err := CreateReport(buf, "text", false, []string{}, reportInfo)
			Expect(err).ShouldNot(HaveOccurred())

			result := stripString(buf.String())
			Expect(result).To(ContainSubstring("Results:Summary"))
		})

		It("sarif formatted report should contain the suppressed issues", func() {
			errors := map[string][]gosec.Error{}
			reportInfo := gosec.NewReportInfo([]*issue.Issue{&suppressedIssue}, &gosec.Metrics{}, errors)

			buf := new(bytes.Buffer)
			err := CreateReport(buf, "sarif", false, []string{}, reportInfo)
			Expect(err).ShouldNot(HaveOccurred())

			result := stripString(buf.String())
			Expect(result).To(ContainSubstring(`"results":[{`))
		})

		It("json formatted report should contain the suppressed issues", func() {
			errors := map[string][]gosec.Error{}
			reportInfo := gosec.NewReportInfo([]*issue.Issue{&suppressedIssue}, &gosec.Metrics{}, errors)

			buf := new(bytes.Buffer)
			err := CreateReport(buf, "json", false, []string{}, reportInfo)
			Expect(err).ShouldNot(HaveOccurred())

			result := stripString(buf.String())
			Expect(result).To(ContainSubstring(`"Issues":[{`))
		})
	})
})
