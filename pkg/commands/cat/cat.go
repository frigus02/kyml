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
		Short: "Concatenate Kubernetes YAML files to standard output.",
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
