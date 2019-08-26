package cat

import (
	"fmt"
	"io"

	"github.com/frigus02/kyml/pkg/cat"
	"github.com/frigus02/kyml/pkg/fs"
	"github.com/spf13/cobra"
)

type catOptions struct {
	files []string
}

// NewCmdCat creates a new cat command.
func NewCmdCat(out io.Writer, fs fs.Filesystem) *cobra.Command {
	var o catOptions

	cmd := &cobra.Command{
		Use:   "cat <file>...",
		Short: "Concatenate Kubernetes YAML files to stdout",
		Long: `Read and concatenate YAML documents from all files in the order they are specified. Then print them to stdout.

YAML documents are changed in the following ways:
- Documents are parsed as Kubernetes YAML documents and then formatted. This will change indentation and ordering of properties.
- Documents are deduplicated. If multiple YAML documents refer to the same Kubernetes resource, only the last one will appear in the result.
- Documents are sorted by dependencies, e.g. namespaces come before deployments.

The result of this command can be piped into other commands like "kyml test" or "kubectl apply".`,
		Example: `  # Cat one folder
  kyml cat production/*

  # Merge YAML documents from two folders
  kyml cat base/* overlay-production/*

  # Specify files individually
  kyml cat prod/deployment.yaml prod/service.yaml prod/ingress.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(out, fs)
		},
	}

	return cmd
}

// Validate validates cat command.
func (o *catOptions) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("specify at least one file")
	}

	o.files = args
	return nil
}

// Run runs cat command.
func (o *catOptions) Run(out io.Writer, fs fs.Filesystem) error {
	return cat.Cat(out, o.files, fs)
}
