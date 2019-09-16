package effectiveconfig

import (
	"io/ioutil"
	"path/filepath"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/pkg/context"
	yaml "gopkg.in/yaml.v2"
)

// Pipe that writes the effective config file to dist
type Pipe struct {
}

func (Pipe) String() string {
	return "writing effective config file"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) (err error) {
	var path = filepath.Join(ctx.Config.Dist, "config.yaml")
	bts, err := yaml.Marshal(ctx.Config)
	if err != nil {
		return err
	}
	log.WithField("config", path).Info("writing")
	return ioutil.WriteFile(path, bts, 0644)
}
