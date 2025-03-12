package golines

import (
	"github.com/golangci/golines"

	"github.com/golangci/golangci-lint/v2/pkg/config"
)

const Name = "golines"

type Formatter struct {
	shortener *golines.Shortener
}

func New(settings *config.GoLinesSettings) *Formatter {
	options := golines.ShortenerConfig{}

	if settings != nil {
		options = golines.ShortenerConfig{
			MaxLen:           settings.MaxLen,
			TabLen:           settings.TabLen,
			KeepAnnotations:  false, // golines debug (not usable inside golangci-lint)
			ShortenComments:  settings.ShortenComments,
			ReformatTags:     settings.ReformatTags,
			IgnoreGenerated:  false, // handle globally
			DotFile:          "",    // golines debug (not usable inside golangci-lint)
			ChainSplitDots:   settings.ChainSplitDots,
			BaseFormatterCmd: "go fmt", // fake cmd
		}
	}

	return &Formatter{shortener: golines.NewShortener(options)}
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(_ string, src []byte) ([]byte, error) {
	return f.shortener.Shorten(src)
}
