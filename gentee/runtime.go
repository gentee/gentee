// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
)

type CallCode struct {
	Code     *Code
	Offset   int // the offset of the current command in the bytecode
	StackOff int // the offset in the stack
}

// RunTime is the structure for running compiled functions
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
	if idFunc >= len(rt.VM.Funcs) {
		return runtimeError(rt, ErrRuntime)
	}
	code := rt.VM.Funcs[idFunc]
	call := CallCode{
		Code:     code,
		Offset:   0,
		StackOff: len(rt.Stack),
	}
	fmt.Println(`RUN`, code.ByteCode)
	rt.Calls = append(rt.Calls, &call)
	for ; call.Offset < len(code.ByteCode); call.Offset++ {
		cmd := code.ByteCode[call.Offset]
		cmdType := cmd.ID >> 24
		cmdID := cmd.ID & 0xFFFFFF
		switch cmdType {
		case cmfStack:
			switch cmdID {
			case cmdPush:
				fmt.Println(`PUSH`, cmd)
				rt.Stack = append(rt.Stack, cmd.Value)
			case cmdReturn:
				fmt.Println(`RET`, cmd)
			}
		}
	}
	return nil
}
