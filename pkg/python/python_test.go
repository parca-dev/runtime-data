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
	"fmt"
	"runtime"
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

var allSupportedArchs = []string{"amd64", "arm64"}

func TestGetLayout(t *testing.T) {
	tests := []struct {
		version string
		archs   []string
		want    *Layout
		wantErr bool
	}{
		{
			version: "2.7.15",
			archs:   allSupportedArchs,
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
			archs:   allSupportedArchs,
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
			archs:   allSupportedArchs,
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
		var (
			version = tt.version
			want    = tt.want
		)
		for _, arch := range tt.archs {
			arch := arch
			t.Run(version, func(t *testing.T) {
				v, err := semver.StrictNewVersion(version)
				if err != nil {
					t.Errorf("StrictNewVersion() error = %v", err)
					return
				}
				_, got, err := getLayoutForArch(v, arch)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetLayout(%s on %s) error = %v, wantErr %v", version, arch, err, tt.wantErr)
					return
				}
				if diff := cmp.Diff(want, got, cmp.AllowUnexported(Layout{})); diff != "" {
					t.Errorf("GetLayout(%s on %s) mismatch (-want +got):\n%s", version, arch, diff)
				}
			})
		}
	}
}

func TestGetInitialState(t *testing.T) {
	tests := []struct {
		version string
		arch    string
		want    *InitialState
		wantErr bool
	}{
		{
			version: "2.7.15",
			want:    nil,
			wantErr: true,
		},
		{
			version: "3.3.7",
			want:    nil,
			wantErr: true,
		},
		{
			version: "3.6.6",
			want:    nil,
			wantErr: true,
		},
		{
			version: "3.7.2",
			want: &InitialState{
				InterpreterHead:    24,
				ThreadStateCurrent: 1392,
				AutoTSSKey:         1416,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.7.4",
			want: &InitialState{
				InterpreterHead:    24,
				ThreadStateCurrent: 1480,
				AutoTSSKey:         1504,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.8.0",
			want: &InitialState{
				InterpreterHead:    32,
				ThreadStateCurrent: 1368,
				AutoTSSKey:         1392,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.9.6",
			want: &InitialState{
				InterpreterHead:    32,
				ThreadStateCurrent: 568,
				AutoTSSKey:         584,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.10.0",
			want: &InitialState{
				InterpreterHead:    32,
				ThreadStateCurrent: 568,
				AutoTSSKey:         584,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.11.0",
			want: &InitialState{
				InterpreterHead:    40,
				ThreadStateCurrent: 576,
				AutoTSSKey:         592,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.12.0",
			want: &InitialState{
				InterpreterHead:    40,
				ThreadStateCurrent: -1,
				AutoTSSKey:         1544,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		// arm64
		{
			version: "2.7.15",
			arch:    "arm64",
			want:    nil,
			wantErr: true,
		},
		{
			version: "3.3.7",
			arch:    "arm64",
			want:    nil,
			wantErr: true,
		},
		{
			version: "3.6.6",
			arch:    "arm64",
			want:    nil,
			wantErr: true,
		},
		{
			version: "3.7.2",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    24,
				ThreadStateCurrent: 1408,
				AutoTSSKey:         1432,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.7.4",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    24,
				ThreadStateCurrent: 1496,
				AutoTSSKey:         1520,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.8.0",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    32,
				ThreadStateCurrent: 1384,
				AutoTSSKey:         1408,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.9.6",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    32,
				ThreadStateCurrent: 584,
				AutoTSSKey:         600,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.10.0",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    32,
				ThreadStateCurrent: 584,
				AutoTSSKey:         600,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.11.0",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    40,
				ThreadStateCurrent: 592,
				AutoTSSKey:         608,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
		{
			version: "3.12.0",
			arch:    "arm64",
			want: &InitialState{
				InterpreterHead:    40,
				ThreadStateCurrent: -1,
				AutoTSSKey:         1544,
				PyTSS: PyTSSKey{
					Key:  4,
					Size: 8,
				},
			},
		},
	}
	for _, tt := range tests {
		var (
			arch    = tt.arch
			version = tt.version
			want    = tt.want
		)
		if tt.arch == "" {
			arch = runtime.GOARCH
		}
		name := fmt.Sprintf("%s on %s", tt.version, arch)
		t.Run(name, func(t *testing.T) {
			v, err := semver.StrictNewVersion(version)
			if err != nil {
				t.Errorf("StrictNewVersion() error = %v", err)
				return
			}
			_, got, err := GetInitialStateForArch(v, arch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInitialState(%s) error = %v, wantErr %v", name, err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(want, got, cmp.AllowUnexported(InitialState{})); diff != "" {
				t.Errorf("GetInitialState(%s) mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}
