// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitFn appends stdlib fn functions to the virtual machine
func InitFn(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{core.Link{AssignºFnFn, core.ASSIGN + 7}, `fn,fn`, `fn`}, // fn = fn
		//		{core.Link{AssignBitAndºFnFn, core.ASSIGNPTR}, `fn*,fn*`, `fn*`}, // fn &= fn
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AssignºFnFn copies one fn to another one
func AssignºFnFn(ptr *interface{}, value *core.Fn) *core.Fn {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Fn)
}

// AssignBitAndºFnFn assigns a pointer to the array
/*func AssignBitAndºFnFn(ptr *interface{}, value *core.Fn) *core.Fn {
	*ptr = value
	return (*ptr).(*core.Fn)
}*/
