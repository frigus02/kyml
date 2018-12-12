package cat

import (
	"fmt"
	"io"

	"github.com/frigus02/kyml/pkg/cat"
	"github.com/spf13/cobra"
)

type catOptions struct {
	files []string
}

// NewCmdCat creates a new cat command.
func NewCmdCat(out io.Writer) *cobra.Command {
	var o catOptions

	cmd := &cobra.Command{
		Use:          "cat <FILE>...",
		Short:        "Concatenate Kubernetes YAML files to standard output.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(out)
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
func (o *catOptions) Run(out io.Writer) error {
	return cat.Cat(out, o.files)
}
