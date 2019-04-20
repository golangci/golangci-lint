package processors

import (
	"go/token"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/golangci/golangci-lint/pkg/lint/astcache"

	"github.com/golangci/golangci-lint/pkg/result"
)

type posMapper func(pos token.Position) token.Position

// FilenameUnadjuster is needed because a lot of linters use fset.Position(f.Pos())
// to get filename. And they return adjusted filename (e.g. *.qtpl) for an issue. We need
// restore real .go filename to properly output it, parse it, etc.
type FilenameUnadjuster struct {
	m   map[string]posMapper // map from adjusted filename to position mapper: adjusted -> unadjusted position
	log logutils.Log
}

var _ Processor = FilenameUnadjuster{}

func NewFilenameUnadjuster(cache *astcache.Cache, log logutils.Log) *FilenameUnadjuster {
	m := map[string]posMapper{}
	for _, f := range cache.GetAllValidFiles() {
		adjustedFilename := f.Fset.PositionFor(f.F.Pos(), true).Filename
		if adjustedFilename == "" {
			continue
		}
		unadjustedFilename := f.Fset.PositionFor(f.F.Pos(), false).Filename
		if unadjustedFilename == "" || unadjustedFilename == adjustedFilename {
			continue
		}
		if !strings.HasSuffix(unadjustedFilename, ".go") {
			continue // file.go -> /caches/cgo-xxx
		}

		f := f
		m[adjustedFilename] = func(adjustedPos token.Position) token.Position {
			tokenFile := f.Fset.File(f.F.Pos())
			if tokenFile == nil {
				log.Warnf("Failed to get token file for %s", adjustedFilename)
				return adjustedPos
			}
			return f.Fset.PositionFor(tokenFile.Pos(adjustedPos.Offset), false)
		}
	}

	return &FilenameUnadjuster{
		m:   m,
		log: log,
	}
}

func (p FilenameUnadjuster) Name() string {
	return "filename_unadjuster"
}

func (p FilenameUnadjuster) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(i *result.Issue) *result.Issue {
		issueFilePath := i.FilePath()
		if !filepath.IsAbs(i.FilePath()) {
			absPath, err := filepath.Abs(i.FilePath())
			if err != nil {
				p.log.Warnf("failed to build abs path for %q: %s", i.FilePath(), err)
				return i
			}
			issueFilePath = absPath
		}

		mapper := p.m[issueFilePath]
		if mapper == nil {
			return i
		}

		newI := *i
		newI.Pos = mapper(i.Pos)
		p.log.Infof("Unadjusted from %v to %v", i.Pos, newI.Pos)
		return &newI
	}), nil
}

func (FilenameUnadjuster) Finish() {}
