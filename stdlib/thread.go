// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// InitThread appends stdlib thread functions to the virtual machine
func InitThread(vm *core.VirtualMachine) {
	for _, item := range []embedInfo{
		{AssignºThreadThread, `thread,thread`, `thread`},      // thread = thread
		{AssignAddºArrInt, `arr.thread,thread`, `arr.thread`}, // arr += thread
		{closeºThread, `thread`, ``},                          // close( thread )
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AssignºThreadThread assigns one thread to another
func AssignºThreadThread(ptr *interface{}, value int64) int64 {
	*ptr = value
	return (*ptr).(int64)
}

// closeºThread closes the thread
func closeºThread(rt *core.RunTime, threadID int64) error {
	root := rt.Root.Threads
	root.ThreadMutex.Lock()
	defer root.ThreadMutex.Unlock()
	if threadID <= 0 || int64(len(root.Threads)) <= threadID {
		return fmt.Errorf(core.ErrorText(core.ErrThreadIndex))
	}
	root.Threads[threadID].Owner.ToBreak = true
	return nil
}
