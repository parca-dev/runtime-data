package python

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"github.com/parca-dev/runtime-data/pkg/version"
)

type InitialState struct {
	InterpreterHead    int64 `offsetof:"_PyRuntimeState.interpreters.head" yaml:"interpreter_head"`
	ThreadStateCurrent int64 `offsetof:"_PyRuntimeState.gilstate.tstate_current" yaml:"tstate_current"`
}

var above3_7 = version.MustParseConstraints(">=3.7.0")

func DataMapForInitialState(v string) *InitialState {
	lookupVersion := semver.MustParse(v)
	if above3_7.Check(lookupVersion) {
		return &InitialState{}
	}
	// There is no offsets to be extracted version below 3.7
	// so we return nil
	// The state kept in the static global variables,
	// we read them from the memory directly.
	// Check the agent for more details.
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
