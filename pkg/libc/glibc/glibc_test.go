package glibc

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/parca-dev/runtime-data/pkg/libc"
)

func Test_getLayoutForArch(t *testing.T) {
	tests := []struct {
		name    string
		v       *semver.Version
		arch    string
		want    *libc.Layout
		wantErr bool
	}{
		{
			name: "2.29.0",
			v:    semver.MustParse("2.29.0"),
			arch: "amd64",
			want: &libc.Layout{
				PThreadSpecific1stblock: 0x310,
				PThreadSize:             2304,
				PThreadKeyData:          0x08,
				PThreadKeyDataSize:      0x10,
			},
		},
		{
			name: "2.39.0",
			v:    semver.MustParse("2.39.0"),
			arch: "amd64",
			want: &libc.Layout{
				PThreadSpecific1stblock: 0x310,
				PThreadSize:             2368,
				PThreadKeyData:          0x08,
				PThreadKeyDataSize:      0x10,
			},
		},
		{
			name: "2.29.0",
			v:    semver.MustParse("2.29.0"),
			arch: "arm64",
			want: &libc.Layout{
				PThreadSpecific1stblock: 0x110,
				PThreadSize:             1792,
				PThreadKeyData:          0x08,
				PThreadKeyDataSize:      0x10,
			},
		},
		{
			name: "2.37.0",
			v:    semver.MustParse("2.37.0"),
			arch: "amd64",
			want: &libc.Layout{
				PThreadSpecific1stblock: 0x310,
				PThreadSize:             2368,
				PThreadKeyData:          0x08,
				PThreadKeyDataSize:      0x10,
			},
		},
		{
			name: "2.37.0",
			v:    semver.MustParse("2.37.0"),
			arch: "arm64",
			want: &libc.Layout{
				PThreadSpecific1stblock: 0x110,
				PThreadSize:             1856,
				PThreadKeyData:          0x08,
				PThreadKeyDataSize:      0x10,
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
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(libc.Layout{})); diff != "" {
				t.Errorf("getLayoutForArch(%s) on %s mismatch (-want +got):\n%s", tt.name, tt.arch, diff)
			}
		})
	}
}
