package tmpl

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_tmplOptions_Validate(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no args",
			args: args{
				args: []string{},
			},
			wantErr: false,
		},
		{
			name: "error when args are specified",
			args: args{
				args: []string{"hello"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &tmplOptions{}
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("tmplOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tmplOptions_Run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		o       *tmplOptions
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "success",
			o: &tmplOptions{
				values: map[string]string{
					"branch": "my-feature",
					"tag":    "latest",
				},
			},
			args: args{
				strings.NewReader("label: {{.branch}}\nimage: hello:{{.tag}}\n"),
			},
			wantOut: "label: my-feature\nimage: hello:latest\n",
			wantErr: false,
		},
		{
			name: "invalid template",
			o: &tmplOptions{
				values: map[string]string{
					"branch": "my-feature",
					"tag":    "latest",
				},
			},
			args: args{
				strings.NewReader("label: {{.branch\nimage: hello:{{.tag}}\n"),
			},
			wantOut: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := tt.o.Run(tt.args.in, out); (err != nil) != tt.wantErr {
				t.Errorf("tmplOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("tmplOptions.Run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
