package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/breml/errchkjson"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func NewErrChkJSONFuncName(linters *config.Linters,
	cfg *config.ErrChkJSONSettings,
	errcheckCfg *config.ErrcheckSettings,
) *goanalysis.Linter {
	a := errchkjson.NewAnalyzer()

	var omitSafe bool
	var reportNoExported bool
	if cfg != nil {
		omitSafe = cfg.OmitSafe
		reportNoExported = cfg.ReportNoExported
	}

	// Modify errcheck config if this linter is enabled and OmitSafe is false
	if isEnabled(linters, a.Name) && !omitSafe {
		if errcheckCfg == nil {
			errcheckCfg = &config.ErrcheckSettings{}
		}
		errcheckCfg.ExcludeFunctions = append(errcheckCfg.ExcludeFunctions,
			"encoding/json.Marshal",
			"encoding/json.MarshalIndent",
			"(*encoding/json.Encoder).Encode",
		)
	}

	cfgMap := map[string]map[string]interface{}{}
	cfgMap[a.Name] = map[string]interface{}{
		"omit-safe":          omitSafe,
		"report-no-exported": reportNoExported,
	}

	return goanalysis.NewLinter(
		"errchkjson",
		"Checks types passed to the json encoding functions. "+
			"Reports unsupported types and reports occations, where the check for the returned error can be omitted.",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func isEnabled(linters *config.Linters, linterName string) bool {
	if linters != nil {
		var enabled bool
		for _, linter := range linters.Enable {
			if linter == linterName {
				enabled = true
				break
			}
		}
		var disabled bool
		for _, linter := range linters.Disable {
			if linter == linterName {
				disabled = true
				break
			}
		}
		var presetEnabled bool
		for _, preset := range linters.Presets {
			if preset == linter.PresetBugs || preset == linter.PresetUnused {
				presetEnabled = true
				break
			}
		}
		return enabled ||
			linters.EnableAll && !disabled ||
			presetEnabled && !disabled
	}
	return false
}
