package golines

import (
	"github.com/golangci/golines"

	"github.com/golangci/golangci-lint/pkg/config"
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
			KeepAnnotations:  false, // debug
			ShortenComments:  settings.ShortenComments,
			ReformatTags:     settings.ReformatTags,
			IgnoreGenerated:  false, // handle globally
			DotFile:          "",    // debug
			ChainSplitDots:   settings.ChainSplitDots,
			BaseFormatterCmd: "fmt", // fake cmd
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
