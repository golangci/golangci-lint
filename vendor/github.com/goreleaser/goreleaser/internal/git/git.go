// Package git provides an integration with the git command
package git

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/apex/log"
)

// IsRepo returns true if current folder is a git repository
func IsRepo() bool {
	out, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// Run runs a git command and returns its output or errors
func Run(args ...string) (string, error) {
	// TODO: use exex.CommandContext here and refactor.
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}
	args = append(extraArgs, args...)
	/* #nosec */
	var cmd = exec.Command("git", args...)
	log.WithField("args", args).Debug("running git")
	bts, err := cmd.CombinedOutput()
	log.WithField("output", string(bts)).
		Debug("git result")
	if err != nil {
		return "", errors.New(string(bts))
	}
	return string(bts), nil
}

// Clean the output
func Clean(output string, err error) (string, error) {
	output = strings.Replace(strings.Split(output, "\n")[0], "'", "", -1)
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
