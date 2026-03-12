package gosec_test

import (
	"fmt"
	"go/ast"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
)

type mockrule struct {
	issue    *issue.Issue
	err      error
	callback func(n ast.Node, ctx *gosec.Context) bool
}

func (m *mockrule) ID() string {
	return "MOCK"
}

func (m *mockrule) Match(n ast.Node, ctx *gosec.Context) (*issue.Issue, error) {
	if m.callback(n, ctx) {
		return m.issue, nil
	}
	return nil, m.err
}

var _ = Describe("Rule", func() {
	Context("when using a ruleset", func() {
		var (
			ruleset        gosec.RuleSet
			dummyErrorRule gosec.Rule
			dummyIssueRule gosec.Rule
		)

		JustBeforeEach(func() {
			ruleset = gosec.NewRuleSet()
			dummyErrorRule = &mockrule{
				issue:    nil,
				err:      fmt.Errorf("An unexpected error occurred"),
				callback: func(n ast.Node, ctx *gosec.Context) bool { return false },
			}
			dummyIssueRule = &mockrule{
				issue: &issue.Issue{
					Severity:   issue.High,
					Confidence: issue.High,
					What:       `Some explanation of the thing`,
					File:       "main.go",
					Code:       `#include <stdio.h> int main(){ puts("hello world"); }`,
					Line:       "42",
				},
				err:      nil,
				callback: func(n ast.Node, ctx *gosec.Context) bool { return true },
			}
		})
		It("should be possible to register a rule for multiple ast.Node", func() {
			registeredNodeA := (*ast.CallExpr)(nil)
			registeredNodeB := (*ast.AssignStmt)(nil)
			unregisteredNode := (*ast.BinaryExpr)(nil)

			ruleset.Register(dummyIssueRule, false, registeredNodeA, registeredNodeB)
			Expect(ruleset.RegisteredFor(unregisteredNode)).Should(BeEmpty())
			Expect(ruleset.RegisteredFor(registeredNodeA)).Should(ContainElement(dummyIssueRule))
			Expect(ruleset.RegisteredFor(registeredNodeB)).Should(ContainElement(dummyIssueRule))
			Expect(ruleset.IsRuleSuppressed(dummyIssueRule.ID())).Should(BeFalse())
		})

		It("should not register a rule when no ast.Nodes are specified", func() {
			ruleset.Register(dummyErrorRule, false)
			Expect(ruleset.Rules).Should(BeEmpty())
		})

		It("should be possible to retrieve a list of rules for a given node type", func() {
			registeredNode := (*ast.CallExpr)(nil)
			unregisteredNode := (*ast.AssignStmt)(nil)
			ruleset.Register(dummyErrorRule, false, registeredNode)
			ruleset.Register(dummyIssueRule, false, registeredNode)
			Expect(ruleset.RegisteredFor(unregisteredNode)).Should(BeEmpty())
			Expect(ruleset.RegisteredFor(registeredNode)).Should(HaveLen(2))
			Expect(ruleset.RegisteredFor(registeredNode)).Should(ContainElement(dummyErrorRule))
			Expect(ruleset.RegisteredFor(registeredNode)).Should(ContainElement(dummyIssueRule))
		})

		It("should register a suppressed rule", func() {
			registeredNode := (*ast.CallExpr)(nil)
			unregisteredNode := (*ast.AssignStmt)(nil)
			ruleset.Register(dummyIssueRule, true, registeredNode)
			Expect(ruleset.RegisteredFor(registeredNode)).Should(ContainElement(dummyIssueRule))
			Expect(ruleset.RegisteredFor(unregisteredNode)).Should(BeEmpty())
			Expect(ruleset.IsRuleSuppressed(dummyIssueRule.ID())).Should(BeTrue())
		})
	})
})
