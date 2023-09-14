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
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/parca-dev/runtime-data/pkg/byteorder"
)

// PyCFrame
type PyCFrame struct {
	CurrentFrame int64 `yaml:"current_frame"`
}

// PyTupleObject
type PyTupleObject struct {
	ObItem int64 `yaml:"ob_item"`
}

// PyThreadState
type PyThreadState struct {
	Next           int64 `yaml:"next"`
	Interp         int64 `yaml:"interp"`
	Frame          int64 `yaml:"frame"`
	ThreadID       int64 `yaml:"thread_id"`
	NativeThreadID int64 `yaml:"native_thread_id"`
	CFrame         int64 `yaml:"cframe"`
}

// PyString
type PyString struct {
	Data int64 `yaml:"data"`
	Size int64 `yaml:"size"`
}

// PyRuntimeState
type PyRuntimeState struct {
	InterpMain int64 `yaml:"interp_main"`
}

// PyTypeObject
type PyTypeObject struct {
	TPName int64 `yaml:"tp_name"`
}

// PyInterpreterState
type PyInterpreterState struct {
	TStateHead int64 `yaml:"tstate_head"`
}

// PyCodeObject
type PyCodeObject struct {
	CoFilename    int64 `yaml:"co_filename"`
	CoName        int64 `yaml:"co_name"`
	CoVarnames    int64 `yaml:"co_varnames"`
	CoFirstlineno int64 `yaml:"co_firstlineno"`
}

// PyObject
type PyObject struct {
	ObType int64 `yaml:"ob_type"`
}

// PyFrameObject
type PyFrameObject struct {
	FBack       int64 `yaml:"f_back"`
	FCode       int64 `yaml:"f_code"`
	FLineno     int64 `yaml:"f_lineno"`
	FLocalsplus int64 `yaml:"f_localsplus"`
}

type VersionOffsets struct {
	MajorVersion       uint32             `yaml:"major_version"`
	MinorVersion       uint32             `yaml:"minor_version"`
	PatchVersion       uint32             `yaml:"patch_version"`
	_padding           [4]byte            // Padding for alignment.
	PyCFrame           PyCFrame           `yaml:"py_cframe"`
	PyCodeObject       PyCodeObject       `yaml:"py_code_object"`
	PyFrameObject      PyFrameObject      `yaml:"py_frame_object"`
	PyInterpreterState PyInterpreterState `yaml:"py_interpreter_state"`
	PyObject           PyObject           `yaml:"py_object"`
	PyRuntimeState     PyRuntimeState     `yaml:"py_runtime_state"`
	PyString           PyString           `yaml:"py_string"`
	PyThreadState      PyThreadState      `yaml:"py_thread_state"`
	PyTupleObject      PyTupleObject      `yaml:"py_tuple_object"`
	PyTypeObject       PyTypeObject       `yaml:"py_type_object"`
}

func (pvo VersionOffsets) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(unsafe.Sizeof(&pvo)))

	if err := binary.Write(buf, byteorder.GetHostByteOrder(), &pvo); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
