package cat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	cmd.Flags().BoolVarP(&o.deduplicate, "deduplicate", "d", false, "If specified, deduplicate files by name. This means when multiple files have the same name, only the file specified last will be printed.")

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
	var files []string
	if o.deduplicate {
		for _, file := range o.files {
			found := false
			for i, seenFile := range files {
				if filepath.Base(file) == filepath.Base(seenFile) {
					files[i] = file
					found = true
					break
				}
			}

			if !found {
				files = append(files, file)
			}
		}
	} else {
		files = o.files
	}

	var documents []unstructured.Unstructured
	for _, filename := range files {
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

		documents = append(documents, docsInFile...)
	}

	return k8syaml.Encode(out, documents)
}
