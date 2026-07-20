package goanalysis

import (
	"sync/atomic"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

type loadingStats struct {
	alreadyLoadedPackages atomic.Int64
	sourceLoads           atomic.Int64
	exportLoads           atomic.Int64
	exportLoadFallbacks   atomic.Int64
	factCacheHits         atomic.Int64
	factCacheMisses       atomic.Int64
}

func (s *loadingStats) log(log logutils.Log) {
	log.Infof("Goanalysis package loading stats: already_loaded=%d, from_source=%d, from_export=%d, "+
		"export_fallbacks=%d, fact_cache_hits=%d, fact_cache_misses=%d",
		s.alreadyLoadedPackages.Load(),
		s.sourceLoads.Load(),
		s.exportLoads.Load(),
		s.exportLoadFallbacks.Load(),
		s.factCacheHits.Load(),
		s.factCacheMisses.Load())
}
