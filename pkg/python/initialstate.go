package python

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/Masterminds/semver/v3"

	"github.com/parca-dev/runtime-data/pkg/byteorder"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"github.com/parca-dev/runtime-data/pkg/version"
)

var (
	above3_7  = version.MustParseConstraints("3.7.x - 3.11.x")
	above3_12 = version.MustParseConstraints(">=3.12.0-0")
)

func DataMapForInitialState(v string) runtimedata.InitialStateMap {
	lookupVersion := semver.MustParse(v)
	switch {
	case above3_7.Check(lookupVersion):
		return &initialState38{}
	case above3_12.Check(lookupVersion):
		return &initialState312{}
	}
	// There is no offsets to be extracted version below 3.7
	// so we return nil
	return nil
}

type InitialStateWithVersion struct {
	Version      runtimedata.Version `yaml:"version"`
	InitialState InitialState        `yaml:"initial_state"`
}

func WithVersion(version string, is InitialState) (InitialStateWithVersion, error) {
	sm, err := semver.NewVersion(version)
	if err != nil {
		return InitialStateWithVersion{}, fmt.Errorf("failed to parse version (%s): %w", version, err)
	}
	v := runtimedata.Version{
		Major: sm.Major(),
		Minor: sm.Minor(),
		Patch: sm.Patch(),
	}
	return InitialStateWithVersion{
		Version:      v,
		InitialState: is,
	}, nil
}

// The state kept in the static global variables,
// we read them from the memory directly.
// Check the agent for more details.
type InitialState struct {
	InterpreterHead    int64    `yaml:"interpreter_head"`
	ThreadStateCurrent int64    `yaml:"tstate_current"`
	AutoTSSKey         int64    `yaml:"auto_tss_key"`
	PyTSSKey           PyTSSKey `yaml:"tss"`
}

type PyTSSKey struct {
	Key  int64 `yaml:"key"`
	Size int64 `yaml:"size"`
}

func (i InitialState) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(unsafe.Sizeof(&i)))

	if err := binary.Write(buf, byteorder.GetHostByteOrder(), &i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type initialState312 struct {
	InterpreterHead int64 `offsetof:"_PyRuntimeState.interpreters.head"`
	AutoTSSKey      int64 `offsetof:"_PyRuntimeState.autoTSSkey"`
	PyTSSKey        int64 `offsetof:"_Py_tss_t._key"`
	PyTSSSize       int64 `sizeof:"_Py_tss_t"`
}

func (i initialState312) InitialState() runtimedata.RuntimeData {
	return &InitialState{
		InterpreterHead: i.InterpreterHead,
		AutoTSSKey:      i.AutoTSSKey,
		// https://github.com/python/cpython/issues/103323
		ThreadStateCurrent: doesNotExist,
		PyTSSKey: PyTSSKey{
			Key:  i.PyTSSKey,
			Size: i.PyTSSSize,
		},
	}
}

// https://pythondev.readthedocs.io/pystate.html
// Python 3.7: PyThreadState_GET() reads _PyThreadState_Current (atomic variable).
// Python 3.8: _PyThreadState_Current becomes _PyRuntime.gilstate.tstate_current.
type initialState38 struct {
	InterpreterHead    int64 `offsetof:"_PyRuntimeState.interpreters.head"`
	ThreadStateCurrent int64 `offsetof:"_PyRuntimeState.gilstate.tstate_current"`
	AutoTSSKey         int64 `offsetof:"_PyRuntimeState.gilstate.autoTSSkey"`
	PyTSSKey           int64 `offsetof:"_Py_tss_t._key"`
	PyTSSSize          int64 `sizeof:"_Py_tss_t"`
}

func (i initialState38) InitialState() runtimedata.RuntimeData {
	return &InitialState{
		InterpreterHead:    i.InterpreterHead,
		ThreadStateCurrent: i.ThreadStateCurrent,
		AutoTSSKey:         i.AutoTSSKey,
		PyTSSKey: PyTSSKey{
			Key:  i.PyTSSKey,
			Size: i.PyTSSSize,
		},
	}
}
