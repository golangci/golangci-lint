package fsutils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/ldez/grignotin/goenv"
)

// Relative path modes.
const (
	RelativePathModeGoMod   = "gomod"
	RelativePathModeGitRoot = "gitroot"
	RelativePathModeCfg     = "cfg"
	RelativePathModeWd      = "wd"
)

func AllRelativePathModes() []string {
	return []string{RelativePathModeGoMod, RelativePathModeGitRoot, RelativePathModeCfg, RelativePathModeWd}
}

func GetBasePath(ctx context.Context, mode, cfgDir string) (string, error) {
	switch mode {
	case RelativePathModeGoMod:
		goMod, err := goenv.GetOne(ctx, goenv.GOMOD)
		if err != nil {
			return "", fmt.Errorf("get go.mod path: %w", err)
		}

		return filepath.Dir(goMod), nil

	case RelativePathModeGitRoot:
		root, err := gitRoot(ctx)
		if err != nil {
			return "", fmt.Errorf("get git root: %w", err)
		}

		return root, nil

	case RelativePathModeCfg:
		if cfgDir == "" {
			return "", fmt.Errorf("missing configuration directory")
		}

		return cfgDir, nil

	default:
		// mode "wd"
		wd, err := Getwd()
		if err != nil {
			return "", fmt.Errorf("get wd: %w", err)
		}

		return wd, nil
	}
}

func gitRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out)), nil
}
