package cat

import (
	"bytes"
	"testing"

	"github.com/frigus02/kyml/pkg/fs"
)

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

func TestCat(t *testing.T) {
	fs, err := fs.NewFakeFilesystemFromDisk(
		"testdata/base/deployment.yml",
		"testdata/base/service.yml",
		"testdata/overlay-prod/deployment.yml",
	)
	if err != nil {
		t.Errorf("error reading testdata: %v", err)
		return
	}

	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "print files to stdout",
			args: args{
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
			args: args{
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
			if err := Cat(out, tt.args.files, fs); (err != nil) != tt.wantErr {
				t.Errorf("Cat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Cat() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
