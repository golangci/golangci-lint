// (c) Copyright 2024 Mercedes-Benz Tech Innovation GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gosec_test

import (
	"errors"
	"fmt"
	"go/build"
	"log"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/analyzers"
	"github.com/securego/gosec/v2/rules"
	"github.com/securego/gosec/v2/testutils"
)

var _ = Describe("Analyzer", func() {
	var (
		analyzer  *gosec.Analyzer
		logger    *log.Logger
		buildTags []string
		tests     bool
	)
	BeforeEach(func() {
		logger, _ = testutils.NewLogger()
		analyzer = gosec.NewAnalyzer(nil, tests, false, false, 1, logger)
	})

	Context("when processing a package", func() {
		It("should not report an error if the package contains no Go files", func() {
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			dir := GinkgoT().TempDir()
			err := analyzer.Process(buildTags, dir)
			Expect(err).ShouldNot(HaveOccurred())
			_, _, errors := analyzer.Report()
			Expect(errors).To(BeEmpty())
		})

		It("should report an error if the package fails to build", func() {
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("wonky.go", `func main(){ println("forgot the package")}`)
			err := pkg.Build()
			Expect(err).Should(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			_, _, errors := analyzer.Report()
			Expect(errors).To(HaveLen(1))
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(1))
			}
		})

		It("should be able to analyze multiple Go files", func() {
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package main
				func main(){
					bar()
				}`)
			pkg.AddFile("bar.go", `
				package main
				func bar(){
					println("package has two files!")
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			_, metrics, _ := analyzer.Report()
			Expect(metrics.NumFiles).To(Equal(2))
		})

		It("should be able to analyze multiple Go files concurrently", func() {
			customAnalyzer := gosec.NewAnalyzer(nil, true, true, false, 32, logger)
			customAnalyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package main
				func main(){
					bar()
				}`)
			pkg.AddFile("bar.go", `
				package main
				func bar(){
					println("package has two files!")
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			_, metrics, _ := customAnalyzer.Report()
			Expect(metrics.NumFiles).To(Equal(2))
		})

		It("should be able to analyze multiple Go packages", func() {
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg1 := testutils.NewTestPackage()
			pkg2 := testutils.NewTestPackage()
			defer pkg1.Close()
			defer pkg2.Close()
			pkg1.AddFile("foo.go", `
				package main
				func main(){
				}`)
			pkg2.AddFile("bar.go", `
				package main
				func bar(){
				}`)
			err := pkg1.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = pkg2.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg1.Path, pkg2.Path)
			Expect(err).ShouldNot(HaveOccurred())
			_, metrics, _ := analyzer.Report()
			Expect(metrics.NumFiles).To(Equal(2))
		})

		It("should find errors when nosec is not in use", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("md5.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			controlIssues, _, _ := analyzer.Report()
			Expect(controlIssues).Should(HaveLen(sample.Errors))
		})

		It("should find errors when nosec is not in use", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("cipher.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			controlIssues, _, _ := analyzer.Report()
			Expect(controlIssues).Should(HaveLen(sample.Errors))
		})

		It("should find errors when nosec is not in use", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("md4.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			controlIssues, _, _ := analyzer.Report()
			Expect(controlIssues).Should(HaveLen(sample.Errors))
		})

		It("should report Go build errors and invalid files", func() {
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package main
				func main()
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			_, _, errors := analyzer.Report()
			foundErr := false
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(1))
				match, err := regexp.MatchString(ferr[0].Err, `expected declaration, found '}'`)
				if !match || err != nil {
					continue
				}
				foundErr = true
				Expect(ferr[0].Line).To(Equal(4))
				Expect(ferr[0].Column).To(Equal(5))
				Expect(ferr[0].Err).Should(MatchRegexp(`expected declaration, found '}'`))
			}
			Expect(foundErr).To(BeTrue())
		})

		It("should not report errors when a nosec line comment is present", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a disable directive is present", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a nosec line comment is present", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a disable directive is present", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a nosec line comment is present", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a disable directive is present", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a nosec block comment is present", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() /* #nosec */", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a nosec block comment is present", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) /* #nosec */", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a nosec block comment is present", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() /* #nosec */", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for the correct rule", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec G401", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for the correct rule", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable G401", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for the correct rule", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec G405", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for the correct rule", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable G405", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for the correct rule", func() {
			// Rule for MD4 deprecated weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec G406", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for the correct rule", func() {
			// Rule for MD4 deprecated weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable G406", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a nosec  block and line comment are present", func() {
			sample := testutils.SampleCodeG101[23]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G101")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("g101.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})
		It("should not report errors when only a nosec  block is present", func() {
			sample := testutils.SampleCodeG101[24]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G101")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("g101.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})
		It("should not report errors when a single line nosec  is present on a multi-line issue", func() {
			sample := testutils.SampleCodeG112[3]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G112")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("g112.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when a disable directive block and line comment are present", func() {
			sample := testutils.SampleCodeG101[26]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G101")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("g101.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})
		It("should not report errors when only a disable directive block is present", func() {
			sample := testutils.SampleCodeG101[27]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G101")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("g101.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})
		It("should not report errors when a single line nosec  is present on a multi-line issue", func() {
			sample := testutils.SampleCodeG112[4]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G112")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("g112.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should report errors when an exclude comment is present for a different rule", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec G301", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when an exclude comment is present for a different rule", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable G301", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when an exclude comment is present for a different rule", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec G301", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when an exclude comment is present for a different rule", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable G301", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when an exclude comment is present for a different rule", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec G301", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when an exclude comment is present for a different rule", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable G301", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should not report errors when an exclude comment is present for multiple rules, including the correct rule", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec G301 G401", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for multiple rules, including the correct rule", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable G301 G401", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for multiple rules, including the correct rule", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec G301 G405", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for multiple rules, including the correct rule", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable G301 G405", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for multiple rules, including the correct rule", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec G301 G406", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when an exclude comment is present for multiple rules, including the correct rule", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable G301 G406", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not panic if a file can not compile", func() {
			sample := testutils.SampleCodeCompilationFail[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()

			pkg.AddFile("main.go", source)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should exclude a reportable file, if excluded by build tags", func() {
			// file has a reportable security issue, but should only be flagged
			// to only being compiled in via a build flag.
			sample := testutils.SampleCodeG501BuildTag[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()

			pkg.AddFile("main.go", source)
			err := pkg.Build()
			Expect(err).To(BeEquivalentTo(&build.NoGoError{Dir: pkg.Path})) // no files should be found for scanning.
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())

			issues, _, _ := analyzer.Report()
			Expect(issues).Should(BeEmpty())
		})

		It("should attempt to analyse a file with build tags", func() {
			sample := testutils.SampleCodeBuildTag[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()

			tags := []string{"tag"}
			pkg.AddFile("main.go", source)
			err := pkg.Build(testutils.WithBuildTags(tags))
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(tags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())

			issues, _, _ := analyzer.Report()
			if len(issues) != sample.Errors {
				fmt.Println(sample.Code)
			}
			Expect(issues).Should(HaveLen(sample.Errors))
		})

		It("should report issues from a file with build tags", func() {
			sample := testutils.SampleCodeG501BuildTag[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()

			tags := []string{"tag"}
			pkg.AddFile("main.go", source)
			err := pkg.Build(testutils.WithBuildTags(tags))
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(tags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())

			issues, _, _ := analyzer.Report()
			if len(issues) != sample.Errors {
				fmt.Println(sample.Code)
			}
			Expect(issues).Should(HaveLen(sample.Errors))
		})

		It("should process an empty package with test file", func() {
			analyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo_test.go", `
				package tests
			    import "testing"
			    func TestFoo(t *testing.T){
			    }`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be possible to overwrite nosec comments, and report issues", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to overwrite disable directive, and report issues", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to overwrite nosec comments, and report issues", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to overwrite disable directive comments, and report issues", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to overwrite nosec comments, and report issues", func() {
			// Rule for MD4 weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to overwrite disable directive comments, and report issues", func() {
			// Rule for MD4 weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to overwrite nosec comments, and report issues but they should not be counted", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "mynosec")
			nosecIgnoreConfig.SetGlobal(gosec.ShowIgnored, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() // #mynosec", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, metrics, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
			Expect(metrics.NumFound).Should(Equal(0))
			Expect(metrics.NumNosec).Should(Equal(1))
		})

		It("should be possible to overwrite nosec comments, and report issues but they should not be counted", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "mynosec")
			nosecIgnoreConfig.SetGlobal(gosec.ShowIgnored, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) // #mynosec", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, metrics, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
			Expect(metrics.NumFound).Should(Equal(0))
			Expect(metrics.NumNosec).Should(Equal(1))
		})

		It("should be possible to overwrite nosec comments, and report issues but they should not be counted", func() {
			// Rule for MD4 weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.Nosec, "mynosec")
			nosecIgnoreConfig.SetGlobal(gosec.ShowIgnored, "true")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() // #mynosec", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, metrics, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
			Expect(metrics.NumFound).Should(Equal(0))
			Expect(metrics.NumNosec).Should(Equal(1))
		})

		It("should not report errors when nosec tag is in front of a line", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "//Some description\n//#nosec G401\nh := md5.New()", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when disable directive is in front of a line", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "//Some description\n//gosec:disable G401\nh := md5.New()", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when nosec tag is in front of a line", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "//Some description\n//#nosec G405\nc, e := des.NewCipher([]byte(\"mySecret\"))", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when disable directive is in front of a line", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "//Some description\n//gosec:disable G405\nc, e := des.NewCipher([]byte(\"mySecret\"))", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when nosec tag is in front of a line", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "//Some description\n//#nosec G406\nh := md4.New()", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when disable directive is in front of a line", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "//Some description\n//gosec:disable G406\nh := md4.New()", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should report errors when nosec tag is not in front of a line", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "//Some description\n//Another description #nosec G401\nh := md5.New()", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when nosec tag is not in front of a line", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "//Some description\n//Another description #nosec G405\nc, e := des.NewCipher([]byte(\"mySecret\"))", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when nosec tag is not in front of a line", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "//Some description\n//Another description #nosec G406\nh := md4.New()", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should not report errors when rules are in front of nosec tag even rules are wrong", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "//G301\n//#nosec\nh := md5.New()", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when rules are in front of nosec tag even rules are wrong", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "//G301\n//#nosec\nc, e := des.NewCipher([]byte(\"mySecret\"))", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should not report errors when rules are in front of nosec tag even rules are wrong", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "//G301\n//#nosec\nh := md4.New()", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should report errors when there are nosec tags after a #nosec WrongRuleList annotation", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "//#nosec\n//G301\n//#nosec\nh := md5.New()", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when there are disable directives after a //gosec:disable WrongRuleList", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "//gosec:disable G301\n//gosec:disable\nh := md5.New()", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when there are nosec tags after a #nosec WrongRuleList annotation", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "//#nosec\n//G301\n//#nosec\nc, e := des.NewCipher([]byte(\"mySecret\"))", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when there are disable directives after a //gosec:disable WrongRuleList", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "//gosec:disable G301\n//gosec:disable\nc, e := des.NewCipher([]byte(\"mySecret\"))", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when there are nosec tags after a #nosec WrongRuleList annotation", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "//#nosec\n//G301\n//#nosec\nh := md4.New()", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should report errors when there are disable directives after a //gosec:disable WrongRuleList", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "//gosec:disable G301\n//gosec:disable\nh := md4.New()", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := analyzer.Report()
			Expect(nosecIssues).Should(HaveLen(sample.Errors))
		})

		It("should be possible to use an alternative nosec tag", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.NoSecAlternative, "falsePositive")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() // #falsePositive", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should be possible to use an alternative nosec tag", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.NoSecAlternative, "falsePositive")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) // #falsePositive", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should be possible to use an alternative nosec tag", func() {
			// Rule for MD4 deprecated weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.NoSecAlternative, "falsePositive")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() // #falsePositive", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should ignore vulnerabilities when the default tag is found", func() {
			// Rule for MD5 weak crypto usage
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.NoSecAlternative, "falsePositive")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should ignore vulnerabilities when the default tag is found", func() {
			// Rule for DES weak crypto usage
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.NoSecAlternative, "falsePositive")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should ignore vulnerabilities when the default tag is found", func() {
			// Rule for MD4 deprecated weak crypto usage
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]

			// overwrite nosec option
			nosecIgnoreConfig := gosec.NewConfig()
			nosecIgnoreConfig.SetGlobal(gosec.NoSecAlternative, "falsePositive")
			customAnalyzer := gosec.NewAnalyzer(nosecIgnoreConfig, tests, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			nosecIssues, _, _ := customAnalyzer.Report()
			Expect(nosecIssues).Should(BeEmpty())
		})

		It("should be able to analyze Go test package", func() {
			customAnalyzer := gosec.NewAnalyzer(nil, true, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package foo
				func foo(){
				}`)
			pkg.AddFile("foo_test.go", `
				package foo_test
				import "testing"
				func test() error {
				  return nil
				}
				func TestFoo(t *testing.T){
					test()
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := customAnalyzer.Report()
			Expect(issues).Should(HaveLen(1))
		})
		It("should be able to scan generated files if NOT excluded when using the rules", func() {
			customAnalyzer := gosec.NewAnalyzer(nil, true, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package foo
				// Code generated some-generator DO NOT EDIT.
				func test() error {
				  return nil
				}
				func TestFoo(t *testing.T){
					test()
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := customAnalyzer.Report()
			Expect(issues).Should(HaveLen(1))
		})
		It("should be able to skip generated files if excluded when using the rules", func() {
			customAnalyzer := gosec.NewAnalyzer(nil, true, true, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false).RulesInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				// Code generated some-generator DO NOT EDIT.
				package foo
				func test() error {
				  return nil
				}
				func TestFoo(t *testing.T){
					test()
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := customAnalyzer.Report()
			Expect(issues).Should(BeEmpty())
		})
		It("should be able to scan generated files if NOT excluded when using the analyzes", func() {
			customAnalyzer := gosec.NewAnalyzer(nil, true, false, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false).RulesInfo())
			customAnalyzer.LoadAnalyzers(analyzers.Generate(false).AnalyzersInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package main
				// Code generated some-generator DO NOT EDIT.
        import (
          "fmt"
        )
        func main() {
          values := []string{}
          fmt.Println(values[0])
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := customAnalyzer.Report()
			Expect(issues).Should(HaveLen(1))
		})
		It("should be able to skip generated files if excluded when using the analyzes", func() {
			customAnalyzer := gosec.NewAnalyzer(nil, true, true, false, 1, logger)
			customAnalyzer.LoadRules(rules.Generate(false).RulesInfo())
			customAnalyzer.LoadAnalyzers(analyzers.Generate(false).AnalyzersInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				// Code generated some-generator DO NOT EDIT.
				package main
        import (
          "fmt"
        )
        func main() {
          values := []string{}
          fmt.Println(values[0])
				}`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = customAnalyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := customAnalyzer.Report()
			Expect(issues).Should(BeEmpty())
		})
	})
	It("should be able to analyze Cgo files", func() {
		analyzer.LoadRules(rules.Generate(false).RulesInfo())
		sample := testutils.SampleCodeCgo[0]
		source := sample.Code[0]

		testPackage := testutils.NewTestPackage()
		defer testPackage.Close()
		testPackage.AddFile("main.go", source)
		err := testPackage.Build()
		Expect(err).ShouldNot(HaveOccurred())
		err = analyzer.Process(buildTags, testPackage.Path)
		Expect(err).ShouldNot(HaveOccurred())
		issues, _, _ := analyzer.Report()
		Expect(issues).Should(BeEmpty())
	})

	Context("when parsing errors from a package", func() {
		It("should return no error when the error list is empty", func() {
			pkg := &packages.Package{}
			_, err := gosec.ParseErrors(pkg)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should properly parse the errors", func() {
			pkg := &packages.Package{
				Errors: []packages.Error{
					{
						Pos: "file:1:2",
						Msg: "build error",
					},
				},
			}
			errors, err := gosec.ParseErrors(pkg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(errors).To(HaveLen(1))
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(1))
				Expect(ferr[0].Line).To(Equal(1))
				Expect(ferr[0].Column).To(Equal(2))
				Expect(ferr[0].Err).Should(MatchRegexp(`build error`))
			}
		})

		It("should properly parse the errors without line and column", func() {
			pkg := &packages.Package{
				Errors: []packages.Error{
					{
						Pos: "file",
						Msg: "build error",
					},
				},
			}
			errors, err := gosec.ParseErrors(pkg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(errors).To(HaveLen(1))
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(1))
				Expect(ferr[0].Line).To(Equal(0))
				Expect(ferr[0].Column).To(Equal(0))
				Expect(ferr[0].Err).Should(MatchRegexp(`build error`))
			}
		})

		It("should properly parse the errors without column", func() {
			pkg := &packages.Package{
				Errors: []packages.Error{
					{
						Pos: "file",
						Msg: "build error",
					},
				},
			}
			errors, err := gosec.ParseErrors(pkg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(errors).To(HaveLen(1))
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(1))
				Expect(ferr[0].Line).To(Equal(0))
				Expect(ferr[0].Column).To(Equal(0))
				Expect(ferr[0].Err).Should(MatchRegexp(`build error`))
			}
		})

		It("should return error when line cannot be parsed", func() {
			pkg := &packages.Package{
				Errors: []packages.Error{
					{
						Pos: "file:line",
						Msg: "build error",
					},
				},
			}
			_, err := gosec.ParseErrors(pkg)
			Expect(err).Should(HaveOccurred())
		})

		It("should return error when column cannot be parsed", func() {
			pkg := &packages.Package{
				Errors: []packages.Error{
					{
						Pos: "file:1:column",
						Msg: "build error",
					},
				},
			}
			_, err := gosec.ParseErrors(pkg)
			Expect(err).Should(HaveOccurred())
		})

		It("should append  error to the same file", func() {
			pkg := &packages.Package{
				Errors: []packages.Error{
					{
						Pos: "file:1:2",
						Msg: "error1",
					},
					{
						Pos: "file:3:4",
						Msg: "error2",
					},
				},
			}
			errors, err := gosec.ParseErrors(pkg)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(errors).To(HaveLen(1))
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(2))
				Expect(ferr[0].Line).To(Equal(1))
				Expect(ferr[0].Column).To(Equal(2))
				Expect(ferr[0].Err).Should(MatchRegexp(`error1`))
				Expect(ferr[1].Line).To(Equal(3))
				Expect(ferr[1].Column).To(Equal(4))
				Expect(ferr[1].Err).Should(MatchRegexp(`error2`))
			}
		})

		It("should set the config", func() {
			config := gosec.NewConfig()
			config["test"] = "test"
			analyzer.SetConfig(config)
			found := analyzer.Config()
			Expect(config).To(Equal(found))
		})

		It("should reset the analyzer", func() {
			analyzer.Reset()
			issues, metrics, errors := analyzer.Report()
			Expect(issues).To(BeEmpty())
			Expect(*metrics).To(Equal(gosec.Metrics{}))
			Expect(errors).To(BeEmpty())
		})
	})

	Context("when appending errors", func() {
		It("should skip error for non-buildable packages", func() {
			err := &build.NoGoError{
				Dir: "pkg/test",
			}
			analyzer.AppendError("test", err)
			_, _, errors := analyzer.Report()
			Expect(errors).To(BeEmpty())
		})

		It("should add a new error", func() {
			analyzer.AppendError("file", errors.New("build error"))
			analyzer.AppendError("file", errors.New("file build error"))
			_, _, errors := analyzer.Report()
			Expect(errors).To(HaveLen(1))
			for _, ferr := range errors {
				Expect(ferr).To(HaveLen(2))
			}
		})
	})

	Context("when tracking suppressions", func() {
		BeforeEach(func() {
			analyzer = gosec.NewAnalyzer(nil, tests, false, true, 1, logger)
		})

		It("should not report an error if the violation is suppressed", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec G401 -- Justification", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Justification"))
		})

		It("should not report an error if the violation is suppressed", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable G401 -- Justification", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Justification"))
		})

		It("should not report an error if the violation is suppressed", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec G405 -- Justification", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Justification"))
		})

		It("should not report an error if the violation is suppressed", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable G405 -- Justification", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Justification"))
		})

		It("should not report an error if the violation is suppressed", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec G406 -- Justification", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Justification"))
		})

		It("should not report an error if the violation is suppressed", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable G406 -- Justification", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Justification"))
		})

		It("should not report an error if the violation is suppressed without certain rules", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //#nosec", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal(""))
		})

		It("should not report an error if the violation is suppressed without certain rules", func() {
			sample := testutils.SampleCodeG401[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G401")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md5.New()", "h := md5.New() //gosec:disable", 1)
			nosecPackage.AddFile("md5.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal(""))
		})

		It("should not report an error if the violation is suppressed without certain rules", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //#nosec", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal(""))
		})

		It("should not report an error if the violation is suppressed without certain rules", func() {
			sample := testutils.SampleCodeG405[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G405")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "c, e := des.NewCipher([]byte(\"mySecret\"))", "c, e := des.NewCipher([]byte(\"mySecret\")) //gosec:disable", 1)
			nosecPackage.AddFile("cipher.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal(""))
		})

		It("should not report an error if the violation is suppressed without certain rules", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //#nosec", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal(""))
		})

		It("should not report an error if the violation is suppressed without certain rules", func() {
			sample := testutils.SampleCodeG406[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G406")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "h := md4.New()", "h := md4.New() //gosec:disable", 1)
			nosecPackage.AddFile("md4.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal(""))
		})

		It("should not report an error if the rule is not included", func() {
			sample := testutils.SampleCodeG101[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(true, rules.NewRuleFilter(false, "G401")).RulesInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("pwd.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			controlIssues, _, _ := analyzer.Report()
			Expect(controlIssues).Should(HaveLen(sample.Errors))
			Expect(controlIssues[0].Suppressions).To(HaveLen(1))
			Expect(controlIssues[0].Suppressions[0].Kind).To(Equal("external"))
			Expect(controlIssues[0].Suppressions[0].Justification).To(Equal("Globally suppressed."))
		})

		It("should not report an error if the rule is excluded", func() {
			sample := testutils.SampleCodeG101[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(true, rules.NewRuleFilter(true, "G101")).RulesInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("pwd.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).Should(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("external"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Globally suppressed."))
		})

		It("should not report an error if the analyzer is not included", func() {
			sample := testutils.SampleCodeG407[0]
			source := sample.Code[0]
			analyzer.LoadAnalyzers(analyzers.Generate(true, analyzers.NewAnalyzerFilter(false, "G115")).AnalyzersInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("cipher.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			controlIssues, _, _ := analyzer.Report()
			Expect(controlIssues).Should(HaveLen(sample.Errors))
			Expect(controlIssues[0].Suppressions).To(HaveLen(1))
			Expect(controlIssues[0].Suppressions[0].Kind).To(Equal("external"))
			Expect(controlIssues[0].Suppressions[0].Justification).To(Equal("Globally suppressed."))
		})

		It("should not report an error if the analyzer is excluded", func() {
			sample := testutils.SampleCodeG407[0]
			source := sample.Code[0]
			analyzer.LoadAnalyzers(analyzers.Generate(true, analyzers.NewAnalyzerFilter(true, "G407")).AnalyzersInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("cipher.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).Should(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("external"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Globally suppressed."))
		})

		It("should not report an error if the analyzer is not included", func() {
			sample := testutils.SampleCodeG602[0]
			source := sample.Code[0]
			analyzer.LoadAnalyzers(analyzers.Generate(true, analyzers.NewAnalyzerFilter(false, "G115")).AnalyzersInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("cipher.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			controlIssues, _, _ := analyzer.Report()
			Expect(controlIssues).Should(HaveLen(sample.Errors))
			Expect(controlIssues[0].Suppressions).To(HaveLen(1))
			Expect(controlIssues[0].Suppressions[0].Kind).To(Equal("external"))
			Expect(controlIssues[0].Suppressions[0].Justification).To(Equal("Globally suppressed."))
		})

		It("should not report an error if the analyzer is excluded", func() {
			sample := testutils.SampleCodeG602[0]
			source := sample.Code[0]
			analyzer.LoadAnalyzers(analyzers.Generate(true, analyzers.NewAnalyzerFilter(true, "G602")).AnalyzersInfo())

			controlPackage := testutils.NewTestPackage()
			defer controlPackage.Close()
			controlPackage.AddFile("cipher.go", source)
			err := controlPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, controlPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).Should(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("external"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("Globally suppressed."))
		})

		It("should track multiple suppressions if the violation is multiply suppressed", func() {
			sample := testutils.SampleCodeG101[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(true, rules.NewRuleFilter(true, "G101")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source, "password := \"f62e5bcda4fae4f82370da0c6f20697b8f8447ef\"", "password := \"f62e5bcda4fae4f82370da0c6f20697b8f8447ef\" //#nosec G101 -- Justification", 1)
			nosecPackage.AddFile("pwd.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).Should(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(2))
		})

		It("should not report an error if the violation is suppressed on a struct filed", func() {
			sample := testutils.SampleCodeG402[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G402")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source,
				"TLSClientConfig: &tls.Config{InsecureSkipVerify: true}",
				"TLSClientConfig: &tls.Config{InsecureSkipVerify: true} // #nosec G402", 1)
			nosecPackage.AddFile("tls.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
		})

		It("should not report an error if the violation is suppressed on a struct filed", func() {
			sample := testutils.SampleCodeG402[0]
			source := sample.Code[0]
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G402")).RulesInfo())

			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecSource := strings.Replace(source,
				"TLSClientConfig: &tls.Config{InsecureSkipVerify: true}",
				"TLSClientConfig: &tls.Config{InsecureSkipVerify: true} //gosec:disable G402", 1)
			nosecPackage.AddFile("tls.go", nosecSource)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(sample.Errors))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
		})

		It("should not report an error if the violation is suppressed on multi-lien issue", func() {
			source := `
package main

import (
	"fmt"
)

const TokenLabel = `
			source += "`" + `
f62e5bcda4fae4f82370da0c6f20697b8f8447ef
      ` + "`" + "//#nosec G101 -- false positive, this is not a private data" + `
func main() {
	fmt.Printf("Label: %s ", TokenLabel)
}
      `
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G101")).RulesInfo())
			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("pwd.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(1))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("false positive, this is not a private data"))
		})

		It("should not report an error if the violation is suppressed on multi-lien issue", func() {
			source := `
package main

import (
	"fmt"
)

const TokenLabel = `
			source += "`" + `
f62e5bcda4fae4f82370da0c6f20697b8f8447ef
      ` + "`" + "//gosec:disable G101 -- false positive, this is not a private data" + `
func main() {
	fmt.Printf("Label: %s ", TokenLabel)
}
      `
			analyzer.LoadRules(rules.Generate(false, rules.NewRuleFilter(false, "G101")).RulesInfo())
			nosecPackage := testutils.NewTestPackage()
			defer nosecPackage.Close()
			nosecPackage.AddFile("pwd.go", source)
			err := nosecPackage.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, nosecPackage.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			Expect(issues).To(HaveLen(1))
			Expect(issues[0].Suppressions).To(HaveLen(1))
			Expect(issues[0].Suppressions[0].Kind).To(Equal("inSource"))
			Expect(issues[0].Suppressions[0].Justification).To(Equal("false positive, this is not a private data"))
		})
	})

	Context("when fixing issue #1240 - nosec with open bracket", func() {
		It("should suppress G115 when #nosec is at the end of an if line with bracket", func() {
			source := `
package main

import "fmt"

func main() {
	ten := 10
	uintTen := uint(10)
	configVal := uint(ten) // #nosec G115 -- this works
	inputSlice := []int{1, 2, 3, 4, 5}

	if len(inputSlice) <= int(uintTen) { // #nosec G115 -- this works
		fmt.Println("hello world!")
	}

	if len(inputSlice) <= int(configVal) { // #nosec G115 -- this should work now (fix for #1240)
		fmt.Println("hello world!")
	}
}
`
			analyzer.LoadAnalyzers(analyzers.Generate(false, analyzers.NewAnalyzerFilter(false, "G115")).AnalyzersInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("main.go", source)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, metrics, _ := analyzer.Report()
			// No G115 issues should be reported as all conversions are suppressed
			for _, issue := range issues {
				if issue.RuleID == "G115" {
					Fail(fmt.Sprintf("G115 should be suppressed but was reported at line %s", issue.Line))
				}
			}
			Expect(metrics.NumNosec).Should(BeNumerically(">=", 3)) // At least 3 nosec comments
		})

		It("should suppress G115 when #nosec is used with block comment before bracket", func() {
			source := `
package main

import "fmt"

func main() {
	configVal := uint(10)
	inputSlice := []int{1, 2, 3, 4, 5}

	if len(inputSlice) <= int(configVal) /* #nosec G115 */ {
		fmt.Println("hello world!")
	}
}
`
			analyzer.LoadAnalyzers(analyzers.Generate(false, analyzers.NewAnalyzerFilter(false, "G115")).AnalyzersInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("main.go", source)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			// No G115 issues should be reported
			for _, issue := range issues {
				if issue.RuleID == "G115" {
					Fail(fmt.Sprintf("G115 should be suppressed but was reported at line %s", issue.Line))
				}
			}
		})

		It("should suppress G115 in for loop with bracket and trailing comment", func() {
			source := `
package main

func main() {
	x := uint(10)
	for i := 0; i < int(x); i++ { // #nosec G115
		println(i)
	}
}
`
			analyzer.LoadAnalyzers(analyzers.Generate(false, analyzers.NewAnalyzerFilter(false, "G115")).AnalyzersInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("main.go", source)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			// No G115 issues should be reported
			for _, issue := range issues {
				if issue.RuleID == "G115" {
					Fail(fmt.Sprintf("G115 should be suppressed but was reported at line %s", issue.Line))
				}
			}
		})

		It("should suppress G115 in switch statement with bracket and trailing comment", func() {
			source := `
package main

func main() {
	x := uint(10)
	switch int(x) { // #nosec G115
	case 10:
		println("ten")
	}
}
`
			analyzer.LoadAnalyzers(analyzers.Generate(false, analyzers.NewAnalyzerFilter(false, "G115")).AnalyzersInfo())
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("main.go", source)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			err = analyzer.Process(buildTags, pkg.Path)
			Expect(err).ShouldNot(HaveOccurred())
			issues, _, _ := analyzer.Report()
			// No G115 issues should be reported
			for _, issue := range issues {
				if issue.RuleID == "G115" {
					Fail(fmt.Sprintf("G115 should be suppressed but was reported at line %s", issue.Line))
				}
			}
		})
	})
})
