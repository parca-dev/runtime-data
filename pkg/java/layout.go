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

package java

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/parca-dev/runtime-data/pkg/byteorder"
)

type Layout struct {
	CollectedHeapReserve uint64 `yaml:"collected_heap_reserve"`

	MemRegionStart uint64 `yaml:"mem_region_start"`
	MemRegionEnd   uint64 `yaml:"mem_region_end"`
	HeapWordSize   uint64 `yaml:"heap_word_size"`

	VMStructEntryTypeName  uint64 `yaml:"vm_struct_entry_type_name"`
	VMStructEntryFieldName uint64 `yaml:"vm_struct_entry_field_name"`
	VMStructEntryAddress   uint64 `yaml:"vm_struct_entry_address"`
	VMStructEntrySize      uint64 `yaml:"vm_struct_entry_size"`

	KlassName uint64 `yaml:"klass_name"`

	ConstantPoolHolder uint64 `yaml:"constant_pool_holder"`
	ConstantPoolSize   uint64 `yaml:"constant_pool_size"`

	OOPDescMetadata uint64 `yaml:"oop_desc_metadata"`
	OPPDescSize     uint64 `yaml:"oop_desc_size"`

	AccessFlags             uint64 `yaml:"access_flags"`
	SymbolLengthAndRefcount uint64 `yaml:"symbol_length_and_refcount"`
	SymbolBody              uint64 `yaml:"symbol_body"`

	MethodConst       uint64 `yaml:"method_const"`
	MethodAccessFlags uint64 `yaml:"method_access_flags"`
	MethodSize        uint64 `yaml:"method_size"`

	ConstMethodConstants      uint64 `yaml:"const_method_constants"`
	ConstMethodFlags          uint64 `yaml:"const_method_flags"`
	ConstMethodCodeSize       uint64 `yaml:"const_method_code_size"`
	ConstMethodNameIndex      uint64 `yaml:"const_method_name_index"`
	ConstMethodSignatureIndex uint64 `yaml:"const_method_signature_index"`
	ConstMethodSize           uint64 `yaml:"const_method_size"`

	CodeHeapMemory          uint64 `yaml:"code_heap_memory"`
	CodeHeapSegmap          uint64 `yaml:"code_heap_segmap"`
	CodeHeapLog2SegmentSize uint64 `yaml:"code_heap_log2_segment_size"`

	VirtualSpaceLowBoundary  uint64 `yaml:"virtual_space_low_boundary"`
	VirtualSpaceHighBoundary uint64 `yaml:"virtual_space_high_boundary"`
	VirtualSpaceLow          uint64 `yaml:"virtual_space_low"`
	VirtualSpaceHigh         uint64 `yaml:"virtual_space_high"`

	CodeBlobName         uint64 `yaml:"code_blob_name"`
	CodeBlobHeaderSize   uint64 `yaml:"code_blob_header_size"`
	CodeBlobContentBegin uint64 `yaml:"code_blob_content_begin"`
	// CodeBlobCodeStart    uint64 `yaml:"code_blob_code_start"`
	CodeBlobCodeBegin  uint64 `yaml:"code_blob_code_begin"`
	CodeBlobCodeEnd    uint64 `yaml:"code_blob_code_end"`
	CodeBlobDataOffset uint64 `yaml:"code_blob_data_offset"`
	CodeBlobFrameSize  uint64 `yaml:"code_blob_frame_size"`
	// CodeBlobFrameCompleteOffset uint64 `yaml:"code_blob_frame_complete_offset"`
	CodeBlobSize uint64 `yaml:"code_blob_size"`

	NMethodMethod             uint64 `yaml:"nmethod_method"`
	NMethodDependenciesOffset uint64 `yaml:"nmethod_dependencies_offset"`
	NMethodMetadataOffset     uint64 `yaml:"nmethod_metadata_offset"`
	NMethodScopesDataBegin    uint64 `yaml:"nmethod_scopes_data_begin"`
	NMethodScopesPCsOffset    uint64 `yaml:"nmethod_scopes_pcs_offset"`
	NMethodHandlerTableOffset uint64 `yaml:"nmethod_handler_table_offset"`
	NMethodDeoptHandlerBegin  uint64 `yaml:"nmethod_deopt_handler_begin"`
	NMethodOrigPCOffset       uint64 `yaml:"nmethod_orig_pc_offset"`
	// NMethodCompileID          uint64 `yaml:"nmethod_compile_id"`
	NMethodSize uint64 `yaml:"nmethod_size"`

	PCDescPCOffset          uint64 `yaml:"pc_desc_pc_offset"`
	PCDescScopeDecodeOffset uint64 `yaml:"pc_desc_scope_decode_offset"`
	PCDescSize              uint64 `yaml:"pc_desc_size"`

	NarrowPtrStructBase  uint64 `yaml:"narrow_ptr_struct_base"`
	NarrowPtrStructShift uint64 `yaml:"narrow_ptr_struct_shift"`

	BufferBlobSize    uint64 `yaml:"buffer_blob_size"`
	SingletonBlobSize uint64 `yaml:"singleton_blob_size"`
	RuntimeStubSize   uint64 `yaml:"runtime_stub_size"`
	SafepointBlobSize uint64 `yaml:"safepoint_blob_size"`

	CodeCacheStart uint64 `yaml:"code_cache_start"`
	CodeCacheEnd   uint64 `yaml:"code_cache_end"`

	DeoptHandler uint64 `yaml:"deopt_handler"`

	HeapBlockSize uint64 `yaml:"heap_block_size"`
	SegmentShift  uint64 `yaml:"segment_shift"`
}

func (jo Layout) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(unsafe.Sizeof(&jo)))

	if err := binary.Write(buf, byteorder.GetHostByteOrder(), &jo); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
