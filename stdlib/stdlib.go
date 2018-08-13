// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"bitbucket.org/novostrim/go-gentee/core"
)

// InitStdlib appends stdlib types and fucntions to the virtual machine
func InitStdlib(vm *core.VirtualMachine) {
	vm.Units[core.DefName] = core.InitUnit(core.UnitPackage)
	InitTypes(vm)
	InitInt(vm)
	InitBool(vm)
	InitStr(vm)

	vm.Units[core.DefName].NewConst(core.ConstDepth, int64(1000), true)
	vm.Units[core.DefName].NewConst(core.ConstCycle, int64(16000000), true)
	vm.Units[core.DefName].NewConst(core.ConstIota, int64(0), false)
}
