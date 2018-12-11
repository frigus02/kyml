package commands

import (
	"os"

	"github.com/frigus02/kyml/pkg/commands/cat"
	"github.com/frigus02/kyml/pkg/commands/edit"
	"github.com/frigus02/kyml/pkg/commands/test"
	"github.com/frigus02/kyml/pkg/commands/tmpl"
	"github.com/frigus02/kyml/pkg/commands/version"
	"github.com/spf13/cobra"
)

// NewDefaultCommand returns the default (aka root) command for kyml command.
func NewDefaultCommand() *cobra.Command {
	stdOut := os.Stdout

	c := &cobra.Command{
		Use:   "kyml",
		Short: "kyml helps you to manage your Kubernetes YAML files.",
	}

	c.AddCommand(
		cat.NewCmdCat(stdOut),
		edit.NewCmdEdit(),
		test.NewCmdTest(),
		tmpl.NewCmdTmpl(),
		version.NewCmdVersion(stdOut),
	)

	return c
}
