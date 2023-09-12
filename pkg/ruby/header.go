// Copyright 2023 The Parca Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ruby

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/parca-dev/runtime-data/pkg/byteorder"
)

type RubyVersionOffsets struct {
	MajorVersion        int32 `yaml:"major_version"`
	MinorVersion        int32 `yaml:"minor_version"`
	PatchVersion        int32 `yaml:"patch_version"`
	VMOffset            int32 `yaml:"vm_offset"`
	VMSizeOffset        int32 `yaml:"vm_size_offset"`
	ControlFrameSizeof  int32 `yaml:"control_frame_t_sizeof"`
	CfpOffset           int32 `yaml:"cfp_offset"`
	LabelOffset         int32 `yaml:"label_offset"`
	PathFlavour         int32 `yaml:"path_flavour"`
	LineInfoSizeOffset  int32 `yaml:"line_info_size_offset"`
	LineInfoTableOffset int32 `yaml:"line_info_table_offset"`
	LinenoOffset        int32 `yaml:"lineno_offset"`
	MainThreadOffset    int32 `yaml:"main_thread_offset"`
	EcOffset            int32 `yaml:"ec_offset"`
}

func (rvo RubyVersionOffsets) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(unsafe.Sizeof(&rvo)))

	if err := binary.Write(buf, byteorder.GetHostByteOrder(), &rvo); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
