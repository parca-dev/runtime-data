// Copyright 2023 The Parca Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ruby

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-cmp/cmp"
)

func TestGetLayouts(t *testing.T) {
	layouts, err := loadLayouts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(layouts)
}

func TestGetLayout(t *testing.T) {
	tests := []struct {
		version string
		want    *Layout
		wantErr bool
	}{
		{
			version: "2.7.3",
			want: &Layout{
				VMSizeOffset:        8,
				ControlFrameSizeof:  56,
				CfpOffset:           16,
				LabelOffset:         16,
				PathFlavour:         1,
				LineInfoSizeOffset:  136,
				LineInfoTableOffset: 120,
				MainThreadOffset:    192,
				EcOffset:            32,
			},
		},
		{
			version: "3.0.4",
			want: &Layout{
				VMOffset:            0,
				VMSizeOffset:        8,
				ControlFrameSizeof:  56,
				CfpOffset:           16,
				LabelOffset:         16,
				PathFlavour:         1,
				LineInfoSizeOffset:  136,
				LineInfoTableOffset: 120,
				LinenoOffset:        0,
				MainThreadOffset:    32, // 40
				EcOffset:            520,
			},
		},
		{
			version: "3.1.2",
			want: &Layout{
				VMSizeOffset:        8,
				ControlFrameSizeof:  64,
				CfpOffset:           16,
				LabelOffset:         16,
				PathFlavour:         1,
				LineInfoSizeOffset:  136,
				LineInfoTableOffset: 120,
				MainThreadOffset:    32, // 40
				EcOffset:            520,
			},
		},
		{
			version: "3.2.1",
			want: &Layout{
				VMSizeOffset:        8,
				ControlFrameSizeof:  64,
				CfpOffset:           16,
				LabelOffset:         16,
				PathFlavour:         1,
				LineInfoSizeOffset:  128,
				LineInfoTableOffset: 112,
				MainThreadOffset:    32, // 40
				EcOffset:            520,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			version, err := semver.StrictNewVersion(tt.version)
			if err != nil {
				t.Errorf("StrictNewVersion() error = %v", err)
				return
			}
			_, got, err := GetLayout(version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got.(*Layout), cmp.AllowUnexported(Layout{})); diff != "" {
				t.Errorf("GetLayout() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
