// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitMap appends stdlib map functions to the virtual machine
func InitMap(vm *core.VirtualMachine) {
	for _, item := range []embedInfo{
		{LenºMap, `map*`, `int`},                   // the length of map
		{AssignºMapMap, `map*,map*`, `map*`},       // map = map
		{AssignBitAndºMapMap, `map*,map*`, `map*`}, // map &= map
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// LenºMap returns the length of the map
func LenºMap(pmap *core.Map) int64 {
	return int64(len(pmap.Data))
}

// AssignºMapMap copies one array to another one
func AssignºMapMap(ptr *interface{}, value *core.Map) *core.Map {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Map)
}

// AssignBitAndºMapMap assigns a pointer to the map
func AssignBitAndºMapMap(ptr *interface{}, value *core.Map) *core.Map {
	*ptr = value
	return (*ptr).(*core.Map)
}
