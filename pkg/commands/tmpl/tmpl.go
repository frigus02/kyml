package tmpl

import (
	"github.com/spf13/cobra"
)

type tmplOptions struct {
	values map[string]string
}

// NewCmdTmpl creates a new tmpl command.
func NewCmdTmpl() *cobra.Command {
	var o tmplOptions

	cmd := &cobra.Command{
		Use:          "tmpl",
		Short:        "Template Kubernetes YAML files. File is read from stdin, executed with the specified values, and printed to stdout.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run()
		},
	}

	cmd.Flags().StringToStringVarP(&o.values, "output", "o", nil, "Key-value pair, which should be replaced in the YAML files")

	return cmd
}

// Validate validates tmpl command.
func (o *tmplOptions) Validate(args []string) error {
	return nil
}

// Run runs tmpl command.
func (o *tmplOptions) Run() error {
	return nil
}
