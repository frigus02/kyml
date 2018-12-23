package resolve

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_isSupportedKind(t *testing.T) {
	type args struct {
		gvk schema.GroupVersionKind
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "supported",
			args: args{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}},
			want: true,
		},
		{
			name: "not supported",
			args: args{schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Service"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedKind(tt.args.gvk); got != tt.want {
				t.Errorf("isSupportedKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
