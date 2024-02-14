//go:build integration
// +build integration

package integration

import (
	"bytes"
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/parca-dev/runtime-data/pkg/datamap"
	"github.com/parca-dev/runtime-data/pkg/ruby"
	"gopkg.in/yaml.v3"
)

const TargetDirRuby = "binaries/ruby"

var rubyVersions = []string{
	"2.6.0",
	"2.6.3",
	"2.7.1",
	"2.7.4",
	"2.7.6",
	"3.0.0",
	"3.0.4",
	"3.1.2",
	"3.1.3",
	"3.2.0",
	"3.2.1",
}

func TestRubyIntegration(t *testing.T) {
	t.Parallel()

	for _, version := range rubyVersions {
		version := version
		t.Run(version, func(t *testing.T) {
			t.Parallel()

			layoutMap := ruby.DataMapForLayout(version)
			if layoutMap == nil {
				t.Fatalf("ruby.DataMapForLayout(%s) = nil", version)
			}

			dm, err := datamap.New(layoutMap)
			if err != nil {
				t.Fatalf("ruby.GenerateDataMap(%s) = %v", version, err)
			}

			input := fmt.Sprintf("%s/libruby.so.%s", TargetDirRuby, version)

			f, err := elf.Open(input)
			if err != nil {
				t.Fatalf("elf.Open() = %v", err)
			}

			dwarfData, err := f.DWARF()
			if err != nil {
				t.Fatalf("f.DWARF() = %v", err)
			}

			if err := dm.ReadFromDWARF(dwarfData); err != nil {
				t.Errorf("input: %s", input)
				t.Fatalf("datamap.ReadFromDWARF() = %v", err)
			}

			got := layoutMap.Layout().(*ruby.Layout)

			golden := filepath.Join("testdata", fmt.Sprintf("ruby_%s.yaml", sanitizeIdentifier(version)))
			if *update {
				var buf bytes.Buffer
				enc := yaml.NewEncoder(&buf)
				enc.SetIndent(2)
				if err := enc.Encode(got); err != nil {
					t.Fatalf("yaml.Encode() = %v", err)
				}
				if err := os.WriteFile(golden, buf.Bytes(), 0o644); err != nil {
					t.Fatalf("os.WriteFile() = %v", err)
				}
			}

			wantData, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("os.ReadFile() = %v", err)
			}

			var want ruby.Layout
			yaml.Unmarshal(wantData, &want)

			if diff := cmp.Diff(want, *got, cmp.AllowUnexported(ruby.Layout{})); diff != "" {
				t.Errorf("input: %s, golden: %s", input, golden)
				t.Errorf("ruby(%s) mismatch (-want +got):\n%s", version, diff)
			}
		})
	}
}
