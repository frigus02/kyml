package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/frigus02/kyml/pkg/fs"
)

func mustCreateFs(t *testing.T) fs.Filesystem {
	fsWithTestdata, err := fs.NewFakeFilesystemFromDisk(
		"testdata/production/deployment.yml",
		"testdata/production/service.yml",
		"testdata/staging/deployment.yml",
		"testdata/staging/service.yml",
	)
	if err != nil {
		t.Fatalf("error reading testdata: %v", err)
	}

	return fsWithTestdata
}

func mustCreateFsWithSnapshot(t *testing.T, snapshot string) fs.Filesystem {
	fsWithTestdata := mustCreateFs(t)
	err := fsWithTestdata.WriteFile("kyml-snapshot.diff", []byte(snapshot), 0644)
	if err != nil {
		t.Fatalf("error writing snapshot: %v", err)
	}

	return fsWithTestdata
}

func mustCreateStream(t *testing.T, files ...string) io.Reader {
	var content []byte
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

func readFileOrEmpty(filename string, fs fs.Filesystem) string {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(data)
}

func Test_testOptions_Validate(t *testing.T) {
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
		o := &testOptions{}
		t.Run(tt.name, func(t *testing.T) {
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("testOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(o.files, tt.wantFiles) {
				t.Errorf("catOptions.files = %v, want %v", o.files, tt.wantFiles)
			}
		})
	}
}

func Test_testOptions_Run(t *testing.T) {
	type args struct {
		in io.Reader
		fs fs.Filesystem
	}
	tests := []struct {
		name             string
		o                *testOptions
		args             args
		wantOut          string
		wantSnapshot     string
		wantErr          bool
		wantErrToContain string
	}{
		{
			name: "input files don't exist",
			o: &testOptions{
				nameMain:       "staging",
				nameComparison: "production",
				files:          []string{"testdata/production/dep.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				in: mustCreateStream(t, "testdata/staging/deployment.yml"),
				fs: mustCreateFs(t),
			},
			wantOut:      "",
			wantSnapshot: "",
			wantErr:      true,
		},
		{
			name: "snapshot file doesn't exist",
			o: &testOptions{
				nameMain:       "staging",
				nameComparison: "production",
				files:          []string{"testdata/production/deployment.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				in: mustCreateStream(t, "testdata/staging/deployment.yml"),
				fs: mustCreateFs(t),
			},
			wantOut:      "",
			wantSnapshot: "",
			wantErr:      true,
		},
		{
			name: "snapshot diff doesn't match",
			o: &testOptions{
				nameMain:       "staging",
				nameComparison: "production",
				files:          []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				in: mustCreateStream(t, "testdata/staging/deployment.yml", "testdata/staging/service.yml"),
				fs: mustCreateFsWithSnapshot(t, "--- staging\n+++ production\n@@ -19 +19 @@\n-  replicas: 1\n+  replicas: 2\n"),
			},
			wantOut:          "",
			wantSnapshot:     "--- staging\n+++ production\n@@ -19 +19 @@\n-  replicas: 1\n+  replicas: 2\n",
			wantErr:          true,
			wantErrToContain: "--- snapshot diff\n+++ this diff\n@@ -5 +5 @@\n-+  replicas: 2\n++  replicas: 3\n",
		},
		{
			name: "snapshot diff doesn't match and update requested",
			o: &testOptions{
				nameMain:       "staging",
				nameComparison: "production",
				files:          []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: true,
			},
			args: args{
				in: mustCreateStream(t, "testdata/staging/deployment.yml", "testdata/staging/service.yml"),
				fs: mustCreateFsWithSnapshot(t, "--- staging\n+++ production\n@@ -19 +19 @@\n-  replicas: 1\n+  replicas: 2\n"),
			},
			wantOut:      "---\napiVersion: v1\nkind: Service\nmetadata:\n  name: the-service\nspec:\n  ports:\n  - port: 80\n    protocol: TCP\n  selector:\n    deployment: hello\n  type: LoadBalancer\n---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: the-deployment\nspec:\n  replicas: 1\n  template:\n    spec:\n      containers:\n      - image: kyml/hello\n        name: the-container\n",
			wantSnapshot: "--- staging\n+++ production\n@@ -19 +19 @@\n-  replicas: 1\n+  replicas: 3\n",
			wantErr:      false,
		},
		{
			name: "snapshot diff matches",
			o: &testOptions{
				nameMain:       "staging",
				nameComparison: "production",
				files:          []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				in: mustCreateStream(t, "testdata/staging/deployment.yml", "testdata/staging/service.yml"),
				fs: mustCreateFsWithSnapshot(t, "--- staging\n+++ production\n@@ -19 +19 @@\n-  replicas: 1\n+  replicas: 3\n"),
			},
			wantOut:      "---\napiVersion: v1\nkind: Service\nmetadata:\n  name: the-service\nspec:\n  ports:\n  - port: 80\n    protocol: TCP\n  selector:\n    deployment: hello\n  type: LoadBalancer\n---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: the-deployment\nspec:\n  replicas: 1\n  template:\n    spec:\n      containers:\n      - image: kyml/hello\n        name: the-container\n",
			wantSnapshot: "--- staging\n+++ production\n@@ -19 +19 @@\n-  replicas: 1\n+  replicas: 3\n",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := tt.o.Run(tt.args.in, out, tt.args.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("testOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.wantErrToContain) {
				t.Errorf("testOptions.Run() error = %v, wantErrToContain %v", err, tt.wantErrToContain)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("testOptions.Run() = %v, want %v", gotOut, tt.wantOut)
				return
			}
			if gotSnapshot := readFileOrEmpty(tt.o.snapshotFile, tt.args.fs); gotSnapshot != tt.wantSnapshot {
				t.Errorf("testOptions.Run() snapshot = %v, wantSnapshot %v", gotSnapshot, tt.wantSnapshot)
			}
		})
	}
}
