package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// NewCmdVersion makes a new version command.
func NewCmdVersion(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print kyml version",
		Example: `kyml version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(out, "%v, commit %v, built at %v", version, commit, date)
		},
	}
}
