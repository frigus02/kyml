package tmpl

import (
	"bufio"
	"fmt"
	"html/template"
	"io"

	"github.com/spf13/cobra"
)

type tmplOptions struct {
	values map[string]string
}

// NewCmdTmpl creates a new tmpl command.
func NewCmdTmpl(in io.Reader, out io.Writer) *cobra.Command {
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

			return o.Run(in, out)
		},
	}

	cmd.Flags().StringToStringVarP(&o.values, "value", "v", nil, "Key-value pair, which should be replaced in the YAML files")

	return cmd
}

// Validate validates tmpl command.
func (o *tmplOptions) Validate(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unknown argument")
	}

	return nil
}

// Run runs tmpl command.
func (o *tmplOptions) Run(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		tmpl, err := template.New("").Parse(line)
		if err != nil {
			return err
		}

		err = tmpl.Execute(out, o.values)
		if err != nil {
			return err
		}

		fmt.Fprintln(out)
	}

	return scanner.Err()
}
