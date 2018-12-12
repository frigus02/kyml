package cat

import (
	"fmt"
	"io"
	"os"

	"github.com/frigus02/kyml/pkg/k8syaml"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type catOptions struct {
	files       []string
	deduplicate bool
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

	cmd.Flags().BoolVarP(&o.deduplicate, "deduplicate", "d", false, "If specified, deduplicate YAML documents. This means when multiple YAML documents have the same apiVersion, kind, namespace and name, only the document specified last will be printed.")

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
	var documents []unstructured.Unstructured
	for _, filename := range o.files {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		docsInFile, err := k8syaml.Decode(file)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}

		if o.deduplicate {
			for _, doc := range docsInFile {
				found := false
				for i, seenDoc := range documents {
					if doc.GetAPIVersion() == seenDoc.GetAPIVersion() &&
						doc.GetKind() == seenDoc.GetKind() &&
						doc.GetNamespace() == seenDoc.GetNamespace() &&
						doc.GetName() == doc.GetName() {
						documents[i] = doc
						found = true
						break
					}
				}

				if !found {
					documents = append(documents, doc)
				}
			}
		} else {
			documents = append(documents, docsInFile...)
		}
	}

	return k8syaml.Encode(out, documents)
}
