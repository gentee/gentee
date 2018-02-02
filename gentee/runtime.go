// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

type CallCode struct {
	Code     *Code
	Offset   int // the offset of the current command in the bytecode
	StackOff int // the offset in the stack
}

type RunTime struct {
	VM    *VirtualMachine
	Stack []interface{} // the stack of values
	Calls []*CallCode   // the stack of calling functions
}

func newRunTime(vm *VirtualMachine) *RunTime {
	rt := &RunTime{
		VM:    vm,
		Stack: make([]interface{}, 0, 1024),
		Calls: make([]*CallCode, 0, 64),
	}
	return rt
}

func (rt *RunTime) run(idFunc int) error {
	if idFunc >= len(rt.vm.Funcs) {
		return runtimeError(rt, ErrRuntime, `run`)
	}
	call = CallCode{
		Code: code,
		Offset: 0,
		StackOff: len(rt.Stack)
	}
}
