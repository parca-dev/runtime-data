package runtimedata

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type Version struct {
	Major uint64 `yaml:"major"`
	Minor uint64 `yaml:"minor"`
	Patch uint64 `yaml:"patch"`
}

type LayoutWithVersion struct {
	Version Version     `yaml:"version"`
	Layout  RuntimeData `yaml:"layout"`
}

type LayoutMap interface {
	Layout() RuntimeData
}

type RuntimeData interface {
	Data() ([]byte, error)
}

func WithVersion(version string, l RuntimeData) (LayoutWithVersion, error) {
	sm, err := semver.NewVersion(version)
	if err != nil {
		return LayoutWithVersion{}, fmt.Errorf("failed to parse version (%s): %w", version, err)
	}
	v := Version{
		Major: sm.Major(),
		Minor: sm.Minor(),
		Patch: sm.Patch(),
	}
	return LayoutWithVersion{
		Version: v,
		Layout:  l,
	}, nil
}
