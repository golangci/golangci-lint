package goanalysis

import (
	"golang.org/x/tools/go/analysis"
)

type SupportedLinter interface {
	Analyzers() []*analysis.Analyzer
	Cfg() map[string]map[string]interface{}
	AnalyzerToLinterNameMapping() map[*analysis.Analyzer]string
}
