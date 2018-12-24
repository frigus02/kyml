package k8syaml

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGVKEquals(t *testing.T) {
	type args struct {
		a schema.GroupVersionKind
		b schema.GroupVersionKind
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal",
			args: args{
				schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
				schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			},
			want: true,
		},
		{
			name: "different group",
			args: args{
				schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
				schema.GroupVersionKind{Group: "extensions", Version: "v1", Kind: "Deployment"},
			},
			want: false,
		},
		{
			name: "different version",
			args: args{
				schema.GroupVersionKind{Group: "apps", Version: "v1beta1", Kind: "Deployment"},
				schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			},
			want: false,
		},
		{
			name: "different kind",
			args: args{
				schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
				schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GVKEquals(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("GVKEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
