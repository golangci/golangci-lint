package before

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Pipe is a global hook pipe
type Pipe struct{}

// String is the name of this pipe
func (Pipe) String() string {
	return "Running before hooks"
}

// Run executes the hooks
func (Pipe) Run(ctx *context.Context) error {
	var tmpl = tmpl.New(ctx)
	/* #nosec */
	for _, step := range ctx.Config.Before.Hooks {
		s, err := tmpl.Apply(step)
		if err != nil {
			return err
		}
		args := strings.Fields(s)
		log.Infof("running %s", color.CyanString(step))
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = ctx.Env.Strings()
		out, err := cmd.CombinedOutput()
		log.Debug(string(out))
		if err != nil {
			return fmt.Errorf("hook failed: %s\n%v", step, string(out))
		}
	}
	return nil
}
