package executors

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/golangci/golangci-shared/pkg/timeutils"
)

type Docker struct {
	mountFromPath, mountToPath string
	image                      string
	wd                         string
	*envStore
}

func NewDocker(mountFromPath, mountToPath string) *Docker {
	return &Docker{
		mountFromPath: mountFromPath,
		mountToPath:   mountToPath,
		image:         "golangci_executor",
		envStore:      newEnvStoreNoOS(),
	}
}

var _ Executor = Docker{}

func (d Docker) Run(ctx context.Context, name string, args ...string) (string, error) {
	// XXX: don't use docker sdk because it's too heavyweight: dep ensure takes minutes on it
	dockerArgs := []string{
		"run",
		"-v", fmt.Sprintf("%s:%s", d.mountFromPath, d.mountToPath),
	}

	for _, e := range d.env {
		dockerArgs = append(dockerArgs, "-e", e)
	}

	dockerArgs = append(dockerArgs, "--rm", d.image, name)
	dockerArgs = append(dockerArgs, args...)
	// TODO: take work dir d.wd into account

	defer timeutils.Track(time.Now(), "docker full execution: docker %v", dockerArgs)

	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (d Docker) WithEnv(k, v string) Executor {
	dCopy := d
	dCopy.SetEnv(k, v)
	return dCopy
}

func (d Docker) Clean() {}

func (d Docker) WithWorkDir(wd string) Executor {
	dCopy := d
	dCopy.wd = wd
	return dCopy
}

func (d Docker) WorkDir() string {
	panic("isn't supported") // TODO
}
