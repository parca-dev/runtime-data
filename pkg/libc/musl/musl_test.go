package musl

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-cmp/cmp"
)

func Test_getLayoutForArch(t *testing.T) {

	tests := []struct {
		name    string
		v       *semver.Version
		arch    string
		want    *Layout
		wantErr bool
	}{
		{
			name: "1.2.2",
			v:    semver.MustParse("1.2.2"),
			arch: "amd64",
			want: &Layout{
				PthreadSize: 200,
				PthreadTSD:  128,
			},
		},
		{
			name: "1.2.2",
			v:    semver.MustParse("1.2.2"),
			arch: "arm64",
			want: &Layout{
				PthreadSize: 200,
				PthreadTSD:  112,
			},
		},
		{
			name: "1.1.19",
			v:    semver.MustParse("1.1.19"),
			arch: "amd64",
			want: &Layout{
				PthreadSize: 280,
				PthreadTSD:  152,
			},
		},
		{
			name: "1.1.19",
			v:    semver.MustParse("1.1.19"),
			arch: "arm64",
			want: &Layout{
				PthreadSize: 280,
				PthreadTSD:  152,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := getLayoutForArch(tt.v, tt.arch)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLayoutForArch(%s) on %s error = %v, wantErr %v", tt.name, tt.arch, err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(Layout{})); diff != "" {
				t.Errorf("getLayoutForArch(%s) on %s mismatch (-want +got):\n%s", tt.name, tt.arch, diff)
			}
		})
	}
}
