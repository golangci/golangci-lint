package analyzers_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/packages"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/analyzers"
	"github.com/securego/gosec/v2/testutils"
)

func benchmarkAnalyzerStress(b *testing.B, analyzerID string, generator func() string) {
	logger, _ := testutils.NewLogger()
	code := generator()

	// SETUP: Create temp dir and main.go
	tmpDir, err := os.MkdirTemp("", "gosec_bench")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	mainGo := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGo, []byte(code), 0o600); err != nil {
		b.Fatalf("failed to write main.go: %v", err)
	}

	// Create a dummy go.mod to ensure we are in a module
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module bench\n\ngo 1.24\n"), 0o600); err != nil {
		b.Fatalf("failed to write go.mod: %v", err)
	}

	conf := &packages.Config{
		Mode: gosec.LoadMode,
		Dir:  tmpDir,
	}
	pkgs, err := packages.Load(conf, ".")
	if err != nil {
		b.Fatalf("failed to load package: %v", err)
	}
	if len(pkgs) == 0 {
		b.Fatalf("no packages loaded")
	}
	if len(pkgs[0].Errors) > 0 {
		b.Fatalf("errors loading package: %v", pkgs[0].Errors)
	}

	// Prepare analysis context
	pass := &analysis.Pass{
		Fset:       pkgs[0].Fset,
		Files:      pkgs[0].Syntax,
		Pkg:        pkgs[0].Types,
		TypesInfo:  pkgs[0].TypesInfo,
		TypesSizes: pkgs[0].TypesSizes,
		ResultOf:   make(map[*analysis.Analyzer]any),
		Report:     func(d analysis.Diagnostic) {},
	}

	pass.Analyzer = inspect.Analyzer
	i, _ := inspect.Analyzer.Run(pass)
	pass.ResultOf[inspect.Analyzer] = i

	pass.Analyzer = ctrlflow.Analyzer
	cf, _ := ctrlflow.Analyzer.Run(pass)
	pass.ResultOf[ctrlflow.Analyzer] = cf

	pass.Analyzer = buildssa.Analyzer
	ssaRes, err := buildssa.Analyzer.Run(pass)
	if err != nil {
		b.Fatalf("failed to build SSA: %v", err)
	}
	ssaResult := ssaRes.(*buildssa.SSA)

	if len(ssaResult.SrcFuncs) == 0 {
		b.Fatalf("SSA has 0 source functions.")
	}

	// Find targeted analyzer
	var target *analysis.Analyzer
	analyzerList := analyzers.Generate(false)
	if def, ok := analyzerList.Analyzers[analyzerID]; ok {
		target = def.Create(def.ID, def.Description)
	} else {
		b.Fatalf("analyzer %s not found", analyzerID)
	}

	resultMap := map[*analysis.Analyzer]any{
		buildssa.Analyzer: &analyzers.SSAAnalyzerResult{
			Config: gosec.NewConfig(),
			Logger: logger,
			SSA:    ssaResult,
		},
	}

	runPass := &analysis.Pass{
		Analyzer:   target,
		Fset:       pkgs[0].Fset,
		Files:      pkgs[0].Syntax,
		Pkg:        pkgs[0].Types,
		TypesInfo:  pkgs[0].TypesInfo,
		TypesSizes: pkgs[0].TypesSizes,
		ResultOf:   resultMap,
		Report:     func(d analysis.Diagnostic) {},
	}

	b.ResetTimer()
	for range b.N {
		_, err := target.Run(runPass)
		if err != nil {
			b.Fatalf("failed to run analyzer: %v", err)
		}
	}
}

// Generators

func generateG115Deep(nesting, conversions int) string {
	var sb strings.Builder
	sb.WriteString("package main\nimport \"math\"\nfunc run_stress(x int64) {\n")
	for i := range nesting {
		fmt.Fprintf(&sb, "if x > %d && x < math.MaxInt64 {\n", i)
	}
	for range conversions {
		fmt.Fprintf(&sb, "_ = int8(x)\n")
	}
	for range nesting {
		sb.WriteString("}\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

func generateG602Wide(levels, accesses int) string {
	var sb strings.Builder
	sb.WriteString("package main\nfunc run_stress() {\n")
	sb.WriteString("s := make([]byte, 100000)\n")
	for i := range levels {
		fmt.Fprintf(&sb, "s%d := s[%d:]\n", i, i)
		for j := range accesses {
			fmt.Fprintf(&sb, "_ = s%d[%d]\n", i, j)
			fmt.Fprintf(&sb, "_ = s%d[%d]\n", i, j+1)
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}

func generateG407Stress(depth int) string {
	var sb strings.Builder
	sb.WriteString("package main\nimport \"crypto/cipher\"\nfunc run_stress(gcm cipher.AEAD, data []byte) {\n")
	sb.WriteString("nonce := []byte(\"hardcoded_nonce_value\")\n")
	// Chain of assignments
	for i := range depth {
		fmt.Fprintf(&sb, "n%d := nonce\n", i)
		if i > 0 {
			fmt.Fprintf(&sb, "n%d = n%d\n", i, i-1)
		}
	}
	// Use the last nonce in the chain
	fmt.Fprintf(&sb, "gcm.Seal(nil, n%d, data, nil)\n", depth-1)
	fmt.Fprintf(&sb, "}\n")
	return sb.String()
}

// Benchmarks (Logic Only)

func BenchmarkAnalysisG115_Deep(b *testing.B) {
	benchmarkAnalyzerStress(b, "G115", func() string { return generateG115Deep(300, 1000) })
}

func BenchmarkAnalysisG602_Wide(b *testing.B) {
	benchmarkAnalyzerStress(b, "G602", func() string { return generateG602Wide(500, 200) })
}

func BenchmarkAnalysisG407_Deep(b *testing.B) {
	benchmarkAnalyzerStress(b, "G407", func() string { return generateG407Stress(1000) })
}

func generateComplex(functions, complexity int) string {
	var sb strings.Builder
	sb.WriteString("package main\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"math\"\n")
	sb.WriteString("\t\"crypto/cipher\"\n")
	sb.WriteString(")\n")

	// Generate helper functions that call each other
	for i := range functions {
		fmt.Fprintf(&sb, "func complexFunction%d(x int64, s []byte, gcm cipher.AEAD) {\n", i)

		// G115 logic: conversions in branches
		for j := range complexity {
			fmt.Fprintf(&sb, "\tif x > %d && x < math.MaxInt64 {\n", j)
			fmt.Fprintf(&sb, "\t\t_ = int8(x)\n")
			fmt.Fprintf(&sb, "\t}\n")
		}

		// G602 logic: slice operations
		fmt.Fprintf(&sb, "\t_ = s[%d]\n", i%10)
		for j := range complexity {
			fmt.Fprintf(&sb, "\tif len(s) > %d {\n", j)
			fmt.Fprintf(&sb, "\t\t_ = s[%d]\n", j)
			fmt.Fprintf(&sb, "\t}\n")
		}

		// G407 logic: nonce passing (simulated)
		fmt.Fprintf(&sb, "\tnonce := []byte(\"hardcoded_nonce_%d\")\n", i)
		fmt.Fprintf(&sb, "\tgcm.Seal(nil, nonce, s, nil)\n")

		// Call next function if not last
		if i < functions-1 {
			fmt.Fprintf(&sb, "\tcomplexFunction%d(x, s, gcm)\n", i+1)
		}
		sb.WriteString("}\n")
	}

	sb.WriteString("func run_stress() {\n")
	sb.WriteString("\ts := make([]byte, 10000)\n")
	sb.WriteString("\tcomplexFunction0(100, s, nil)\n")
	sb.WriteString("}\n")

	return sb.String()
}

func BenchmarkAnalysisG115_Complex(b *testing.B) {
	benchmarkAnalyzerStress(b, "G115", func() string { return generateComplex(50, 20) })
}

func BenchmarkAnalysisG602_Complex(b *testing.B) {
	benchmarkAnalyzerStress(b, "G602", func() string { return generateComplex(50, 20) })
}

func BenchmarkAnalysisG407_Complex(b *testing.B) {
	benchmarkAnalyzerStress(b, "G407", func() string { return generateComplex(50, 20) })
}
