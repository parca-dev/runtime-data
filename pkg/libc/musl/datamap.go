package musl

import "github.com/parca-dev/runtime-data/pkg/runtimedata"

type musl struct {
	PThreadSize int64 `sizeof:"__pthread" yaml:"pthread_size"`
	PThreadTSD  int64 `offsetof:"__pthread.tsd" yaml:"pthread_tsd"`
}

func (m *musl) Layout() runtimedata.RuntimeData {
	return &Layout{
		PthreadSize: m.PThreadSize,
		PthreadTSD:  m.PThreadTSD,
	}
}

func DataMapForLayout(version string) runtimedata.LayoutMap {
	return &musl{}
}
