package commands

import (
	"os"

	"github.com/frigus02/kyml/pkg/commands/cat"
	"github.com/frigus02/kyml/pkg/commands/completion"
	"github.com/frigus02/kyml/pkg/commands/resolve"
	"github.com/frigus02/kyml/pkg/commands/test"
	"github.com/frigus02/kyml/pkg/commands/tmpl"
	"github.com/frigus02/kyml/pkg/fs"
	"github.com/spf13/cobra"
)

var version = "dev"

// NewRootCommand returns the root command for kyml.
func NewRootCommand() *cobra.Command {
	osFs := fs.NewOSFilesystem()

	c := &cobra.Command{
		Use:          "kyml",
		Short:        "A CLI, which helps you to work with and deploy plain Kubernetes YAML files.",
		SilenceUsage: true,
		Version:      version,
	}

	c.AddCommand(
		cat.NewCmdCat(os.Stdout, osFs),
		completion.NewCmdCompletion(os.Stdout, c),
		resolve.NewCmdResolve(os.Stdin, os.Stdout),
		test.NewCmdTest(os.Stdin, os.Stdout, osFs),
		tmpl.NewCmdTmpl(os.Stdin, os.Stdout),
	)

	return c
}
