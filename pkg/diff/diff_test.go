package diff

import "testing"

func TestDiff(t *testing.T) {
	type args struct {
		nameA string
		a     string
		nameB string
		b     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				nameA: "staging",
				a:     "1\n2\n3",
				nameB: "production",
				b:     "1\na\n3",
			},
			want:    "--- staging\n+++ production\n@@ -2 +2 @@\n-2\n+a\n",
			wantErr: false,
		},
		{
			name: "no difference",
			args: args{
				nameA: "staging",
				a:     "1\n2\n3",
				nameB: "production",
				b:     "1\n2\n3",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Diff(tt.args.nameA, tt.args.a, tt.args.nameB, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
