// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitArray appends stdlib array functions to the virtual machine
func InitArray(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		AssignAddºArrStr, // arr += str
		//		AssignAddºArrInt,  // arr += int
		//		AssignAddºArrBool, // arr += bool
		LenºArr, // the length of array
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// LenºArr returns the length of the array
func LenºArr(arr *core.Array) int64 {
	return int64(len(arr.Data))
}

// AssignAddºArrStr appends one string to array
func AssignAddºArrStr(ptr *interface{}, value string) *core.Array {
	(*ptr).(*core.Array).Data = append((*ptr).(*core.Array).Data, value)
	return (*ptr).(*core.Array)
}
