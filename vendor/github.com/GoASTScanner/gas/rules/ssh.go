package rules

import (
	"go/ast"

	"github.com/GoASTScanner/gas"
)

type sshHostKey struct {
	gas.MetaData
	pkg   string
	calls []string
}

func (r *sshHostKey) ID() string {
	return r.MetaData.ID
}

func (r *sshHostKey) Match(n ast.Node, c *gas.Context) (gi *gas.Issue, err error) {
	if _, matches := gas.MatchCallByPackage(n, c, r.pkg, r.calls...); matches {
		return gas.NewIssue(c, n, r.ID(), r.What, r.Severity, r.Confidence), nil
	}
	return nil, nil
}

// NewSSHHostKey rule detects the use of insecure ssh HostKeyCallback.
func NewSSHHostKey(id string, conf gas.Config) (gas.Rule, []ast.Node) {
	return &sshHostKey{
		pkg:   "golang.org/x/crypto/ssh",
		calls: []string{"InsecureIgnoreHostKey"},
		MetaData: gas.MetaData{
			ID:         id,
			What:       "Use of ssh InsecureIgnoreHostKey should be audited",
			Severity:   gas.Medium,
			Confidence: gas.High,
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
