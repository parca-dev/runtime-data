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
	"reflect"
	"testing"
)

func TestGetVersions(t *testing.T) {
	versions, err := GetVersions()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(versions)
}

func TestGetVersionMap(t *testing.T) {
	tests := []struct {
		version string
		want    Layout
		wantErr bool
	}{
		{
			version: "3.0.4",
			want: Layout{
				MajorVersion:        3,
				MinorVersion:        0,
				PatchVersion:        4,
				VMOffset:            0,
				VMSizeOffset:        8,
				ControlFrameSizeof:  56,
				CfpOffset:           16,
				LabelOffset:         16,
				PathFlavour:         1,
				LineInfoSizeOffset:  136,
				LineInfoTableOffset: 120,
				LinenoOffset:        0,
				MainThreadOffset:    32,
				EcOffset:            520,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			m, err := GetVersionMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersionMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := m[tt.version]
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVersionMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
