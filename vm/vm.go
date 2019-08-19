// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"github.com/gentee/gentee/core"
)

const (
	STACKSIZE = 128
	// CYCLE is the limit of loops
	CYCLE = uint64(16000000)
	// DEPTH is the maximum size of blocks stack
	DEPTH = uint32(1000)
)

type Settings struct {
	CmdLine []string
	Cycle   uint64 // limit of loops
	Depth   uint32 // limit of blocks stack
}

type Const struct {
	Type  uint16
	Value interface{}
}

// VM is the main structure of the virtual machine
type VM struct {
	Settings Settings
	Exec     *core.Exec
	Consts   map[int32]Const
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
	Cycle  uint64
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
		Consts:   make(map[int32]Const),
	}
	if vm.Settings.Cycle == 0 {
		vm.Settings.Cycle = CYCLE
	}
	if vm.Settings.Depth == 0 {
		vm.Settings.Depth = DEPTH
	}
	//	fmt.Println(`CODE`, vm.Exec.Code)
	//fmt.Println(`POS`, vm.Exec.Pos)
	//fmt.Println(`STRING`, vm.Exec.Strings)
	for i, id := range vm.Exec.Init {
		if i == 0 {
			vm.Consts[id] = Const{Type: core.TYPEINT, Value: int64(0)}
			continue
		}
		val, err := vm.RunThread(int64(vm.Exec.Funcs[id]))
		if err != nil {
			return nil, err
		}
		var constType uint16
		switch v := val.(type) {
		case int64:
			constType = core.TYPEINT
		case bool:
			constType = core.TYPEBOOL
			if v {
				val = int64(1)
			} else {
				val = int64(0)
			}
			//				case reflect.TypeOf(float64(0.0)):
			//					retType = core.STACKFLOAT
		case rune:
			constType = core.TYPECHAR
			val = int64(v)
		case string:
			constType = core.TYPESTR
		}
		vm.Consts[id] = Const{Type: constType, Value: val}
	}
	//	fmt.Println(`CONST`, vm.Consts)
	return vm.RunThread(0)
}
