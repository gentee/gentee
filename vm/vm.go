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
	Owner  *VM
	States []State
	// These are stacks for different types
	SInt   [STACKSIZE]int64       // int, char, bool
	SFloat [STACKSIZE]float64     // float
	SStr   [STACKSIZE]string      // str
	SAny   [STACKSIZE]interface{} // all other types
}

// State stores tops of stacks
type State struct {
	topInt   int
	topFloat int
	topStr   int
	topAny   int
}

func (state *State) Get() (int, int, int, int) {
	return state.topInt, state.topFloat, state.topStr, state.topAny
}

func (rt *Runtime) PushState(topInt, topFloat, topStr, topAny int) {
	rt.States = append(rt.States, State{
		topInt:   topInt,
		topFloat: topFloat,
		topStr:   topStr,
		topAny:   topAny,
	})
}

func (rt *Runtime) PopState() (int, int, int, int) {
	state := rt.States[len(rt.States)-1]
	rt.States = rt.States[:len(rt.States)-1]
	return state.topInt, state.topFloat, state.topStr, state.topAny
}

func (vm *VM) RunThread(offset int64) (interface{}, error) {
	rt := &Runtime{
		Owner:  vm,
		States: []State{State{}},
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
