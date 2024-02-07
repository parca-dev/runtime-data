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
package python

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
			version: "3.11.0",
			want: Layout{
				PyObject: PyObject{
					ObType: 8,
				},
				PyString: PyString{
					Data: 48,
					Size: -1,
				},
				PyTypeObject: PyTypeObject{
					TPName: 24,
				},
				PyThreadState: PyThreadState{
					Next:           8,
					Interp:         16,
					Frame:          -1,
					ThreadID:       152,
					NativeThreadID: 160,
					CFrame:         56,
				},
				PyCFrame: PyCFrame{
					CurrentFrame: 8,
				},
				PyInterpreterState: PyInterpreterState{
					TStateHead: 16,
				},
				PyRuntimeState: PyRuntimeState{
					InterpMain: 48,
				},
				PyFrameObject: PyFrameObject{
					FBack:       48,
					FCode:       32,
					FLineno:     -1,
					FLocalsplus: 72,
				},
				PyCodeObject: PyCodeObject{
					CoFilename:    112,
					CoName:        120,
					CoVarnames:    96,
					CoFirstlineno: 72,
				},
				PyTupleObject: PyTupleObject{
					ObItem: 24,
				},
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
