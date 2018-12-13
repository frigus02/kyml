package cat

import (
	"reflect"
	"testing"
)

func Test_catOptions_Validate(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			o := &catOptions{}
			if err := o.Validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("catOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(o.files, tt.wantFiles) {
				t.Errorf("catOptions.files = %v, want %v", o.files, tt.wantFiles)
			}
		})
	}
}
