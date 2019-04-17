// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitThread appends stdlib thread functions to the virtual machine
func InitThread(vm *core.VirtualMachine) {
	for _, item := range []embedInfo{
		{AssignºThreadThread, `thread,thread`, `thread`},      // thread = thread
		{AssignAddºArrInt, `arr.thread,thread`, `arr.thread`}, // arr += thread
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AssignºThreadThread assigns one thread to another
func AssignºThreadThread(ptr *interface{}, value int64) int64 {
	*ptr = value
	return (*ptr).(int64)
}
