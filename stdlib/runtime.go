// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// InitRuntime appends stdlib runtime functions to the virtual machine
func InitRuntime(vm *core.VirtualMachine) {
	NewStructType(vm, `trace`, []string{
		`Path:str`, `Entry:str`, `Func:str`, `Line:int`, `Pos:int`,
	})
	for _, item := range []embedInfo{
		{errorºIntStr, `int,str`, ``},    // error( int, str )
		{ErrID, `error`, `int`},          // ErrID( error ) int
		{ErrText, `error`, `str`},        // ErrText( error ) str
		{ErrTrace, `error`, `arr.trace`}, // ErrTrace( error ) arr.trace
		{Trace, ``, `arr.trace`},         // Trace() arr.trace
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// errorºIntStr throws an error
func errorºIntStr(code int64, text string, pars ...interface{}) error {
	if len(pars) > 0 {
		text = fmt.Sprintf(text, pars...)
	}
	return &core.RuntimeError{
		ID:      int(code),
		Message: text,
	}
}

// ErrID returns the id of the error
func ErrID(err *core.RuntimeError) int64 {
	return int64(err.ID)
}

// ErrText returns the text of the error
func ErrText(err *core.RuntimeError) string {
	return err.Message
}

func getTrace(rt *core.RunTime, list []core.TraceInfo, it *core.Array) *core.Array {
	for _, item := range list {
		trace := core.NewStruct(rt.VM.StdLib().FindType(`trace`).(*core.TypeObject))
		trace.Values[0] = item.Path
		trace.Values[1] = item.Entry
		trace.Values[2] = item.Func
		trace.Values[3] = int64(item.Line)
		trace.Values[4] = int64(item.Pos)
		it.Data = append(it.Data, trace)
	}
	return it
}

// ErrTrace returns the trace of the error
func ErrTrace(rt *core.RunTime, err *core.RuntimeError) *core.Array {
	return getTrace(rt, err.Trace, core.NewArray())
}

// Trace gets trace information
func Trace(rt *core.RunTime) *core.Array {
	return getTrace(rt, core.GetTrace(rt, nil), core.NewArray())
}
