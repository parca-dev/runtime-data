package datamap

import (
	"debug/elf"
	"fmt"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func arch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		return runtime.GOARCH
	}
}

type testMap struct {
	Size       int `sizeof:"test_t"`
	A          int `offsetof:"test_t.a"`
	B          int `offsetof:"test_t.b"`
	NestedSize int `sizeof:"test_t.nested"`

	NestedA          int `offsetof:"test_t.nested.nested_a"`
	NestedB          int `offsetof:"test_t.nested.nested_b"`
	DeeplyNestedSize int `sizeof:"test_t.nested.deeply_nested"`

	DeeplyNestedA int `offsetof:"test_t.nested.deeply_nested.deeply_nested_a"`
	DeeplyNestedB int `offsetof:"test_t.nested.deeply_nested.deeply_nested_b"`
}

func TestDataMap_ReadFromDWARF(t *testing.T) {
	type args struct {
		inputPath string
		lm        any
	}
	tests := []struct {
		name    string
		args    args
		want    *testMap
		wantErr bool
	}{
		{
			name: "Test ReadFromDWARF",
			args: args{
				inputPath: fmt.Sprintf("testdata/%s/test", arch()),
				lm:        &testMap{},
			},
			want: &testMap{
				Size:       24,
				A:          0,
				B:          4,
				NestedSize: 16,

				NestedA:          8,
				NestedB:          12,
				DeeplyNestedSize: 8,

				DeeplyNestedA: 16,
				DeeplyNestedB: 20,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layoutMap := tt.args.lm
			dm, err := New(layoutMap)
			if err != nil {
				t.Fatalf("failed to generate query: %v", err)
			}

			input := tt.args.inputPath
			ef, err := elf.Open(input)
			if err != nil {
				t.Fatalf("failed to open ELF file: %v", err)
			}
			defer ef.Close()

			if err := dm.ReadFromDWARF(ef); err != nil {
				t.Fatalf("failed to read DWARF data: %v", err)
			}

			diff := cmp.Diff(tt.want, layoutMap, cmp.AllowUnexported(testMap{}))
			if diff != "" {
				t.Errorf("ReadFromDWARF() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
