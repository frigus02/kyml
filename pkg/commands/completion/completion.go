package completion

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		return getPowerShellCompletion(out, rootCommand)
	case "zsh":
		return rootCommand.GenZshCompletion(out)
	default:
		return fmt.Errorf("invalid shell \"%s\" (supported are bash, powershell and zsh)", o.shell)
	}
}

func getPowerShellCompletion(out io.Writer, rootCommand *cobra.Command) error {
	mainTemplate := `using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Register-ArgumentCompleter -Native -CommandName '%s' -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commandElements = $commandAst.CommandElements
    $command = @(
        '%s'
        for ($i = 1; $i -lt $commandElements.Count; $i++) {
            $element = $commandElements[$i]
            if ($element -isnot [StringConstantExpressionAst] -or
                $element.StringConstantType -ne [StringConstantType]::BareWord -or
                $element.Value.StartsWith('-')) {
                break
            }
            $element.Value
        }
    ) -join ';'
    $completions = @(switch ($command) {%s
    })
    $completions.Where{ $_.CompletionText -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}`

	var subCommandCases bytes.Buffer
	generatePowerShellSubcommandCases(&subCommandCases, rootCommand, "")

	fmt.Fprintf(out, mainTemplate, rootCommand.Name(), rootCommand.Name(), subCommandCases.String())
	return nil
}

func generatePowerShellSubcommandCases(out io.Writer, cmd *cobra.Command, previousCommandName string) {
	var cmdName string
	if previousCommandName == "" {
		cmdName = cmd.Name()
	} else {
		cmdName = fmt.Sprintf("%s;%s", previousCommandName, cmd.Name())
	}

	fmt.Fprintf(out, "\n        '%s' {", cmdName)
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden || len(flag.Deprecated) > 0 {
			return
		}
		usage := escapeStringForPowerShell(flag.Usage)
		if len(flag.Shorthand) > 0 {
			fmt.Fprintf(out, "\n            [CompletionResult]::new('-%s', '%s', [CompletionResultType]::ParameterName, '%s')", flag.Shorthand, flag.Shorthand, usage)
		}
		fmt.Fprintf(out, "\n            [CompletionResult]::new('--%s', '%s', [CompletionResultType]::ParameterName, '%s')", flag.Name, flag.Name, usage)
	})
	for _, subCmd := range cmd.Commands() {
		usage := escapeStringForPowerShell(subCmd.Short)
		fmt.Fprintf(out, "\n            [CompletionResult]::new('%s', '%s', [CompletionResultType]::ParameterValue, '%s')", subCmd.Name(), subCmd.Name(), usage)
	}
	fmt.Fprint(out, "\n            break\n        }")

	for _, subCmd := range cmd.Commands() {
		generatePowerShellSubcommandCases(out, subCmd, cmdName)
	}
}

func escapeStringForPowerShell(str string) string {
	return strings.Replace(str, "'", "''", -1)
}
