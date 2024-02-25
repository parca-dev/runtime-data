package glibc

import (
	"github.com/parca-dev/runtime-data/pkg/libc"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
)

type glibc struct {
	PThreadSpecific1stblock int64 `offsetof:"pthread.specific_1stblock" yaml:"pthread_specific_1stblock"`
	PThreadSize             int64 `sizeof:"pthread" yaml:"pthread_size"`
	PThreadKeyData          int64 `offsetof:"pthread_key_data.data" yaml:"pthread_key_data"`
	PThreadKeyDataSize      int64 `sizeof:"pthread_key_data" yaml:"pthread_key_data_size"`
}

func (g *glibc) Layout() runtimedata.RuntimeData {
	return &libc.Layout{
		PThreadSpecific1stblock: g.PThreadSpecific1stblock,
		PThreadSize:             g.PThreadSize,
		PThreadKeyData:          g.PThreadKeyData,
		PThreadKeyDataSize:      g.PThreadKeyDataSize,
	}
}

func DataMapForLayout(version string) runtimedata.LayoutMap {
	return &glibc{}
}
