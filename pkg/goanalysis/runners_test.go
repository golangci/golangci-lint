package goanalysis

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/v2/internal/cache"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis/pkgerrors"
)

func TestIssuesCacheHashMode(t *testing.T) {
	testCases := []struct {
		desc     string
		loadMode LoadMode
		want     cache.HashMode
	}{
		{
			desc:     "none uses only self",
			loadMode: LoadModeNone,
			want:     cache.HashModeNeedOnlySelf,
		},
		{
			desc:     "syntax uses only self",
			loadMode: LoadModeSyntax,
			want:     cache.HashModeNeedOnlySelf,
		},
		{
			desc:     "types info uses all dependencies",
			loadMode: LoadModeTypesInfo,
			want:     cache.HashModeNeedAllDeps,
		},
		{
			desc:     "whole program uses all dependencies",
			loadMode: LoadModeWholeProgram,
			want:     cache.HashModeNeedAllDeps,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			require.Equal(t, tc.want, issuesCacheHashMode(tc.loadMode))
		})
	}
}

func TestPackagesToSaveIssuesFor(t *testing.T) {
	good := &packages.Package{Name: "good"}
	bad := &packages.Package{Name: "bad", IllTyped: true}
	pkgs := []*packages.Package{good, bad}

	testCases := []struct {
		desc   string
		errs   []error
		want   []*packages.Package
		wantOk bool
	}{
		{
			desc:   "no errors saves all packages",
			want:   pkgs,
			wantOk: true,
		},
		{
			desc:   "ill typed errors save only well typed packages",
			errs:   []error{fmt.Errorf("analyzer: %w", &pkgerrors.IllTypedError{Pkg: bad})},
			want:   []*packages.Package{good},
			wantOk: true,
		},
		{
			desc:   "non typecheck errors skip saving",
			errs:   []error{errors.New("analyzer failed")},
			wantOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, gotOk := packagesToSaveIssuesFor(pkgs, tc.errs)

			require.Equal(t, tc.wantOk, gotOk)
			require.Equal(t, tc.want, got)
		})
	}
}
