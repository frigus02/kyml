package tmpl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/frigus02/cobra"
	"github.com/frigus02/kyml/pkg/cat"
	"github.com/frigus02/kyml/pkg/k8syaml"
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

Templates are only supported in values of type string. They use the go template syntax (https://golang.org/pkg/text/template/). You can add data to the template context using the options "--value" and "--env". Please note that keys (including environment variable names) are case sensitive.

The command parses the data as Kubernetes YAML documents before templating. While doing so it applies the same transformations as "kyml cat". Since the document is parsed, you need to make sure it is still valid YAML, even with the template characters inside.`,
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

	documents, err := cat.StreamDecodeOnly(in)
	if err != nil {
		return err
	}

	var execTmpl = func(text, name string) (string, error) {
		tmpl, err := template.New(name).Option("missingkey=error").Parse(text)
		if err != nil {
			return "", err
		}

		var result bytes.Buffer
		if err = tmpl.Execute(&result, vars); err != nil {
			return "", err
		}

		return result.String(), nil
	}

	for _, doc := range documents {
		templated, err := templateValuesInMap(doc.UnstructuredContent(), execTmpl)
		if err != nil {
			return err
		}

		doc.SetUnstructuredContent(templated)
	}

	return k8syaml.Encode(out, documents)
}

type valueTemplater func(text, name string) (string, error)

func templateValuesInMap(m map[string]interface{}, execTmpl valueTemplater) (map[string]interface{}, error) {
	newMap := make(map[string]interface{}, len(m))
	for key, value := range m {
		templated, err := templateValue(value, key, execTmpl)
		if err != nil {
			return nil, err
		}

		newMap[key] = templated
	}

	return newMap, nil
}

func templateValuesInSlice(s []interface{}, name string, execTmpl valueTemplater) ([]interface{}, error) {
	newSlice := make([]interface{}, len(s))
	for index, value := range s {
		templated, err := templateValue(value, fmt.Sprintf("%s[%d]", name, index), execTmpl)
		if err != nil {
			return nil, err
		}

		newSlice[index] = templated
	}

	return newSlice, nil
}

func templateValue(value interface{}, name string, execTmpl valueTemplater) (interface{}, error) {
	switch value := value.(type) {
	case map[string]interface{}:
		return templateValuesInMap(value, execTmpl)
	case []interface{}:
		return templateValuesInSlice(value, name, execTmpl)
	case string:
		return execTmpl(value, name)
	default:
		return value, nil
	}
}
