// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

type Settings struct {
	CmdLine []string
}

// VM is the main structure of the virtual machine
type VM struct {
	Settings *Settings
	BCode    []int16
	Runtimes []Runtime
}

// Runtime is the one thread structure
type Runtime struct {
	Owner *VM
}

func (vm *VM) RunThread(offset int64) (interface{}, error) {
	vm.Runtimes = append(vm.Runtimes, Runtime{Owner: vm})
	return nil, nil
}

func Run(bcode []int16, settings *Settings) (interface{}, error) {
	vm := &VM{
		Settings: settings,
		BCode:    bcode,
	}
	return vm.RunThread(0)
}
