// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitStruct appends stdlib map functions to the virtual machine
func InitStruct(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{AssignºStructStruct, `struct,struct`, `struct`},       // struct = struct
		{AssignBitAndºStructStruct, `struct,struct`, `struct`}, // struct &= struct
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AssignºStructStruct copies one struct to another one
func AssignºStructStruct(ptr *interface{}, value *core.Struct) *core.Struct {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Struct)
}

// AssignBitAndºStructStruct assigns a pointer to data of one struct to another struct
func AssignBitAndºStructStruct(ptr *interface{}, value *core.Struct) *core.Struct {
	*ptr = value
	return (*ptr).(*core.Struct)
}
