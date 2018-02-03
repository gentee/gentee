// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

// Run executes run block
func (vm *VirtualMachine) Run() (interface{}, error) {
	rt := newRunTime(vm)
	if vm.RunID == Undefined {
		return nil, runtimeError(rt, ErrNoRun)
	}
	rt.run(vm.RunID)
	var result interface{}
	if len(rt.Stack) > 0 {
		result = rt.Stack[len(rt.Stack)-1]
	}
	return result, nil
}
