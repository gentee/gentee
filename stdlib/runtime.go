// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitRuntime appends stdlib runtime functions to the virtual machine
func InitRuntime(vm *core.VirtualMachine) {
	NewStructType(vm, `trace`, []string{
		`Path:str`, `Entry:str`, `Func:str`, `Line:int`, `Pos:int`,
	})
	for _, item := range []embedInfo{
		{TraceºTrace, `arr.trace`, `arr.trace`}, // Trace( trace )
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// TraceºTrace gets trace information
func TraceºTrace(rt *core.RunTime, it *core.Array) *core.Array {
	for _, item := range core.GetTrace(rt, nil) {
		trace := core.NewStruct(rt.VM.StdLib().Names[`trace`].(*core.TypeObject))
		trace.Values[0] = item.Path
		trace.Values[1] = item.Entry
		trace.Values[2] = item.Func
		trace.Values[3] = int64(item.Line)
		trace.Values[4] = int64(item.Pos)
		it.Data = append(it.Data, trace)
	}
	return it
}
