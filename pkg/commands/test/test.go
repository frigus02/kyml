package test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/frigus02/kyml/pkg/cat"
	"github.com/frigus02/kyml/pkg/diff"
	"github.com/frigus02/kyml/pkg/fs"
	"github.com/spf13/cobra"
)

type testOptions struct {
	files          []string
	nameComparison string
	nameMain       string
	snapshotFile   string
	updateSnapshot bool
}

// NewCmdTest creates a new test command.
func NewCmdTest(in io.Reader, out io.Writer, fs fs.Filesystem) *cobra.Command {
	var o testOptions

	cmd := &cobra.Command{
		Use:   "test <file>...",
		Short: "Run a snapshot test on the diff between Kubernetes YAML files of two environments",
		Long: `Run a snapshot test on the diff between Kubernetes YAML files of two environments.

Integrate this in your CI builds to make sure your environments don't accidentially drift apart.

The main environment is specified on stdin. Use "kyml cat" to concatenate multiple files and pipe the result into "kyml test".

The comparison environment is specified using filenames. Files are concatenated using the same rules as in "kyml cat".

The command compares the diff between these environments to a previous diff stored in the specified snapshot file. If it matches, it prints the main environment to stdout, so it can be piped into followup commands like "kyml tmpl" or "kubectl apply". If it doesn't match, it prints the diff to stderr and exits with a non-zero exit code.`,
		Example: `  # Make sure production and staging don't drift apart unknowingly
  kyml cat production/* | kyml test staging/* \
    --name-main production \
    --name-staging staging \
    --snapshot-file tests/prod-vs-staging.diff

  # Update snapshot file when the change was deliberate
  kyml cat production/* | kyml test staging/* \
    --name-main production \
    --name-staging staging \
    --snapshot-file tests/prod-vs-staging.diff \
    --update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(in, out, fs)
		},
	}

	cmd.Flags().StringVar(&o.nameMain, "name-main", "main", "Name of the main environment read from stdin")
	cmd.Flags().StringVar(&o.nameComparison, "name-comparison", "comparison", "Name of the comparison environment read from files")
	cmd.Flags().StringVarP(&o.snapshotFile, "snapshot-file", "s", "kyml-snapshot.diff", "Snapshot file")
	cmd.Flags().BoolVarP(&o.updateSnapshot, "update", "u", false, "If specified, update snapshot file and exit successfully in case of non-match")

	return cmd
}

// Validate validates test command.
func (o *testOptions) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("specify at least one file for the comparison environment")
	}

	o.files = args
	return nil
}

// Run runs test command.
func (o *testOptions) Run(in io.Reader, out io.Writer, fs fs.Filesystem) error {
	var bufferMain bytes.Buffer
	if err := cat.Stream(&bufferMain, in); err != nil {
		return err
	}

	var bufferComparison bytes.Buffer
	if err := cat.Cat(&bufferComparison, o.files, fs); err != nil {
		return err
	}

	diffStr, err := diff.Diff(
		o.nameMain, bufferMain.String(),
		o.nameComparison, bufferComparison.String())
	if err != nil {
		return err
	}

	snapshotBytes, err := fs.ReadFile(o.snapshotFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot open snapshot file: %v", err)
		}

		if !o.updateSnapshot {
			return fmt.Errorf("snapshot file does not exist")
		}
	}

	snapshotDiffStr, err := diff.Diff(
		"snapshot diff", string(snapshotBytes),
		"this diff", diffStr)
	if err != nil {
		return err
	}

	if snapshotDiffStr != "" {
		if o.updateSnapshot {
			if err := fs.WriteFile(o.snapshotFile, []byte(diffStr), 0644); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("snapshot diff does not match this diff\n%s", snapshotDiffStr)
		}
	}

	fmt.Fprint(out, bufferMain.String())
	return nil
}
