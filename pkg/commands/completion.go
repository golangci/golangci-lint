package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (e *Executor) initCompletion() {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		RunE:  e.executeCompletion,
	}
	e.rootCmd.AddCommand(completionCmd)
}

func (e *Executor) executeCompletion(cmd *cobra.Command, args []string) error {
	err := cmd.Root().GenBashCompletion(os.Stdout)
	if err != nil {
		return fmt.Errorf("unable to generate bash completions: %v", err)
	}

	return nil
}
