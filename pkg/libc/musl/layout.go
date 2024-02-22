package musl

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/parca-dev/runtime-data/pkg/byteorder"
)

type Layout struct {
	PthreadSize int64 `yaml:"pthread_size"`
	PthreadTSD  int64 `yaml:"pthread_tsd"`
}

func (m Layout) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(unsafe.Sizeof(&m)))

	if err := binary.Write(buf, byteorder.GetHostByteOrder(), &m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
