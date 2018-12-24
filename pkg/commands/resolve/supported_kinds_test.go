package resolve

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_getPathToPodSpec(t *testing.T) {
	type args struct {
		gvk schema.GroupVersionKind
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "supported",
			args: args{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}},
			want: []string{"spec", "template", "spec"},
		},
		{
			name: "not supported",
			args: args{schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Service"}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPathToPodSpec(tt.args.gvk); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPathToPodSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
