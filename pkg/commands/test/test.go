package test

import (
	"github.com/spf13/cobra"
)

type testOptions struct {
	filesEnv1 []string
	filesEnv2 []string
}

// NewCmdTest creates a new test command.
func NewCmdTest() *cobra.Command {
	var o testOptions

	cmd := &cobra.Command{
		Use:          "test <FILES_ENV1> <FILES_ENV2>",
		Short:        "Run a snapshot test for a diff between the specified environments.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}

			return o.Run()
		},
	}

	return cmd
}

// Validate validates test command.
func (o *testOptions) Validate(args []string) error {
	return nil
}

// Run runs test command.
func (o *testOptions) Run() error {
	return nil
}
