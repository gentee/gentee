// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitStdlib appends stdlib types and functions to the virtual machine
func InitStdlib(vm *core.VirtualMachine) {
	stdlib := vm.InitUnit()
	stdlib.Pub = core.PubAll
	vm.Units = append(vm.Units, stdlib)
	vm.UnitNames[core.DefName] = len(vm.Units) - 1
	InitTypes(vm)
	InitInt(vm)
	InitFloat(vm)
	InitBool(vm)
	InitChar(vm)
	InitStr(vm)
	InitKeyValue(vm)
	InitRange(vm)
	InitArray(vm)
	InitBuffer(vm)
	InitMap(vm)
	InitStruct(vm)
	InitFn(vm)
	InitSystem(vm)
	InitTime(vm)
	InitFile(vm)
	InitConsole(vm)
	InitRuntime(vm)
	InitRegExp(vm)

	stdlib.NewConst(core.ConstDepth, int64(1000), true)
	stdlib.NewConst(core.ConstCycle, int64(16000000), true)
	stdlib.NewConst(core.ConstIota, int64(0), false)
	stdlib.NewConst(core.ConstVersion, core.Version, false)
}
