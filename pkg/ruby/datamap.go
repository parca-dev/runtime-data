package ruby

import (
	"github.com/Masterminds/semver/v3"

	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"github.com/parca-dev/runtime-data/pkg/version"
)

func DataMapForLayout(v string) runtimedata.LayoutMap {
	// Keys are version constraints defined in semver format,
	// check github.com/Masterminds/semver for more details.
	rubyVersions := map[*semver.Constraints]runtimedata.LayoutMap{
		version.MustParseConstraints("2.6.x - 2.7.x"): &ruby26_27{},
		version.MustParseConstraints("3.x"):           &ruby30{},
	}
	lookupVersion := semver.MustParse(v)
	for c, mapper := range rubyVersions {
		if c.Check(lookupVersion) {
			return mapper
		}
	}
	return nil
}

type ruby26_27 struct {
	VMOffset                   int64 `offsetof:"rb_execution_context_struct.vm_stack"`
	VMSizeOffset               int64 `offsetof:"rb_execution_context_struct.vm_stack_size"`
	ControlFrameSizeof         int64 `sizeof:"rb_control_frame_struct"`
	CFPOffset                  int64 `offsetof:"rb_execution_context_struct.cfp"`
	LabelOffset                int64 `offsetof:"rb_iseq_location_struct.label"`
	LineInfoTableOffset        int64 `offsetof:"rb_iseq_constant_body.insns_info"`
	LineInfoIseqInfoSizeOffset int64 `offsetof:"iseq_insn_info.size"`
	MainThreadOffset           int64 `offsetof:"rb_vm_struct.main_thread"`
	EcOffset                   int64 `offsetof:"rb_thread_struct.ec"`
}

func (r ruby26_27) Layout() runtimedata.RuntimeData {
	return &Layout{
		VMOffset:            r.VMOffset,
		VMSizeOffset:        r.VMSizeOffset,
		ControlFrameSizeof:  r.ControlFrameSizeof,
		CfpOffset:           r.CFPOffset,
		LabelOffset:         r.LabelOffset,
		PathFlavour:         1,
		LineInfoTableOffset: r.LineInfoTableOffset,
		LineInfoSizeOffset:  r.LineInfoTableOffset + r.LineInfoIseqInfoSizeOffset,
		MainThreadOffset:    r.MainThreadOffset,
		EcOffset:            r.EcOffset,
	}
}

type ruby30 struct {
	VMOffset                   int64 `offsetof:"rb_execution_context_struct.vm_stack"`
	VMSizeOffset               int64 `offsetof:"rb_execution_context_struct.vm_stack_size"`
	ControlFrameSizeof         int64 `sizeof:"rb_control_frame_struct"`
	CFPOffset                  int64 `offsetof:"rb_execution_context_struct.cfp"`
	LabelOffset                int64 `offsetof:"rb_iseq_location_struct.label"`
	LineInfoTableOffset        int64 `offsetof:"rb_iseq_constant_body.insns_info"`
	LineInfoIseqInfoSizeOffset int64 `offsetof:"iseq_insn_info.size"`
	MainThreadOffset           int64 `offsetof:"rb_vm_struct.ractor.main_thread"`
	EcOffset                   int64 `offsetof:"rb_ractor_struct.threads.running_ec"`
}

func (r ruby30) Layout() runtimedata.RuntimeData {
	return &Layout{
		VMOffset:            r.VMOffset,
		VMSizeOffset:        r.VMSizeOffset,
		ControlFrameSizeof:  r.ControlFrameSizeof,
		CfpOffset:           r.CFPOffset,
		LabelOffset:         r.LabelOffset,
		PathFlavour:         1,
		LineInfoTableOffset: r.LineInfoTableOffset,
		LineInfoSizeOffset:  r.LineInfoTableOffset + r.LineInfoIseqInfoSizeOffset,
		// TODO(kakkoyun): This is a temporary fix, we need to find a better way to get the main thread.
		// - https://github.com/javierhonduco/rbperf/issues/78
		MainThreadOffset:    r.MainThreadOffset - 8, 		// ruby_current_vm_ptr->ractor->main_thread
		EcOffset: r.EcOffset, // ruby_current_vm_ptr->ractor->main_thread->ractor(->threads)->running_ec
	}
}
