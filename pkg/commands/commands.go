package commands

import (
	"os"

	"github.com/frigus02/kyml/pkg/commands/cat"
	"github.com/frigus02/kyml/pkg/commands/completion"
	"github.com/frigus02/kyml/pkg/commands/test"
	"github.com/frigus02/kyml/pkg/commands/tmpl"
	"github.com/frigus02/kyml/pkg/fs"
	"github.com/frigus02/cobra"
)

var version = "dev"

// NewRootCommand returns the root command for kyml.
func NewRootCommand() *cobra.Command {
	osFs := fs.NewOSFilesystem()

	c := &cobra.Command{
		Use:          "kyml",
		Short:        "kyml helps you to manage your Kubernetes YAML files.",
		SilenceUsage: true,
		Version:      version,
	}

	c.AddCommand(
		cat.NewCmdCat(os.Stdout, osFs),
		completion.NewCompletionCommand(os.Stdout, c),
		test.NewCmdTest(os.Stdin, os.Stdout, osFs),
		tmpl.NewCmdTmpl(os.Stdin, os.Stdout),
	)

	return c
}
