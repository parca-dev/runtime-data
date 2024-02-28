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

package openjdk

import (
	"github.com/Masterminds/semver/v3"

	"github.com/parca-dev/runtime-data/pkg/java"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"github.com/parca-dev/runtime-data/pkg/version"
)

func DataMapForLayout(v string) runtimedata.LayoutMap {
	javaVersions := map[*semver.Constraints]runtimedata.LayoutMap{
		version.MustParseConstraints(">=17.0.0-0"): &openjdk{},
	}
	lookupVersion := semver.MustParse(v)
	for c, m := range javaVersions {
		if c.Check(lookupVersion) {
			return m
		}
	}
	return nil
}

type openjdk struct {
	CollectedHeapReserve uint64 `offsetof:"CollectedHeap._reserved"`

	MemRegionStart uint64 `offsetof:"MemRegion._start"`
	MemRegionEnd   uint64 `offsetof:"MemRegion._end"`
	HeapWordSize   uint64 `sizeof:"HeapWord"`

	VMStructEntryTypeName  uint64 `offsetof:"VMStructEntry.typeName"`
	VMStructEntryFieldName uint64 `offsetof:"VMStructEntry.fieldName"`
	VMStructEntryAddress   uint64 `offsetof:"VMStructEntry.address"`
	VMStructEntrySize      uint64 `sizeof:"VMStructEntry"`

	KlassName uint64 `offsetof:"Klass._name"`

	ConstantPoolHolder uint64 `offsetof:"ConstantPool._pool_holder"`
	ConstantPoolSize   uint64 `offsetof:"ConstantPool._size"`

	OOPDescMetadata uint64 `offsetof:"oopDesc._metadata"`
	OPPDescSize     uint64 `sizeof:"oopDesc"`

	AccessFlags             uint64 `offsetof:"AccessFlags._flags"`
	SymbolLengthAndRefcount uint64 `offsetof:"Symbol._length_and_refcount"`
	SymbolBody              uint64 `offsetof:"Symbol._body"`

	MethodConst       uint64 `offsetof:"Method._constMethod"`
	MethodAccessFlags uint64 `offsetof:"Method._access_flags"`
	MethodSize        uint64 `sizeof:"Method"`

	ConstMethodConstants      uint64 `offsetof:"ConstMethod._constants"`
	ConstMethodFlags          uint64 `offsetof:"ConstMethod._flags"`
	ConstMethodCodeSize       uint64 `offsetof:"ConstMethod._code_size"`
	ConstMethodNameIndex      uint64 `offsetof:"ConstMethod._name_index"`
	ConstMethodSignatureIndex uint64 `offsetof:"ConstMethod._signature_index"`
	ConstMethodSize           uint64 `sizeof:"ConstMethod"`

	CodeHeapMemory          uint64 `offsetof:"CodeHeap._memory"`
	CodeHeapSegmap          uint64 `offsetof:"CodeHeap._segmap"`
	CodeHeapLog2SegmentSize uint64 `offsetof:"CodeHeap._log2_segment_size"`

	VirtualSpaceLowBoundary  uint64 `offsetof:"VirtualSpace._low_boundary"`
	VirtualSpaceHighBoundary uint64 `offsetof:"VirtualSpace._high_boundary"`
	VirtualSpaceLow          uint64 `offsetof:"VirtualSpace._low"`
	VirtualSpaceHigh         uint64 `offsetof:"VirtualSpace._high"`

	CodeBlobName         uint64 `offsetof:"CodeBlob._name"`
	CodeBlobHeaderSize   uint64 `offsetof:"CodeBlob._header_size"`
	CodeBlobContentBegin uint64 `offsetof:"CodeBlob._content_begin"`
	// CodeBlobCodeStart    uint64 `offsetof:"CodeBlob._code_start"`
	CodeBlobCodeBegin  uint64 `offsetof:"CodeBlob._code_begin"`
	CodeBlobCodeEnd    uint64 `offsetof:"CodeBlob._code_end"`
	CodeBlobDataOffset uint64 `offsetof:"CodeBlob._data_offset"`
	CodeBlobFrameSize  uint64 `offsetof:"CodeBlob._frame_size"`
	// CodeBlobFrameCompleteOffset uint64 `offsetof:"CodeBlob._frame_complete_offset"`
	CodeBlobSize uint64 `sizeof:"CodeBlob"`

	NMethodMethod             uint64 `offsetof:"nmethod._method"`
	NMethodDependenciesOffset uint64 `offsetof:"nmethod._dependencies_offset"`
	NMethodMetadataOffset     uint64 `offsetof:"nmethod._metadata_offset"`
	NMethodScopesDataBegin    uint64 `offsetof:"nmethod._scopes_data_begin"`
	NMethodScopesPCsOffset    uint64 `offsetof:"nmethod._scopes_pcs_offset"`
	NMethodHandlerTableOffset uint64 `offsetof:"nmethod._handler_table_offset"`
	NMethodDeoptHandlerBegin  uint64 `offsetof:"nmethod._deopt_handler_begin"`
	NMethodOrigPCOffset       uint64 `offsetof:"nmethod._orig_pc_offset"`
	// NMethodCompileID          uint64 `offsetof:"nmethod._compile_id"`
	NMethodSize uint64 `sizeof:"nmethod"`

	PCDescPCOffset          uint64 `offsetof:"PcDesc._pc_offset"`
	PCDescScopeDecodeOffset uint64 `offsetof:"PcDesc._scope_decode_offset"`
	PCDescSize              uint64 `sizeof:"PcDesc"`

	NarrowPtrStructBase  uint64 `offsetof:"NarrowPtrStruct._base"`
	NarrowPtrStructShift uint64 `offsetof:"NarrowPtrStruct._shift"`

	// BufferBlobSize    uint64 `sizeof:"BufferBlob"`
	// SingletonBlobSize uint64 `sizeof:"SingletonBlob"`
	// RuntimeStubSize   uint64 `sizeof:"RuntimeStub"`
	// SafepointBlobSize uint64 `sizeof:"SafepointBlob"`

	CodeCacheStart uint64 `offsetof:"CodeCache._low_bound"`
	CodeCacheEnd   uint64 `offsetof:"CodeCache._high_bound"`

	DeoptHandler uint64 `offsetof:"CompiledMethod.deopt_handler"`

	HeapBlockSize uint64 `sizeof:"HeapBlock"`
	SegmentShift  uint64 `offsetof:"ZLiveMap._segment_shift"`
}

func (oj openjdk) Layout() runtimedata.RuntimeData {
	return &java.Layout{
		CollectedHeapReserve: oj.CollectedHeapReserve,
		MemRegionStart:       oj.MemRegionStart,
		MemRegionEnd:         oj.MemRegionEnd,
		HeapWordSize:         oj.HeapWordSize,

		VMStructEntryTypeName:  oj.VMStructEntryTypeName,
		VMStructEntryFieldName: oj.VMStructEntryFieldName,
		VMStructEntryAddress:   oj.VMStructEntryAddress,
		VMStructEntrySize:      oj.VMStructEntrySize,

		KlassName: oj.KlassName,

		ConstantPoolHolder: oj.ConstantPoolHolder,
		ConstantPoolSize:   oj.ConstantPoolSize,

		OOPDescMetadata: oj.OOPDescMetadata,
		OPPDescSize:     oj.OPPDescSize,

		AccessFlags:             oj.AccessFlags,
		SymbolLengthAndRefcount: oj.SymbolLengthAndRefcount,
		SymbolBody:              oj.SymbolBody,

		MethodConst:       oj.MethodConst,
		MethodAccessFlags: oj.MethodAccessFlags,
		MethodSize:        oj.MethodSize,

		ConstMethodConstants:      oj.ConstMethodConstants,
		ConstMethodFlags:          oj.ConstMethodFlags,
		ConstMethodCodeSize:       oj.ConstMethodCodeSize,
		ConstMethodNameIndex:      oj.ConstMethodNameIndex,
		ConstMethodSignatureIndex: oj.ConstMethodSignatureIndex,
		ConstMethodSize:           oj.ConstMethodSize,

		CodeHeapMemory:          oj.CodeHeapMemory,
		CodeHeapSegmap:          oj.CodeHeapSegmap,
		CodeHeapLog2SegmentSize: oj.CodeHeapLog2SegmentSize,

		VirtualSpaceLowBoundary:  oj.VirtualSpaceLowBoundary,
		VirtualSpaceHighBoundary: oj.VirtualSpaceHighBoundary,
		VirtualSpaceLow:          oj.VirtualSpaceLow,
		VirtualSpaceHigh:         oj.VirtualSpaceHigh,

		CodeBlobName:         oj.CodeBlobName,
		CodeBlobHeaderSize:   oj.CodeBlobHeaderSize,
		CodeBlobContentBegin: oj.CodeBlobContentBegin,

		CodeBlobCodeBegin:  oj.CodeBlobCodeBegin,
		CodeBlobCodeEnd:    oj.CodeBlobCodeEnd,
		CodeBlobDataOffset: oj.CodeBlobDataOffset,
		CodeBlobFrameSize:  oj.CodeBlobFrameSize,

		CodeBlobSize: oj.CodeBlobSize,

		NMethodMethod:             oj.NMethodMethod,
		NMethodDependenciesOffset: oj.NMethodDependenciesOffset,
		NMethodMetadataOffset:     oj.NMethodMetadataOffset,
		NMethodScopesDataBegin:    oj.NMethodScopesDataBegin,
		NMethodScopesPCsOffset:    oj.NMethodScopesPCsOffset,
		NMethodHandlerTableOffset: oj.NMethodHandlerTableOffset,
		NMethodDeoptHandlerBegin:  oj.NMethodDeoptHandlerBegin,
		NMethodOrigPCOffset:       oj.NMethodOrigPCOffset,

		NMethodSize: oj.NMethodSize,

		PCDescPCOffset:          oj.PCDescPCOffset,
		PCDescScopeDecodeOffset: oj.PCDescScopeDecodeOffset,
		PCDescSize:              oj.PCDescSize,

		NarrowPtrStructBase:  oj.NarrowPtrStructBase,
		NarrowPtrStructShift: oj.NarrowPtrStructShift,

		// BufferBlobSize:    oj.BufferBlobSize,
		// SingletonBlobSize: oj.SingletonBlobSize,
		// RuntimeStubSize:   oj.RuntimeStubSize,
		// SafepointBlobSize: oj.SafepointBlobSize,

		CodeCacheStart: oj.CodeCacheStart,
		CodeCacheEnd:   oj.CodeCacheEnd,

		DeoptHandler: oj.DeoptHandler,

		HeapBlockSize: oj.HeapBlockSize,
		SegmentShift:  oj.SegmentShift,
	}
}
