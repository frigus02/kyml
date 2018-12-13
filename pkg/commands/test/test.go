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
	env1Name       string
	env1Files      []string
	env2Name       string
	env2Files      []string
	snapshotFile   string
	updateSnapshot bool
}

// NewCmdTest creates a new test command.
func NewCmdTest(out io.Writer, fs fs.Filesystem) *cobra.Command {
	var o testOptions

	cmd := &cobra.Command{
		Use:          "test",
		Short:        "Run a snapshot test for a diff between the specified environments.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run(out, fs)
		},
	}

	cmd.Flags().StringVar(&o.env1Name, "name1", "1", "Name of first environment")
	cmd.Flags().StringArrayVar(&o.env1Files, "file1", nil, "File in first environment")
	cmd.Flags().StringVar(&o.env2Name, "name2", "2", "Name of second environment")
	cmd.Flags().StringArrayVar(&o.env2Files, "file2", nil, "File in second environment")
	cmd.Flags().StringVarP(&o.snapshotFile, "snapshot", "s", "kyml-snapshot.diff", "Snapshot file")
	cmd.Flags().BoolVarP(&o.updateSnapshot, "update", "u", false, "Update snapshot file")

	return cmd
}

// Validate validates test command.
func (o *testOptions) Validate(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unknown argument")
	}

	return nil
}

// Run runs test command.
func (o *testOptions) Run(out io.Writer, fs fs.Filesystem) error {
	diffStr, err := createEnvDiff(o.env1Name, o.env1Files, o.env2Name, o.env2Files, fs)
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

				fmt.Fprintln(out, "Updated snapshot diff")
			} else {
				fmt.Fprint(out, snapshotDiffStr)
				return fmt.Errorf("snapshot diff does not match this diff")
			}
		}
	} else {
		if err := fs.WriteFile(o.snapshotFile, []byte(diffStr), 0644); err != nil {
			return err
		}

		fmt.Fprintln(out, "Wrote initial snapshot diff")
	}

	return nil
}

func createEnvDiff(
	env1Name string, env1Files []string,
	env2Name string, env2Files []string,
	fs fs.Filesystem,
) (string, error) {
	var buffer1 bytes.Buffer
	if err := cat.Cat(&buffer1, env1Files, fs); err != nil {
		return "", err
	}

	var buffer2 bytes.Buffer
	if err := cat.Cat(&buffer2, env2Files, fs); err != nil {
		return "", err
	}

	return diff.Diff(
		env1Name, string(buffer1.Bytes()),
		env2Name, string(buffer2.Bytes()))
}
