// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import "github.com/gentee/gentee/core"

const (
	STACKSIZE = 128
)

type Settings struct {
	CmdLine []string
}

// VM is the main structure of the virtual machine
type VM struct {
	Settings Settings
	Exec     *core.Exec
	Runtimes []*Runtime
}

// Runtime is the one thread structure
type Runtime struct {
	Owner    *VM
	ParCount int32
	Calls    []Call
	//	Consts
	// These are stacks for different types
	SInt   [STACKSIZE]int64       // int, char, bool
	SFloat [STACKSIZE]float64     // float
	SStr   [STACKSIZE]string      // str
	SAny   [STACKSIZE]interface{} // all other types
}

// Call stores stack of blocks
type Call struct {
	IsFunc bool
	Offset int32
	Int    int32
	Float  int32
	Str    int32
	Any    int32
}

func (vm *VM) RunThread(offset int64) (interface{}, error) {
	rt := &Runtime{
		Owner: vm,
	}
	vm.Runtimes = append(vm.Runtimes, rt)

	return rt.Run(offset)
}

func Run(exec *core.Exec, settings Settings) (interface{}, error) {
	vm := &VM{
		Settings: settings,
		Exec:     exec,
	}
	return vm.RunThread(0)
}
