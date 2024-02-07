package libc

import "github.com/parca-dev/runtime-data/pkg/runtimedata"

func DataMapForVersion(version string) runtimedata.LayoutMap {
	return &LibC{}
}

type LibC struct{}

func (LibC) Layout() runtimedata.RuntimeData {
	return &Layout{}
}
