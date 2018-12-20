package tmpl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/frigus02/cobra"
	"github.com/frigus02/kyml/pkg/cat"
)

type tmplOptions struct {
	values  map[string]string
	envVars []string
}

// NewCmdTmpl creates a new tmpl command.
func NewCmdTmpl(in io.Reader, out io.Writer) *cobra.Command {
	var o tmplOptions

	cmd := &cobra.Command{
		Use:   "tmpl",
		Short: "Template Kubernetes YAML files",
		Long: `Template Kubernetes YAML files. Data is read from stdin, executed with the specified context, and printed to stdout.

Templates use the go template syntax (https://golang.org/pkg/text/template/). You can add data to the template context using the options "--value" and "--env". Please note that keys (including environment variable names) are case sensitive in the template.

The command parses the data as Kubernetes YAML documents before templating. While doing this it applies the same transformations as "kyml cat". Since the document is parsed, you need to make sure it is still valid YAML, even with the template characters inside.`,
		Example: `  # Template feature branch files and deploy to cluster
  kyml cat feature/* |
    kyml tmpl \
      -v ImageTag=$(git rev-parse --short HEAD) \
      -e TRAVIS_BRANCH |
    kubectl apply -f -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(in, out)
		},
	}

	cmd.Flags().StringToStringVarP(&o.values, "value", "v", nil, "Add a key-value pair to the template context")
	cmd.Flags().StringArrayVarP(&o.envVars, "env", "e", nil, "Add an environment variable to the template context")

	return cmd
}

// Validate validates tmpl command.
func (o *tmplOptions) Validate(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("this command takes no positional arguments")
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
