package completion

import (
	"fmt"
	"io"

	"github.com/frigus02/cobra"
)

type completionOptions struct {
	shell string
}

// NewCompletionCommand creates a new completion command.
func NewCompletionCommand(out io.Writer, rootCommand *cobra.Command) *cobra.Command {
	var o completionOptions

	var command = &cobra.Command{
		Use:   "completion <shell>",
		Short: "Generate completion scripts for your shell",
		Long: `Write bash, powershell or zsh shell completion code for kyml to stdout.

bash: Ensure you have bash completions installed and enabled. Then output to a file and load it from your .bash_profile.

powershell: Ensure you have a PowerShell profile. Then output to a file and source it from the file provided by "$profile".

zsh: Output to a file in a directory referenced by the $fpath shell variable.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(out, rootCommand)
		},
	}

	return command
}

// Validate validates completion command.
func (o *completionOptions) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please specify a shell (bash, powershell or zsh)")
	}

	o.shell = args[0]
	return nil
}

// Run runs completion command.
func (o *completionOptions) Run(out io.Writer, rootCommand *cobra.Command) error {
	switch o.shell {
	case "bash":
		return rootCommand.GenBashCompletion(out)
	case "powershell":
		return rootCommand.GenPowerShellCompletion(out)
	case "zsh":
		return rootCommand.GenZshCompletion(out)
	default:
		return fmt.Errorf("invalid shell \"%s\" (supported are bash, powershell and zsh)", o.shell)
	}
}
