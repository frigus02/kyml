package commands

import (
	"os"

	"github.com/frigus02/kyml/pkg/commands/cat"
	"github.com/frigus02/kyml/pkg/commands/test"
	"github.com/frigus02/kyml/pkg/commands/tmpl"
	"github.com/frigus02/kyml/pkg/commands/version"
	"github.com/frigus02/kyml/pkg/fs"
	"github.com/spf13/cobra"
)

// NewDefaultCommand returns the default (aka root) command for kyml.
func NewDefaultCommand() *cobra.Command {
	osFs := fs.NewOSFilesystem()

	c := &cobra.Command{
		Use:   "kyml",
		Short: "kyml helps you to manage your Kubernetes YAML files.",
	}

	c.AddCommand(
		cat.NewCmdCat(os.Stdout, osFs),
		test.NewCmdTest(os.Stdin, os.Stdout, osFs),
		tmpl.NewCmdTmpl(os.Stdin, os.Stdout),
		version.NewCmdVersion(os.Stdout),
	)

	return c
}
