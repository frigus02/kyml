package resolve

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

var testManifestService = `---
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

var testManifestDeploymentNoContainers = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  template:
    spec:
      serviceAccountName: john
`

var testManifestDeploymentMalformedContainer = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  template:
    spec:
      containers:
      - this should not be a string
`

var testManifestDeploymentNoImage = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  template:
    spec:
      containers:
      - name: the-container
`

var testManifestDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  template:
    spec:
      containers:
      - image: kyml/hello
        name: the-container
      initContainers:
      - image: kyml/init
        name: the-init-container
`

var testManifestDeploymentResolved = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  template:
    spec:
      containers:
      - image: kyml/hello@sha256:2cbb95c7479634c53bc2be243554a98d6928c189360fa958d2c970974e7f131f
        name: the-container
      initContainers:
      - image: kyml/init@sha256:2cbb95c7479634c53bc2be243554a98d6928c189360fa958d2c970974e7f131f
        name: the-init-container
`

func Test_resolveOptions_Validate(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error if any args",
			args: args{
				args: []string{"foo"},
			},
			wantErr: true,
		},
		{
			name: "success if no args",
			args: args{
				args: []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &resolveOptions{}
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("resolveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resolveOptions_Run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name          string
		args          args
		resolveOut    map[string]string
		resolveErr    error
		wantOut       string
		wantImageRefs []string
		wantErr       bool
	}{
		{
			name:          "skip unsupported kinds",
			args:          args{strings.NewReader(testManifestService)},
			resolveOut:    nil,
			resolveErr:    nil,
			wantOut:       testManifestService,
			wantImageRefs: nil,
			wantErr:       false,
		},
		{
			name:          "containers array does not exist",
			args:          args{strings.NewReader(testManifestDeploymentNoContainers)},
			resolveOut:    nil,
			resolveErr:    nil,
			wantOut:       testManifestDeploymentNoContainers,
			wantImageRefs: nil,
			wantErr:       false,
		},
		{
			name:          "container is no yaml object",
			args:          args{strings.NewReader(testManifestDeploymentMalformedContainer)},
			resolveOut:    nil,
			resolveErr:    nil,
			wantOut:       testManifestDeploymentMalformedContainer,
			wantImageRefs: nil,
			wantErr:       false,
		},
		{
			name:          "container has no image",
			args:          args{strings.NewReader(testManifestDeploymentNoImage)},
			resolveOut:    nil,
			resolveErr:    nil,
			wantOut:       testManifestDeploymentNoImage,
			wantImageRefs: nil,
			wantErr:       false,
		},
		{
			name:          "resolve errors",
			args:          args{strings.NewReader(testManifestDeployment)},
			resolveOut:    nil,
			resolveErr:    errors.New("oh no"),
			wantOut:       "",
			wantImageRefs: []string{"kyml/init"},
			wantErr:       true,
		},
		{
			name:          "resolve doesn't find image",
			args:          args{strings.NewReader(testManifestDeployment)},
			resolveOut:    nil,
			resolveErr:    nil,
			wantOut:       "",
			wantImageRefs: []string{"kyml/init"},
			wantErr:       true,
		},
		{
			name: "image gets resolved",
			args: args{strings.NewReader(testManifestDeployment)},
			resolveOut: map[string]string{
				"kyml/init":  "kyml/init@sha256:2cbb95c7479634c53bc2be243554a98d6928c189360fa958d2c970974e7f131f",
				"kyml/hello": "kyml/hello@sha256:2cbb95c7479634c53bc2be243554a98d6928c189360fa958d2c970974e7f131f",
			},
			resolveErr:    nil,
			wantOut:       testManifestDeploymentResolved,
			wantImageRefs: []string{"kyml/init", "kyml/hello"},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &resolveOptions{}
			out := &bytes.Buffer{}
			var gotImageRefs []string
			resolveImageMock := func(imageRef string) (string, error) {
				gotImageRefs = append(gotImageRefs, imageRef)
				if tt.resolveErr != nil {
					return "", tt.resolveErr
				} else {
					return tt.resolveOut[imageRef], nil
				}
			}

			if err := o.Run(tt.args.in, out, resolveImageMock); (err != nil) != tt.wantErr {
				t.Errorf("resolveOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotImageRefs, tt.wantImageRefs) {
				t.Errorf("resolveOptions.Run() imageRefs = %v, wantImageRefs %v", gotImageRefs, tt.wantImageRefs)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("resolveOptions.Run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
