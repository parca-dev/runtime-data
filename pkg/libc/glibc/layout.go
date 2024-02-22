package glibc

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/parca-dev/runtime-data/pkg/byteorder"
)

type Layout struct {
	PthreadSpecific1stblock int64 `yaml:"pthread_specific_1stblock"`
	PthreadSize             int64 `yaml:"pthread_size"`
	PthreadKeyData          int64 `yaml:"pthread_key_data"`
	PthreadKeyDataSize      int64 `yaml:"pthread_key_data_size"`
}

func (l Layout) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(unsafe.Sizeof(&l)))

	if err := binary.Write(buf, byteorder.GetHostByteOrder(), &l); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
