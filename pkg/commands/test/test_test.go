package test

import (
	"bytes"
	"testing"

	"github.com/frigus02/kyml/pkg/fs"
)

func Test_testOptions_Validate(t *testing.T) {
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
		o := &testOptions{}
		t.Run(tt.name, func(t *testing.T) {
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("testOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

func readFileOrEmpty(filename string, fs fs.Filesystem) string {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(data)
}

func Test_testOptions_Run(t *testing.T) {
	type args struct {
		fs fs.Filesystem
	}
	tests := []struct {
		name         string
		o            *testOptions
		args         args
		wantOut      string
		wantSnapshot string
		wantErr      bool
	}{
		{
			name: "input files don't exist --> error",
			o: &testOptions{
				env1Name:       "staging",
				env1Files:      []string{"testdata/staging/deployment.yml"},
				env2Name:       "production",
				env2Files:      []string{"testdata/production/dep.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				fs: mustCreateFs(t),
			},
			wantOut:      "",
			wantSnapshot: "",
			wantErr:      true,
		},
		{
			name: "snapshot file doesn't exist --> is created",
			o: &testOptions{
				env1Name:       "staging",
				env1Files:      []string{"testdata/staging/deployment.yml", "testdata/staging/service.yml"},
				env2Name:       "production",
				env2Files:      []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				fs: mustCreateFs(t),
			},
			wantOut:      "Wrote initial snapshot diff\n",
			wantSnapshot: "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 3\n",
			wantErr:      false,
		},
		{
			name: "snapshot diff doesn't match --> error",
			o: &testOptions{
				env1Name:       "staging",
				env1Files:      []string{"testdata/staging/deployment.yml", "testdata/staging/service.yml"},
				env2Name:       "production",
				env2Files:      []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				fs: mustCreateFsWithSnapshot(t, "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 2\n"),
			},
			wantOut:      "--- snapshot diff\n+++ this diff\n@@ -5 +5 @@\n-+  replicas: 2\n++  replicas: 3\n",
			wantSnapshot: "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 2\n",
			wantErr:      true,
		},
		{
			name: "snapshot diff doesn't match and update requested --> is updated",
			o: &testOptions{
				env1Name:       "staging",
				env1Files:      []string{"testdata/staging/deployment.yml", "testdata/staging/service.yml"},
				env2Name:       "production",
				env2Files:      []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: true,
			},
			args: args{
				fs: mustCreateFsWithSnapshot(t, "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 2\n"),
			},
			wantOut:      "Updated snapshot diff\n",
			wantSnapshot: "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 3\n",
			wantErr:      false,
		},
		{
			name: "snapshot diff matches",
			o: &testOptions{
				env1Name:       "staging",
				env1Files:      []string{"testdata/staging/deployment.yml", "testdata/staging/service.yml"},
				env2Name:       "production",
				env2Files:      []string{"testdata/production/deployment.yml", "testdata/production/service.yml"},
				snapshotFile:   "kyml-snapshot.diff",
				updateSnapshot: false,
			},
			args: args{
				fs: mustCreateFsWithSnapshot(t, "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 3\n"),
			},
			wantOut:      "",
			wantSnapshot: "--- staging\n+++ production\n@@ -7 +7 @@\n-  replicas: 1\n+  replicas: 3\n",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := tt.o.Run(out, tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("testOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
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
