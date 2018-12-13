package tmpl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/frigus02/kyml/pkg/cat"
	"github.com/spf13/cobra"
)

type tmplOptions struct {
	values  map[string]string
	envVars []string
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
	cmd.Flags().StringArrayVarP(&o.envVars, "env", "e", nil, "Environment variable, which should be replaced in the YAML files")

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
	vars := make(map[string]string)
	for key, value := range o.values {
		vars[key] = value
	}
	for _, env := range o.envVars {
		vars[env] = os.Getenv(env)
	}

	var buffer bytes.Buffer
	if err := cat.Stream(&buffer, in); err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(buffer.String())
	if err != nil {
		return err
	}

	if err = tmpl.Execute(out, vars); err != nil {
		return err
	}

	fmt.Fprint(out)
	return nil
}
