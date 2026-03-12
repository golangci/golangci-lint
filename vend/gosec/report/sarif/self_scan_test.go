package sarif_test

import (
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/report/sarif"
	"github.com/securego/gosec/v2/rules"
)

var _ = Describe("Sarif Self Scan", func() {
	It("produces locally valid sarif without null relationships when scanning gosec source", func() {
		repoRoot := currentRepoRoot()

		config := gosec.NewConfig()
		logger := log.New(io.Discard, "", 0)
		analyzer := gosec.NewAnalyzer(config, false, true, false, 4, logger)

		ruleList := rules.Generate(false, rules.NewRuleFilter(false, "G401"))
		analyzer.LoadRules(ruleList.RulesInfo())

		excludedDirs := gosec.ExcludedDirsRegExp([]string{"vendor", ".git"})
		packagePaths, err := gosec.PackagePaths(filepath.Join(repoRoot, "..."), excludedDirs)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(packagePaths).ShouldNot(BeEmpty())

		err = analyzer.Process(nil, packagePaths...)
		Expect(err).ShouldNot(HaveOccurred())

		issues, metrics, errors := analyzer.Report()
		Expect(issues).ShouldNot(BeEmpty())
		reportInfo := gosec.NewReportInfo(issues, metrics, errors).WithVersion("test")

		sarifReport, err := sarif.GenerateReport([]string{repoRoot}, reportInfo)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(validateSarifSchema(sarifReport)).To(Succeed())

		encoded, err := json.Marshal(sarifReport)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(encoded).NotTo(ContainSubstring(`"relationships":[null]`))
		Expect(encoded).NotTo(ContainSubstring(`"relationships": [null]`))
	})
})

func currentRepoRoot() string {
	programCounter, currentFile, line, ok := runtime.Caller(0)
	if !ok || programCounter == 0 || line <= 0 {
		return filepath.Clean(".")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
}
