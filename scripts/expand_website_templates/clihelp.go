package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func getHelpLinters() ([][]byte, error) {
	lintersOut, err := exec.Command("./golangci-lint", "help", "linters").Output()
	if err != nil {
		return nil, fmt.Errorf("can't run linters cmd: %w", err)
	}

	return bytes.Split(lintersOut, []byte("\n\n")), nil
}

func getHelpRun() ([]byte, error) {
	helpCmd := exec.Command("./golangci-lint", "run", "-h")
	helpCmd.Env = append(helpCmd.Env, os.Environ()...)
	helpCmd.Env = append(helpCmd.Env, "HELP_RUN=1") // make default concurrency stable: don't depend on machine CPU number
	help, err := helpCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't run help cmd: %w", err)
	}

	helpLines := bytes.Split(help, []byte("\n"))

	return bytes.Join(helpLines[2:], []byte("\n")), nil
}
