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
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-cmp/cmp"
)

func TestGetVersions(t *testing.T) {
	versions, err := loadVersionLayouts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(versions)
}

func TestGetLayout(t *testing.T) {
	tests := []struct {
		version string
		want    *Layout
		wantErr bool
	}{
		{
			version: "2.7.15",
			want: &Layout{
				PyCodeObject: PyCodeObject{
					CoFilename:    80,
					CoName:        88,
					CoVarnames:    56,
					CoFirstlineno: 96},
				PyFrameObject: PyFrameObject{
					FBack:       24,
					FCode:       32,
					FLineno:     124,
					FLocalsplus: 376},
				PyInterpreterState: PyInterpreterState{TStateHead: 8},
				PyObject:           PyObject{ObType: 8},
				PyRuntimeState:     PyRuntimeState{InterpMain: -1},
				PyString:           PyString{Data: 36, Size: 16},
				PyThreadState: PyThreadState{
					Interp:         8,
					Frame:          16,
					ThreadID:       144,
					NativeThreadID: -1,
					CFrame:         -1,
				},
				PyTupleObject: PyTupleObject{ObItem: 24},
				PyTypeObject:  PyTypeObject{TPName: 24},
			},
		},
		{
			version: "3.6.6",
			want: &Layout{
				PyCodeObject: PyCodeObject{
					CoFilename:    96,
					CoName:        104,
					CoVarnames:    64,
					CoFirstlineno: 36,
				},
				PyFrameObject: PyFrameObject{
					FBack:       24,
					FCode:       32,
					FLineno:     124,
					FLocalsplus: 376,
				},
				PyInterpreterState: PyInterpreterState{TStateHead: 8},
				PyObject:           PyObject{ObType: 8},
				PyRuntimeState:     PyRuntimeState{InterpMain: -1},
				PyString: PyString{
					Data: 48,
					Size: 16,
				},
				PyThreadState: PyThreadState{
					Next:           8,
					Interp:         16,
					Frame:          24,
					ThreadID:       152,
					NativeThreadID: -1,
					CFrame:         -1,
				},
				PyTupleObject: PyTupleObject{ObItem: 24},
				PyTypeObject:  PyTypeObject{TPName: 24},
			},
		},
		{
			version: "3.11.0",
			want: &Layout{
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
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(Layout{})); diff != "" {
				t.Errorf("GetLayout() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
