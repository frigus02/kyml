package k8syaml

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var validYamlIncludingNils = `
---
apiVersion: v1
kind: Namespace
metadata:
    name: the-namespace
---
---
apiVersion: v1
kind: Service

metadata:
  name: the-service
  namespace: the-namespace

spec:
  selector:
    deployment: the-deployment
`

var validYamlNicelyFormatted = `---
apiVersion: v1
kind: Namespace
metadata:
  name: the-namespace
---
apiVersion: v1
kind: Service
metadata:
  name: the-service
  namespace: the-namespace
spec:
  selector:
    deployment: the-deployment
`

var invalidYaml = `
hello world!!
`

var unstructuredDocuments = []unstructured.Unstructured{
	unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				"name": "the-namespace",
			},
		},
	},
	unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name":      "the-service",
				"namespace": "the-namespace",
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"deployment": "the-deployment",
				},
			},
		},
	},
}

func TestDecode(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []unstructured.Unstructured
		wantErr bool
	}{
		{
			name:    "valid yaml, including nil documents",
			args:    args{strings.NewReader(validYamlIncludingNils)},
			want:    unstructuredDocuments,
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			args:    args{strings.NewReader(invalidYaml)},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		documents []unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name:    "adds --- between yaml documents",
			args:    args{unstructuredDocuments},
			wantOut: validYamlNicelyFormatted,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := Encode(out, tt.args.documents); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Encode() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
