// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitStdlib appends stdlib types and fucntions to the virtual machine
func InitStdlib(vm *core.VirtualMachine) {
	vm.Units = append(vm.Units, core.InitUnit(core.UnitPackage))
	vm.Names[core.DefName] = len(vm.Units) - 1
	InitTypes(vm)
	InitInt(vm)
	InitBool(vm)
	InitChar(vm)
	InitStr(vm)
	InitKeyValue(vm)
	InitRange(vm)
	InitArray(vm)
	InitMap(vm)
	InitStruct(vm)
	InitSystem(vm)

	vm.StdLib().NewConst(core.ConstDepth, int64(1000), true)
	vm.StdLib().NewConst(core.ConstCycle, int64(16000000), true)
	vm.StdLib().NewConst(core.ConstIota, int64(0), false)
	vm.StdLib().NewConst(core.ConstVersion, core.Version, false)
}
