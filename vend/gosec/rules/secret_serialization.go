package rules

import (
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
)

type secretSerialization struct {
	issue.MetaData
	pattern *regexp.Regexp
}

func (r *secretSerialization) Match(n ast.Node, ctx *gosec.Context) (*issue.Issue, error) {
	field, ok := n.(*ast.Field)
	if !ok || len(field.Names) == 0 {
		return nil, nil // skip embedded (anonymous) fields
	}

	// Parse the JSON tag to determine behavior
	omitted := false
	jsonKey := ""

	if field.Tag != nil {
		if tagVal, err := strconv.Unquote(field.Tag.Value); err == nil && tagVal != "" {
			st := reflect.StructTag(tagVal)
			if tag := st.Get("json"); tag != "" {
				if tag == "-" {
					omitted = true
				} else {
					// "name,omitempty" -> "name"
					// "-," -> "-" (A field literally named "-")
					parts := strings.SplitN(tag, ",", 2)
					jsonKey = parts[0]
				}
			}
		}
	}

	if omitted {
		return nil, nil
	}

	// Iterate over all names in this field definition
	// e.g., type T struct { Pass, Salt string }
	isSensitiveType := false
	switch t := field.Type.(type) {
	case *ast.Ident:
		if t.Name == "string" {
			isSensitiveType = true
		}
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok && ident.Name == "string" {
			isSensitiveType = true
		}
	case *ast.ArrayType:
		if star, ok := t.Elt.(*ast.StarExpr); ok {
			if ident, ok := star.X.(*ast.Ident); ok && ident.Name == "string" {
				isSensitiveType = true // []*string
			}
		} else if ident, ok := t.Elt.(*ast.Ident); ok {
			if ident.Name == "string" || ident.Name == "byte" {
				isSensitiveType = true // []string or []byte
			}
		}
	}

	if !isSensitiveType {
		return nil, nil
	}

	// Check each named exported field
	for _, nameIdent := range field.Names {
		fieldName := nameIdent.Name
		if fieldName == "_" || !ast.IsExported(fieldName) {
			continue
		}

		effectiveKey := jsonKey
		if effectiveKey == "" {
			effectiveKey = fieldName
		}

		if gosec.RegexMatchWithCache(r.pattern, fieldName) || gosec.RegexMatchWithCache(r.pattern, effectiveKey) {
			msg := fmt.Sprintf("Exported struct field %q (JSON key %q) matches secret pattern", fieldName, effectiveKey)
			return ctx.NewIssue(field, r.ID(), msg, r.Severity, r.Confidence), nil
		}
	}

	return nil, nil
}

func NewSecretSerialization(id string, conf gosec.Config) (gosec.Rule, []ast.Node) {
	patternStr := `(?i)\b((?:api|access|auth|bearer|client|oauth|private|refresh|session|jwt)[_-]?(?:key|secret|token)s?|password|passwd|pwd|pass|secret|cred|jwt)\b`

	if val, ok := conf[id]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			if p, ok := m["pattern"].(string); ok && p != "" {
				patternStr = p
			}
		}
	}

	return &secretSerialization{
		pattern:  regexp.MustCompile(patternStr),
		MetaData: issue.NewMetaData(id, "Exported struct field appears to be a secret and is not ignored by JSON marshaling", issue.Medium, issue.Medium),
	}, []ast.Node{(*ast.Field)(nil)}
}
