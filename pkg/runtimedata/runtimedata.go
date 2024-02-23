package runtimedata

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type LayoutMap interface {
	Layout() RuntimeData
}

type InitialStateMap interface {
	InitialState() RuntimeData
}

type RuntimeData interface {
	Data() ([]byte, error)
}

type Version struct {
	Major uint64 `yaml:"major"`
	Minor uint64 `yaml:"minor"`
	Patch uint64 `yaml:"patch"`
}

type DataWithVersion struct {
	Version Version        `yaml:"version"`
	Data    map[string]any `yaml:"data"`
}

func WithVersion(version string, data map[string]any) (DataWithVersion, error) {
	sm, err := semver.NewVersion(version)
	if err != nil {
		return DataWithVersion{}, fmt.Errorf("failed to parse version (%s): %w", version, err)
	}
	v := Version{
		Major: sm.Major(),
		Minor: sm.Minor(),
		Patch: sm.Patch(),
	}
	return DataWithVersion{
		Version: v,
		Data:    data,
	}, nil
}
