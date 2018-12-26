package resolve

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func Test_resolveWithDockerInspect(t *testing.T) {
	type args struct {
		imageRef string
	}
	tests := []struct {
		name    string
		args    args
		cmdOut  []byte
		cmdErr  error
		want    string
		wantCmd string
		wantErr bool
	}{
		{
			name:   "image not found",
			args:   args{"registry:5000/path/hello:latest"},
			cmdOut: []byte(""),
			cmdErr: &exec.ExitError{
				Stderr: []byte("\nError: No such object: registry:5000/path/hello:latest\n"),
			},
			want:    "",
			wantCmd: "docker inspect --format {{json .RepoDigests}} registry:5000/path/hello:latest",
			wantErr: false,
		},
		{
			name:   "Docker daemon not running",
			args:   args{"registry:5000/path/hello:latest"},
			cmdOut: []byte(""),
			cmdErr: &exec.ExitError{
				Stderr: []byte("\nError response from daemon: Bad response from Docker engine\n"),
			},
			want:    "",
			wantCmd: "docker inspect --format {{json .RepoDigests}} registry:5000/path/hello:latest",
			wantErr: true,
		},
		{
			name:    "docker command not found",
			args:    args{"registry:5000/path/hello:latest"},
			cmdOut:  []byte(""),
			cmdErr:  errors.New("docker inspect: exec: \"docker\": executable file not found in $PATH"),
			want:    "",
			wantCmd: "docker inspect --format {{json .RepoDigests}} registry:5000/path/hello:latest",
			wantErr: true,
		},
		{
			name:    "success",
			args:    args{"registry:5000/path/hello:latest"},
			cmdOut:  []byte("[\"another-registry.example.com/hello@sha256:2d8b22d01ca51eef988ff3ae8dcf37c182553b662ea47d3d62ce8208a3b83aef\",\"registry:5000/path/hello@sha256:2d8b22d01ca51eef988ff3ae8dcf37c182553b662ea47d3d62ce8208a3b83aef\"]\n"),
			cmdErr:  nil,
			want:    "registry:5000/path/hello@sha256:2d8b22d01ca51eef988ff3ae8dcf37c182553b662ea47d3d62ce8208a3b83aef",
			wantCmd: "docker inspect --format {{json .RepoDigests}} registry:5000/path/hello:latest",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotCmd string
			execCmdMock := func(name string, arg ...string) ([]byte, error) {
				gotCmd = name + " " + strings.Join(arg, " ")
				return tt.cmdOut, tt.cmdErr
			}

			got, err := resolveWithDockerInspect(tt.args.imageRef, execCmdMock)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveWithDockerInspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("resolveWithDockerInspect() cmd = %v, wantCmd %v", gotCmd, tt.wantCmd)
				return
			}
			if got != tt.want {
				t.Errorf("resolveWithDockerInspect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolveWithDockerManifestInspect(t *testing.T) {
	type args struct {
		imageRef string
	}
	tests := []struct {
		name    string
		args    args
		cmdOut  []byte
		cmdErr  error
		want    string
		wantCmd string
		wantErr bool
	}{
		{
			name:   "image not found",
			args:   args{"registry:5000/path/hello:latest"},
			cmdOut: []byte(""),
			cmdErr: &exec.ExitError{
				Stderr: []byte("no such manifest: registry:5000/path/hello:latest\n"),
			},
			want:    "",
			wantCmd: "docker manifest inspect --verbose registry:5000/path/hello:latest",
			wantErr: false,
		},
		{
			name:   "experimental cli features not enabled",
			args:   args{"registry:5000/path/hello:latest"},
			cmdOut: []byte(""),
			cmdErr: &exec.ExitError{
				Stderr: []byte("docker manifest inspect is only supported on a Docker cli with experimental cli features enabled\n"),
			},
			want:    "",
			wantCmd: "docker manifest inspect --verbose registry:5000/path/hello:latest",
			wantErr: true,
		},
		{
			name:    "docker command not found",
			args:    args{"registry:5000/path/hello:latest"},
			cmdOut:  []byte(""),
			cmdErr:  errors.New("docker inspect: exec: \"docker\": executable file not found in $PATH"),
			want:    "",
			wantCmd: "docker manifest inspect --verbose registry:5000/path/hello:latest",
			wantErr: true,
		},
		{
			name: "manifest list does not contain linux amd64",
			args: args{"openjdk:latest"},
			cmdOut: []byte(`[
				{
					"Ref": "docker.io/library/openjdk:latest@sha256:ff3da04131714a6e03d02684a33a3858e622923344534de87ff453d03181337a",
					"Descriptor": {
						"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
						"digest": "sha256:ff3da04131714a6e03d02684a33a3858e622923344534de87ff453d03181337a",
						"size": 2000,
						"platform": {
							"architecture": "arm",
							"os": "linux",
							"variant": "v5"
						}
					}
				}
			]`),
			cmdErr:  nil,
			want:    "",
			wantCmd: "docker manifest inspect --verbose openjdk:latest",
			wantErr: false,
		},
		{
			name: "success with manifest list",
			args: args{"openjdk:latest"},
			cmdOut: []byte(`[
				{
					"Ref": "docker.io/library/openjdk:latest@sha256:c7381bfd53670f1211314885b03b98f5e13fddf6958afeec61092b07c56ddef1",
					"Descriptor": {
						"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
						"digest": "sha256:c7381bfd53670f1211314885b03b98f5e13fddf6958afeec61092b07c56ddef1",
						"size": 2000,
						"platform": {
							"architecture": "amd64",
							"os": "linux"
						}
					}
				},
				{
					"Ref": "docker.io/library/openjdk:latest@sha256:ff3da04131714a6e03d02684a33a3858e622923344534de87ff453d03181337a",
					"Descriptor": {
						"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
						"digest": "sha256:ff3da04131714a6e03d02684a33a3858e622923344534de87ff453d03181337a",
						"size": 2000,
						"platform": {
							"architecture": "arm",
							"os": "linux",
							"variant": "v5"
						}
					}
				}
			]`),
			cmdErr:  nil,
			want:    "openjdk@sha256:c7381bfd53670f1211314885b03b98f5e13fddf6958afeec61092b07c56ddef1",
			wantCmd: "docker manifest inspect --verbose openjdk:latest",
			wantErr: false,
		},
		{
			name: "success",
			args: args{"registry:5000/path/hello:latest"},
			cmdOut: []byte(`{
				"Ref": "registry:5000/path/hello:latest",
				"Descriptor": {
					"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
					"digest": "sha256:2d8b22d01ca51eef988ff3ae8dcf37c182553b662ea47d3d62ce8208a3b83aef",
					"size": 3661,
					"platform": {
						"architecture": "amd64",
						"os": "linux"
					}
				}
			}`),
			cmdErr:  nil,
			want:    "registry:5000/path/hello@sha256:2d8b22d01ca51eef988ff3ae8dcf37c182553b662ea47d3d62ce8208a3b83aef",
			wantCmd: "docker manifest inspect --verbose registry:5000/path/hello:latest",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotCmd string
			execCmdMock := func(name string, arg ...string) ([]byte, error) {
				gotCmd = name + " " + strings.Join(arg, " ")
				return tt.cmdOut, tt.cmdErr
			}

			got, err := resolveWithDockerManifestInspect(tt.args.imageRef, execCmdMock)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveWithDockerManifestInspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("resolveWithDockerInspect() cmd = %v, wantCmd %v", gotCmd, tt.wantCmd)
				return
			}
			if got != tt.want {
				t.Errorf("resolveWithDockerManifestInspect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeTagAndDigest(t *testing.T) {
	type args struct {
		imageRef string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "only name",
			args: args{"hello"},
			want: "hello",
		},
		{
			name: "only name with domain and path",
			args: args{"registry:5000/path/hello"},
			want: "registry:5000/path/hello",
		},
		{
			name: "tag",
			args: args{"hello:latest"},
			want: "hello",
		},
		{
			name: "tag with domain and path",
			args: args{"registry:5000/path/hello:latest"},
			want: "registry:5000/path/hello",
		},
		{
			name: "digest",
			args: args{"hello@sha256:e3227b2d3d50d02fb80194e6b82ca2df1ece0608e2c1fa9e73356775c6a7c095"},
			want: "hello",
		},
		{
			name: "digest with domain and path",
			args: args{"registry:5000/path/hello@sha256:e3227b2d3d50d02fb80194e6b82ca2df1ece0608e2c1fa9e73356775c6a7c095"},
			want: "registry:5000/path/hello",
		},
		{
			name: "tag and digest",
			args: args{"hello:latest@sha256:e3227b2d3d50d02fb80194e6b82ca2df1ece0608e2c1fa9e73356775c6a7c095"},
			want: "hello",
		},
		{
			name: "tag and digest with domain and path",
			args: args{"registry:5000/path/hello:latest@sha256:e3227b2d3d50d02fb80194e6b82ca2df1ece0608e2c1fa9e73356775c6a7c095"},
			want: "registry:5000/path/hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeTagAndDigest(tt.args.imageRef); got != tt.want {
				t.Errorf("removeTagAndDigest() = %v, want %v", got, tt.want)
			}
		})
	}
}
