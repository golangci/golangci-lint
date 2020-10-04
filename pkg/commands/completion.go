package commands

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func (e *Executor) initCompletion() {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Output completion script",
	}
	e.rootCmd.AddCommand(completionCmd)

	bashCmd := &cobra.Command{
		Use:   "bash",
		Short: "Output bash completion script",
		RunE:  e.executeBashCompletion,
	}
	completionCmd.AddCommand(bashCmd)

	zshCmd := &cobra.Command{
		Use:   "zsh",
		Short: "Output zsh completion script",
		RunE:  e.executeZshCompletion,
	}
	completionCmd.AddCommand(zshCmd)

	fishCmd := &cobra.Command{
		Use:   "fish",
		Short: "Output fish completion script",
		RunE:  e.executeFishCompletion,
	}
	completionCmd.AddCommand(fishCmd)
}

func (e *Executor) executeBashCompletion(cmd *cobra.Command, args []string) error {
	err := cmd.Root().GenBashCompletion(os.Stdout)
	if err != nil {
		return errors.Wrap(err, "unable to generate bash completions: %v")
	}

	return nil
}

func (e *Executor) executeZshCompletion(cmd *cobra.Command, args []string) error {
	err := cmd.Root().GenZshCompletion(os.Stdout)
	if err != nil {
		return errors.Wrap(err, "unable to generate zsh completions: %v")
	}
	// Add extra compdef directive to support sourcing command directly.
	// https://github.com/spf13/cobra/issues/881
	// https://github.com/spf13/cobra/pull/887
	fmt.Println("compdef _golangci-lint golangci-lint")

	return nil
}

func (e *Executor) executeFishCompletion(cmd *cobra.Command, args []string) error {
	err := cmd.Root().GenFishCompletion(os.Stdout, true)
	if err != nil {
		return errors.Wrap(err, "generate fish completion")
	}

	return nil
}
