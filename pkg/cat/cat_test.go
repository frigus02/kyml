package cat

import (
	"bytes"
	"io"
	"io/ioutil"
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

func mustCreateFs(t *testing.T) fs.Filesystem {
	fsWithTestdata, err := fs.NewFakeFilesystemFromDisk(
		"testdata/base/deployment.yml",
		"testdata/base/service.yml",
		"testdata/overlay-prod/deployment.yml",
	)
	if err != nil {
		t.Fatalf("error reading testdata: %v", err)
	}

	return fsWithTestdata
}

func mustCreateStream(t *testing.T) io.Reader {
	var content []byte
	files := []string{
		"testdata/base/deployment.yml",
		"testdata/base/service.yml",
		"testdata/overlay-prod/deployment.yml",
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("error reading testdata: %v", err)
		}

		content = append(content, []byte("\n---\n")...)
		content = append(content, data...)
	}

	return bytes.NewReader(content)
}

func TestCat(t *testing.T) {
	type args struct {
		files []string
		fs    fs.Filesystem
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "print deduplicated docs from files",
			args: args{
				files: []string{
					"testdata/base/deployment.yml",
					"testdata/base/service.yml",
					"testdata/overlay-prod/deployment.yml",
				},
				fs: mustCreateFs(t),
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
				fs: mustCreateFs(t),
			},
			wantOut: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := Cat(out, tt.args.files, tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("Cat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Cat() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestStream(t *testing.T) {
	type args struct {
		stream io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "print deduplicated docs from stream",
			args: args{
				stream: mustCreateStream(t),
			},
			wantOut: testDataManifests,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := Stream(out, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("Stream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Stream() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
