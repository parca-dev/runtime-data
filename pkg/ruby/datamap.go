package ruby

import "github.com/parca-dev/runtime-data/pkg/runtimedata"

func DataMapForVersion(version string) runtimedata.LayoutMap {
	return &Ruby{}
}

type Ruby struct{}

func (Ruby) Layout() runtimedata.RuntimeData {
	return &Layout{}
}
