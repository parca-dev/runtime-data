package musl

import (
	"github.com/parca-dev/runtime-data/pkg/libc"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
)

type musl struct {
	PThreadSize int64 `sizeof:"__pthread"`
	PThreadTSD  int64 `offsetof:"__pthread.tsd"`
}

func (m *musl) Layout() runtimedata.RuntimeData {
	return &libc.Layout{
		PThreadSize:             m.PThreadSize,
		PThreadSpecific1stblock: m.PThreadTSD,
		PThreadKeyData:          0, // unused.
		// pthread_key_t: TYPEDEF unsigned pthread_key_t;
		PThreadKeyDataSize: 8,
	}
}

func DataMapForLayout(version string) runtimedata.LayoutMap {
	return &musl{}
}
