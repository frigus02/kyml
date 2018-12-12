package cat

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func Test_catOptions_Validate(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "error if no args",
			args: args{
				args: []string{},
			},
			wantErr:   true,
			wantFiles: nil,
		},
		{
			name: "files set to args",
			args: args{
				args: []string{"foo", "bar", "baz"},
			},
			wantErr:   false,
			wantFiles: []string{"foo", "bar", "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &catOptions{}
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("catOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(o.files, tt.wantFiles) {
				t.Errorf("catOptions.files = %v, want %v", o.files, tt.wantFiles)
			}
		})
	}
}

var testDataManifests = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  replicas: 3
  template:
    spec:
      containers:
      - image: monopole/hello
        name: the-container
---
apiVersion: v1
kind: Service
metadata:
  name: the-service
spec:
  ports:
  - port: 80
    protocol: TCP
  selector:
    deployment: hello
  type: LoadBalancer
`

func Test_catOptions_Run(t *testing.T) {
	tests := []struct {
		name    string
		o       *catOptions
		wantOut string
		wantErr bool
	}{
		{
			name: "print files to stdout",
			o: &catOptions{
				files: []string{
					"testdata/base/deployment.yml",
					"testdata/base/service.yml",
					"testdata/overlay-prod/deployment.yml",
				},
			},
			wantOut: testDataManifests,
			wantErr: false,
		},
		{
			name: "file does not exist",
			o: &catOptions{
				files: []string{
					"testdata/something.yml",
				},
			},
			wantOut: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := tt.o.Run(out); (err != nil) != tt.wantErr {
				t.Errorf("catOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := strings.Replace(out.String(), "\r\n", "\n", -1); gotOut != tt.wantOut {
				t.Errorf("catOptions.Run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
