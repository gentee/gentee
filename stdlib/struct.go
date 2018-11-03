// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitStruct appends stdlib map functions to the virtual machine
func InitStruct(vm *core.VirtualMachine) {
	for _, item := range []embedInfo{
		{AssignºStructStruct, `struct,struct`, `struct`}, // struct = struct
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AssignºStructStruct copies one struct to another one
func AssignºStructStruct(ptr *interface{}, value *core.Struct) *core.Struct {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Struct)
}
