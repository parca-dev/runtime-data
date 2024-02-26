package python

import (
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"github.com/parca-dev/runtime-data/pkg/version"

	"github.com/Masterminds/semver/v3"
)

const doesNotExist = -1

func DataMapForLayout(v string) runtimedata.LayoutMap {
	// Keys are version constraints defined in semver format,
	// check github.com/Masterminds/semver for more details.
	pythonVersions := map[*semver.Constraints]runtimedata.LayoutMap{
		version.MustParseConstraints("2.7.x"):         &python27{},
		version.MustParseConstraints("3.3.x - 3.9.x"): &python33_39{},
		version.MustParseConstraints("3.10.x"):        &python310{},
		version.MustParseConstraints("3.11.x"):        &python311{},
		version.MustParseConstraints("3.12.x"):        &python312{},
		version.MustParseConstraints(">=3.13.x-0"):    &python313{},
	}

	lookupVersion := semver.MustParse(v)
	for c, mapper := range pythonVersions {
		if c.Check(lookupVersion) {
			return mapper
		}
	}
	return nil
}

type python27 struct {
	PyObjectObType               int64 `offsetof:"PyObject.ob_type"`
	PyStringData                 int64 `offsetof:"PyStringObject.ob_sval"`
	PyStringSize                 int64 `offsetof:"PyStringObject.ob_size"`
	PyTypeObjectTpName           int64 `offsetof:"PyTypeObject.tp_name"`
	PyThreadStateInterp          int64 `offsetof:"PyThreadState.interp"`
	PyThreadStateNext            int64 `offsetof:"PyThreadState.next"`
	PyThreadStateFrame           int64 `offsetof:"PyThreadState.frame"`
	PyThreadStateThreadID        int64 `offsetof:"PyThreadState.thread_id"`
	PyInterpreterStateTstateHead int64 `offsetof:"PyInterpreterState.tstate_head"`
	PyFrameObjectFBack           int64 `offsetof:"PyFrameObject.f_back"`
	PyFrameObjectFCode           int64 `offsetof:"PyFrameObject.f_code"`
	PyFrameObjectFLineNo         int64 `offsetof:"PyFrameObject.f_lineno"`
	PyFrameObjectFLocalsplus     int64 `offsetof:"PyFrameObject.f_localsplus"`
	PyCodeObjectCoFilename       int64 `offsetof:"PyCodeObject.co_filename"`
	PyCodeObjectCoName           int64 `offsetof:"PyCodeObject.co_name"`
	PyCodeObjectCoVarNames       int64 `offsetof:"PyCodeObject.co_varnames"`
	PyCodeObjectCoFirstlineno    int64 `offsetof:"PyCodeObject.co_firstlineno"`
	PyTupleObjectObItem          int64 `offsetof:"PyTupleObject.ob_item"`
}

func (p python27) Layout() runtimedata.RuntimeData {
	return &Layout{
		PyObject: PyObject{
			ObType: p.PyObjectObType,
		},
		PyString: PyString{
			Data: p.PyStringData,
			Size: p.PyStringSize,
		},
		PyTypeObject: PyTypeObject{
			TPName: p.PyTypeObjectTpName,
		},
		PyThreadState: PyThreadState{
			Interp:         p.PyThreadStateInterp,
			Next:           p.PyThreadStateNext,
			Frame:          p.PyThreadStateFrame,
			ThreadID:       p.PyThreadStateThreadID,
			NativeThreadID: doesNotExist,
			CFrame:         doesNotExist,
		},
		PyCFrame: PyCFrame{
			CurrentFrame: 0,
		},
		PyInterpreterState: PyInterpreterState{
			TStateHead: p.PyInterpreterStateTstateHead,
		},
		PyRuntimeState: PyRuntimeState{
			InterpMain: doesNotExist,
		},
		PyFrameObject: PyFrameObject{
			FBack:       p.PyFrameObjectFBack,
			FCode:       p.PyFrameObjectFCode,
			FLineno:     p.PyFrameObjectFLineNo,
			FLocalsplus: p.PyFrameObjectFLocalsplus,
		},
		PyCodeObject: PyCodeObject{
			CoFilename:    p.PyCodeObjectCoFilename,
			CoName:        p.PyCodeObjectCoName,
			CoVarnames:    p.PyCodeObjectCoVarNames,
			CoFirstlineno: p.PyCodeObjectCoFirstlineno,
		},
		PyTupleObject: PyTupleObject{
			ObItem: p.PyTupleObjectObItem,
		},
	}
}

type python33_39 struct {
	PyObjectObType               int64 `offsetof:"PyObject.ob_type"`
	PyStringData                 int64 `sizeof:"PyASCIIObject"`
	PyStringSize                 int64 `offsetof:"PyVarObject.ob_size"`
	PyTypeObjectTpName           int64 `offsetof:"PyTypeObject.tp_name"`
	PyThreadStateInterp          int64 `offsetof:"PyThreadState.interp"`
	PyThreadStateNext            int64 `offsetof:"PyThreadState.next"`
	PyThreadStateFrame           int64 `offsetof:"PyThreadState.frame"`
	PyThreadStateThreadID        int64 `offsetof:"PyThreadState.thread_id"`
	PyInterpreterStateTstateHead int64 `offsetof:"PyInterpreterState.tstate_head"`
	PyFrameObjectFBack           int64 `offsetof:"PyFrameObject.f_back"`
	PyFrameObjectFCode           int64 `offsetof:"PyFrameObject.f_code"`
	PyFrameObjectFLineNo         int64 `offsetof:"PyFrameObject.f_lineno"`
	PyFrameObjectFLocalsplus     int64 `offsetof:"PyFrameObject.f_localsplus"`
	PyCodeObjectCoFilename       int64 `offsetof:"PyCodeObject.co_filename"`
	PyCodeObjectCoName           int64 `offsetof:"PyCodeObject.co_name"`
	PyCodeObjectCoVarNames       int64 `offsetof:"PyCodeObject.co_varnames"`
	PyCodeObjectCoFirstlineno    int64 `offsetof:"PyCodeObject.co_firstlineno"`
	PyTupleObjectObItem          int64 `offsetof:"PyTupleObject.ob_item"`
}

func (p python33_39) Layout() runtimedata.RuntimeData {
	return &Layout{
		PyObject: PyObject{
			ObType: p.PyObjectObType,
		},
		PyString: PyString{
			Data: p.PyStringData,
			Size: p.PyStringSize,
		},
		PyTypeObject: PyTypeObject{
			TPName: p.PyTypeObjectTpName,
		},
		PyThreadState: PyThreadState{
			Interp:         p.PyThreadStateInterp,
			Next:           p.PyThreadStateNext,
			Frame:          p.PyThreadStateFrame,
			ThreadID:       p.PyThreadStateThreadID,
			NativeThreadID: doesNotExist,
			CFrame:         doesNotExist,
		},
		PyCFrame: PyCFrame{
			CurrentFrame: 0,
		},
		PyInterpreterState: PyInterpreterState{
			TStateHead: p.PyInterpreterStateTstateHead,
		},
		PyRuntimeState: PyRuntimeState{
			InterpMain: doesNotExist,
		},
		PyFrameObject: PyFrameObject{
			FBack:       p.PyFrameObjectFBack,
			FCode:       p.PyFrameObjectFCode,
			FLineno:     p.PyFrameObjectFLineNo,
			FLocalsplus: p.PyFrameObjectFLocalsplus,
		},
		PyCodeObject: PyCodeObject{
			CoFilename:    p.PyCodeObjectCoFilename,
			CoName:        p.PyCodeObjectCoName,
			CoVarnames:    p.PyCodeObjectCoVarNames,
			CoFirstlineno: p.PyCodeObjectCoFirstlineno,
		},
		PyTupleObject: PyTupleObject{
			ObItem: p.PyTupleObjectObItem,
		},
	}
}

// TODO(kakkoyun): https://github.com/python/cpython/blob/3.10/Include/cpython/unicodeobject.h#L82-L84
type python310 struct {
	PyObjectObType               int64 `offsetof:"PyObject.ob_type"`
	PyStringData                 int64 `sizeof:"PyASCIIObject"`
	PyTypeObjectTpName           int64 `offsetof:"PyTypeObject.tp_name"`
	PyThreadStateInterp          int64 `offsetof:"PyThreadState.interp"`
	PyThreadStateNext            int64 `offsetof:"PyThreadState.next"`
	PyThreadStateFrame           int64 `offsetof:"PyThreadState.frame"`
	PyThreadStateThreadID        int64 `offsetof:"PyThreadState.thread_id"`
	PyInterpreterStateTstateHead int64 `offsetof:"PyInterpreterState.tstate_head"`
	PyFrameObjectFBack           int64 `offsetof:"PyFrameObject.f_back"`
	PyFrameObjectFCode           int64 `offsetof:"PyFrameObject.f_code"`
	PyFrameObjectFLineNo         int64 `offsetof:"PyFrameObject.f_lineno"`
	PyFrameObjectFLocalsplus     int64 `offsetof:"PyFrameObject.f_localsplus"`
	PyCodeObjectCoFilename       int64 `offsetof:"PyCodeObject.co_filename"`
	PyCodeObjectCoName           int64 `offsetof:"PyCodeObject.co_name"`
	PyCodeObjectCoVarNames       int64 `offsetof:"PyCodeObject.co_varnames"`
	PyCodeObjectCoFirstlineno    int64 `offsetof:"PyCodeObject.co_firstlineno"`
	PyTupleObjectObItem          int64 `offsetof:"PyTupleObject.ob_item"`
}

func (p python310) Layout() runtimedata.RuntimeData {
	return &Layout{
		PyObject: PyObject{
			ObType: p.PyObjectObType,
		},
		PyString: PyString{
			Data: p.PyStringData,
			Size: doesNotExist,
		},
		PyTypeObject: PyTypeObject{
			TPName: p.PyTypeObjectTpName,
		},
		PyThreadState: PyThreadState{
			Interp:         p.PyThreadStateInterp,
			Next:           p.PyThreadStateNext,
			Frame:          p.PyThreadStateFrame,
			ThreadID:       p.PyThreadStateThreadID,
			NativeThreadID: doesNotExist,
			CFrame:         doesNotExist,
		},
		PyCFrame: PyCFrame{
			CurrentFrame: 0,
		},
		PyInterpreterState: PyInterpreterState{
			TStateHead: p.PyInterpreterStateTstateHead,
		},
		PyRuntimeState: PyRuntimeState{
			InterpMain: doesNotExist,
		},
		PyFrameObject: PyFrameObject{
			FBack:       p.PyFrameObjectFBack,
			FCode:       p.PyFrameObjectFCode,
			FLineno:     p.PyFrameObjectFLineNo,
			FLocalsplus: p.PyFrameObjectFLocalsplus,
		},
		PyCodeObject: PyCodeObject{
			CoFilename:    p.PyCodeObjectCoFilename,
			CoName:        p.PyCodeObjectCoName,
			CoVarnames:    p.PyCodeObjectCoVarNames,
			CoFirstlineno: p.PyCodeObjectCoFirstlineno,
		},
		PyTupleObject: PyTupleObject{
			ObItem: p.PyTupleObjectObItem,
		},
	}
}

type python311 struct {
	PyObjectObType                    int64 `offsetof:"PyObject.ob_type"`
	PyStringData                      int64 `sizeof:"PyASCIIObject"`
	PyTypeObjectTpName                int64 `offsetof:"PyTypeObject.tp_name"`
	PyThreadStateInterp               int64 `offsetof:"PyThreadState.interp"`
	PyThreadStateNext                 int64 `offsetof:"PyThreadState.next"`
	PyThreadStateThreadID             int64 `offsetof:"PyThreadState.thread_id"`
	PyThreadStateNativeThreadID       int64 `offsetof:"PyThreadState.native_thread_id"`
	PyThreadStateCFrame               int64 `offsetof:"PyThreadState.cframe"`
	PyCFrameCurrentFrame              int64 `offsetof:"_PyCFrame.current_frame"`
	PyInterpreterStateTstateHead      int64 `offsetof:"PyInterpreterState.threads"`
	PyInterpreterStateIsPythreadsHead int64 `offsetof:"pythreads.head"`
	PyRuntimeStateInterpreters        int64 `offsetof:"pyruntimestate.interpreters"`
	PyRuntimeStatePyInterpretersMain  int64 `offsetof:"pyinterpreters.main"`
	PyFrameObjectFBack                int64 `offsetof:"_PyInterpreterFrame.previous"`
	PyFrameObjectFCode                int64 `offsetof:"_PyInterpreterFrame.f_code"`
	PyFrameObjectFLocalsplus          int64 `offsetof:"_PyInterpreterFrame.localsplus"`
	PyCodeObjectCoFilename            int64 `offsetof:"PyCodeObject.co_filename"`
	PyCodeObjectCoName                int64 `offsetof:"PyCodeObject.co_name"`
	PyCodeObjectCoVarNames            int64 `offsetof:"PyCodeObject.co_localsplusnames"`
	PyCodeObjectCoFirstlineno         int64 `offsetof:"PyCodeObject.co_firstlineno"`
	PyTupleObjectObItem               int64 `offsetof:"PyTupleObject.ob_item"`
}

func (p python311) Layout() runtimedata.RuntimeData {
	return &Layout{
		PyObject: PyObject{
			ObType: p.PyObjectObType,
		},
		PyString: PyString{
			Data: p.PyStringData,
			Size: doesNotExist,
		},
		PyTypeObject: PyTypeObject{
			TPName: p.PyTypeObjectTpName,
		},
		PyThreadState: PyThreadState{
			Interp:         p.PyThreadStateInterp,
			Next:           p.PyThreadStateNext,
			Frame:          doesNotExist,
			ThreadID:       p.PyThreadStateThreadID,
			NativeThreadID: p.PyThreadStateNativeThreadID,
			CFrame:         p.PyThreadStateCFrame,
		},
		PyCFrame: PyCFrame{
			CurrentFrame: p.PyCFrameCurrentFrame,
		},
		PyInterpreterState: PyInterpreterState{
			TStateHead: p.PyInterpreterStateTstateHead + p.PyInterpreterStateIsPythreadsHead,
		},
		PyRuntimeState: PyRuntimeState{
			InterpMain: p.PyRuntimeStateInterpreters + p.PyRuntimeStatePyInterpretersMain,
		},
		PyFrameObject: PyFrameObject{
			FBack:       p.PyFrameObjectFBack,
			FCode:       p.PyFrameObjectFCode,
			FLineno:     doesNotExist,
			FLocalsplus: p.PyFrameObjectFLocalsplus,
		},
		PyCodeObject: PyCodeObject{
			CoFilename:    p.PyCodeObjectCoFilename,
			CoName:        p.PyCodeObjectCoName,
			CoVarnames:    p.PyCodeObjectCoVarNames,
			CoFirstlineno: p.PyCodeObjectCoFirstlineno,
		},
		PyTupleObject: PyTupleObject{
			ObItem: p.PyTupleObjectObItem,
		},
	}
}

type python312 struct {
	PyObjectObType                    int64 `offsetof:"PyObject.ob_type"`
	PyStringData                      int64 `sizeof:"PyASCIIObject"`
	PyTypeObjectTpName                int64 `offsetof:"PyTypeObject.tp_name"`
	PyThreadStateInterp               int64 `offsetof:"PyThreadState.interp"`
	PyThreadStateNext                 int64 `offsetof:"PyThreadState.next"`
	PyThreadStateThreadID             int64 `offsetof:"PyThreadState.thread_id"`
	PyThreadStateNativeThreadID       int64 `offsetof:"PyThreadState.native_thread_id"`
	PyThreadStateCFrame               int64 `offsetof:"PyThreadState.cframe"`
	PyInterpreterStateTstateHead      int64 `offsetof:"PyInterpreterState.threads"`
	PyInterpreterStateIsPythreadsHead int64 `offsetof:"pythreads.head"`
	PyRuntimeStateInterpreters        int64 `offsetof:"pyruntimestate.interpreters"`
	PyRuntimeStatePyInterpretersMain  int64 `offsetof:"pyinterpreters.main"`
	PyFrameObjectFBack                int64 `offsetof:"_PyInterpreterFrame.previous"`
	PyFrameObjectFCode                int64 `offsetof:"_PyInterpreterFrame.f_code"`
	PyFrameObjectFLocalsplus          int64 `offsetof:"_PyInterpreterFrame.localsplus"`
	PyCodeObjectCoFilename            int64 `offsetof:"PyCodeObject.co_filename"`
	PyCodeObjectCoName                int64 `offsetof:"PyCodeObject.co_name"`
	PyCodeObjectCoVarNames            int64 `offsetof:"PyCodeObject.co_varnames"`
	PyCodeObjectCoFirstlineno         int64 `offsetof:"PyCodeObject.co_firstlineno"`
	PyTupleObjectObItem               int64 `offsetof:"PyTupleObject.ob_item"`
}

func (p python312) Layout() runtimedata.RuntimeData {
	return &Layout{
		PyObject: PyObject{
			ObType: p.PyObjectObType,
		},
		PyString: PyString{
			Data: p.PyStringData,
			Size: doesNotExist,
		},
		PyTypeObject: PyTypeObject{
			TPName: p.PyTypeObjectTpName,
		},
		PyThreadState: PyThreadState{
			Interp:         p.PyThreadStateInterp,
			Next:           p.PyThreadStateNext,
			Frame:          doesNotExist,
			ThreadID:       p.PyThreadStateThreadID,
			NativeThreadID: p.PyThreadStateNativeThreadID,
			CFrame:         p.PyThreadStateCFrame,
		},
		PyCFrame: PyCFrame{
			CurrentFrame: 0,
		},
		PyInterpreterState: PyInterpreterState{
			TStateHead: p.PyInterpreterStateTstateHead + p.PyInterpreterStateIsPythreadsHead,
		},
		PyRuntimeState: PyRuntimeState{
			InterpMain: p.PyRuntimeStateInterpreters + p.PyRuntimeStatePyInterpretersMain,
		},
		PyFrameObject: PyFrameObject{
			FBack:       p.PyFrameObjectFBack,
			FCode:       p.PyFrameObjectFCode,
			FLineno:     doesNotExist,
			FLocalsplus: p.PyFrameObjectFLocalsplus,
		},
		PyCodeObject: PyCodeObject{
			CoFilename:    p.PyCodeObjectCoFilename,
			CoName:        p.PyCodeObjectCoName,
			CoVarnames:    p.PyCodeObjectCoVarNames,
			CoFirstlineno: p.PyCodeObjectCoFirstlineno,
		},
		PyTupleObject: PyTupleObject{
			ObItem: p.PyTupleObjectObItem,
		},
	}
}

type python313 struct {
	PyObjectObType                    int64 `offsetof:"PyObject.ob_type"`
	PyStringData                      int64 `sizeof:"PyASCIIObject"`
	PyTypeObjectTpName                int64 `offsetof:"PyTypeObject.tp_name"`
	PyThreadStateInterp               int64 `offsetof:"PyThreadState.interp"`
	PyThreadStateNext                 int64 `offsetof:"PyThreadState.next"`
	PyThreadStateThreadID             int64 `offsetof:"PyThreadState.thread_id"`
	PyThreadStateNativeThreadID       int64 `offsetof:"PyThreadState.native_thread_id"`
	PyThreadStateCurrentFrame         int64 `offsetof:"PyThreadState.current_frame"`
	PyInterpreterStateTstateHead      int64 `offsetof:"PyInterpreterState.threads"`
	PyInterpreterStateIsPythreadsHead int64 `offsetof:"pythreads.head"`
	PyRuntimeStateInterpreters        int64 `offsetof:"pyruntimestate.interpreters"`
	PyRuntimeStatePyInterpretersMain  int64 `offsetof:"pyinterpreters.main"`
	PyFrameObjectFBack                int64 `offsetof:"_PyInterpreterFrame.previous"`
	PyFrameObjectFExecutable          int64 `offsetof:"_PyInterpreterFrame.f_executable"`
	PyFrameObjectFLocalsplus          int64 `offsetof:"_PyInterpreterFrame.localsplus"`
	PyCodeObjectCoFilename            int64 `offsetof:"PyCodeObject.co_filename"`
	PyCodeObjectCoName                int64 `offsetof:"PyCodeObject.co_name"`
	PyCodeObjectCoFirstlineno         int64 `offsetof:"PyCodeObject.co_firstlineno"`
	PyCodeObjectCoVarNames            int64 `offsetof:"PyCodeObject.co_varnames"`
	PyTupleObjectObItem               int64 `offsetof:"PyTupleObject.ob_item"`
}

func (p python313) Layout() runtimedata.RuntimeData {
	return &Layout{
		PyObject: PyObject{
			ObType: p.PyObjectObType,
		},
		PyString: PyString{
			Data: p.PyStringData,
			Size: doesNotExist,
		},
		PyTypeObject: PyTypeObject{
			TPName: p.PyTypeObjectTpName,
		},
		PyThreadState: PyThreadState{
			Interp:         p.PyThreadStateInterp,
			Next:           p.PyThreadStateNext,
			Frame:          doesNotExist,
			ThreadID:       p.PyThreadStateThreadID,
			NativeThreadID: p.PyThreadStateNativeThreadID,
			CFrame:         doesNotExist,
		},
		PyCFrame: PyCFrame{
			CurrentFrame: p.PyThreadStateCurrentFrame,
		},
		PyInterpreterState: PyInterpreterState{
			TStateHead: p.PyInterpreterStateTstateHead + p.PyInterpreterStateIsPythreadsHead,
		},
		PyRuntimeState: PyRuntimeState{
			InterpMain: p.PyRuntimeStateInterpreters + p.PyRuntimeStatePyInterpretersMain,
		},
		PyFrameObject: PyFrameObject{
			FBack:       p.PyFrameObjectFBack,
			FCode:       p.PyFrameObjectFExecutable,
			FLineno:     doesNotExist,
			FLocalsplus: p.PyFrameObjectFLocalsplus,
		},
		PyCodeObject: PyCodeObject{
			CoFilename:    p.PyCodeObjectCoFilename,
			CoName:        p.PyCodeObjectCoName,
			CoVarnames:    p.PyCodeObjectCoVarNames,
			CoFirstlineno: p.PyCodeObjectCoFirstlineno,
		},
		PyTupleObject: PyTupleObject{
			ObItem: p.PyTupleObjectObItem,
		},
	}
}
