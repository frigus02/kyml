package test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/frigus02/kyml/pkg/cat"
	"github.com/frigus02/kyml/pkg/diff"
	"github.com/frigus02/kyml/pkg/fs"
	"github.com/spf13/cobra"
)

type testOptions struct {
	files          []string
	nameFiles      string
	nameStdin      string
	snapshotFile   string
	updateSnapshot bool
}

// NewCmdTest creates a new test command.
func NewCmdTest(in io.Reader, out io.Writer, fs fs.Filesystem) *cobra.Command {
	var o testOptions

	cmd := &cobra.Command{
		Use:          "test <FILE>...",
		Short:        "Run a snapshot test for a diff between the specified environments.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(in, out, fs)
		},
	}

	cmd.Flags().StringVarP(&o.nameFiles, "name-files", "n", "args", "Name of the environment specified in args")
	cmd.Flags().StringVarP(&o.nameStdin, "name-stdin", "o", "stdin", "Name of the environment read from stdin")
	cmd.Flags().StringVarP(&o.snapshotFile, "snapshot-file", "s", "kyml-snapshot.diff", "Snapshot file")
	cmd.Flags().BoolVarP(&o.updateSnapshot, "update", "u", false, "Update snapshot file")

	return cmd
}

// Validate validates test command.
func (o *testOptions) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("specify at least one file")
	}

	o.files = args
	return nil
}

// Run runs test command.
func (o *testOptions) Run(in io.Reader, out io.Writer, fs fs.Filesystem) error {
	var bufferIn bytes.Buffer
	if err := cat.Stream(&bufferIn, in); err != nil {
		return err
	}

	var bufferFiles bytes.Buffer
	if err := cat.Cat(&bufferFiles, o.files, fs); err != nil {
		return err
	}

	diffStr, err := diff.Diff(
		o.nameStdin, bufferIn.String(),
		o.nameFiles, bufferFiles.String())
	if err != nil {
		return err
	}

	snapshotFileInfo, err := fs.Stat(o.snapshotFile)
	if err == nil && snapshotFileInfo.Mode().IsRegular() {
		snapshotBytes, err := fs.ReadFile(o.snapshotFile)
		if err != nil {
			return err
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
				fmt.Fprint(out, snapshotDiffStr)
				return fmt.Errorf("snapshot diff does not match this diff")
			}
		}
	} else {
		if err := fs.WriteFile(o.snapshotFile, []byte(diffStr), 0644); err != nil {
			return err
		}
	}

	fmt.Fprint(out, bufferIn.String())
	return nil
}
