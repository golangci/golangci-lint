package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/tools/go/analysis"
)

const (
	LoadModeSyntax    = "syntax"
	LoadModeTypesInfo = "typesinfo"
)

var (
	pluginsMu sync.RWMutex
	plugins   = make(map[string]NewPlugin)
)

type LinterPlugin interface {
	BuildAnalyzers() ([]*analysis.Analyzer, error)
	GetLoadMode() string
}

type NewPlugin func(conf any) (LinterPlugin, error)

func Plugin(name string, p NewPlugin) {
	pluginsMu.Lock()

	plugins[name] = p

	pluginsMu.Unlock()
}

func GetPlugin(name string) (NewPlugin, error) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	p, ok := plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %q not found", name)
	}

	return p, nil
}

func DecodeSettings[T any](rawSettings any) (T, error) {
	var buffer bytes.Buffer

	if err := json.NewEncoder(&buffer).Encode(rawSettings); err != nil {
		var zero T
		return zero, fmt.Errorf("encoding settings: %w", err)
	}

	decoder := json.NewDecoder(&buffer)
	decoder.DisallowUnknownFields()

	s := new(T)
	if err := decoder.Decode(s); err != nil {
		var zero T
		return zero, fmt.Errorf("decoding settings: %w", err)
	}

	return *s, nil
}
