package completion

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func Test_completionOptions_Validate(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantShell string
	}{
		{
			name: "error if no args",
			args: args{
				args: []string{},
			},
			wantErr:   true,
			wantShell: "",
		},
		{
			name: "error if too many args",
			args: args{
				args: []string{"zsh", "bash"},
			},
			wantErr:   true,
			wantShell: "",
		},
		{
			name: "shell set to first arg",
			args: args{
				args: []string{"zsh"},
			},
			wantErr:   false,
			wantShell: "zsh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &completionOptions{}
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("completionOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if o.shell != tt.wantShell {
				t.Errorf("completionOptions.shell = %v, want %v", o.shell, tt.wantShell)
			}
		})
	}
}

func Test_completionOptions_Run(t *testing.T) {
	cmd := &cobra.Command{
		Use: "kyml",
	}

	type args struct {
		rootCommand *cobra.Command
	}
	tests := []struct {
		name             string
		o                *completionOptions
		args             args
		wantOutToContain string
		wantErr          bool
	}{
		{
			name:             "bash",
			o:                &completionOptions{"bash"},
			args:             args{cmd},
			wantOutToContain: "# bash completion for kyml",
			wantErr:          false,
		},
		{
			name:             "powershell",
			o:                &completionOptions{"powershell"},
			args:             args{cmd},
			wantOutToContain: "Register-ArgumentCompleter -CommandName 'kyml'",
			wantErr:          false,
		},
		{
			name:             "zsh",
			o:                &completionOptions{"zsh"},
			args:             args{cmd},
			wantOutToContain: "#compdef _kyml kyml",
			wantErr:          false,
		},
		{
			name:             "invalid shell",
			o:                &completionOptions{"invalid"},
			args:             args{cmd},
			wantOutToContain: "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := tt.o.Run(out, tt.args.rootCommand); (err != nil) != tt.wantErr {
				t.Errorf("completionOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); !strings.Contains(gotOut, tt.wantOutToContain) {
				t.Errorf("completionOptions.Run() = %v, want it to contain %v", gotOut, tt.wantOutToContain)
			}
		})
	}
}
