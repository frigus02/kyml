package edit

import (
	"github.com/spf13/cobra"
)

type editOptions struct {
	set map[string]string
	del []string
}

// NewCmdEdit creates a new edit command.
func NewCmdEdit() *cobra.Command {
	var o editOptions

	cmd := &cobra.Command{
		Use:          "edit",
		Short:        "Edit Kubernetes YAML files. File is read from stdin, edited based on the specified options, and printed to stdout.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run()
		},
	}

	cmd.Flags().StringToStringVarP(&o.set, "set", "s", nil, "Set the specified path to the specified value")
	cmd.Flags().StringArrayVarP(&o.del, "delete", "d", nil, "Delete the specified path")

	return cmd
}

// Validate validates edit command.
func (o *editOptions) Validate(args []string) error {
	return nil
}

// Run runs edit command.
func (o *editOptions) Run() error {
	return nil
}
