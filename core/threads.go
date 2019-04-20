// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"sync"
)

const (
	// ThQueue means that the thread is in the queue to start
	ThQueue = iota
	// ThWork means that the thread is running
	ThWork
	// ThFinish means that the thread finished
	ThFinish
	// ThClosed means that the thread has been closed
	ThClosed
)

// Thread contains information about a thread
type Thread struct {
	Status byte
	//	Chan   chan int
	Owner *RunTime // The runtime of the thread
}

// RootThread is a structure for thread management
type RootThread struct {
	ConstMutex  sync.RWMutex
	CtxMutex    sync.RWMutex
	ThreadMutex sync.RWMutex
	Context     map[string]string
	Threads     []*Thread
	Count       int64 // count of active threads
	ChCount     chan int64
	ChError     chan error
}

func (rt *RunTime) newRootThread() {
	rt.Threads = &RootThread{
		Context: make(map[string]string),
		Threads: make([]*Thread, 0, 32),
		ChCount: make(chan int64, 16),
		ChError: make(chan error, 16),
	}
	rt.newThread(ThWork)
	go func() {
		x := int64(1)
		for x != 0 {
			select {
			case x = <-rt.Threads.ChCount:
				if x != 0 {
					rt.Threads.ThreadMutex.Lock()
					rt.Threads.Count--
					rt.Threads.ThreadMutex.Unlock()
				}
			}
		}
	}()
}

func (rt *RunTime) newThread(status byte) bool {
	root := rt.Root.Threads
	root.ThreadMutex.Lock()
	defer root.ThreadMutex.Unlock()
	if rt.Root.ToBreak {
		return false
	}
	ret := &Thread{
		Status: status,
		//		Chan:   make(chan int),
		Owner: rt,
	}
	root.Threads = append(root.Threads, ret)
	rt.ThreadID = int64(len(root.Threads) - 1)
	if status == ThQueue {
		root.Count++
	}
	return true
}

func (rt *RunTime) changeStatus(status byte) {
	root := rt.Root.Threads
	root.ThreadMutex.Lock()
	root.Threads[rt.ThreadID].Status = status
	root.ThreadMutex.Unlock()
}

func (rt *RunTime) closeAll() {
	root := rt.Threads
	root.ThreadMutex.Lock()
	for i := range rt.Threads.Threads {
		root.Threads[i].Owner.ToBreak = true
	}
	root.ThreadMutex.Unlock()
}

// Thread executes a new thread
func (rt *RunTime) Thread(funcObj *FuncObject) int64 {
	thread := &RunTime{
		VM:    rt.VM,
		Stack: make([]interface{}, 0, 1024),
		Calls: make([]ICmd, 0, 64),
		Root:  rt.Root,
		Cycle: rt.Cycle,
		Depth: rt.Depth,
	}
	thread.newThread(ThQueue)
	go func() {
		thread.changeStatus(ThWork)
		if err := thread.runCmd(&funcObj.Block); err != nil {
			if thread.ToBreak {
				thread.changeStatus(ThClosed)
			} else {
				rt.Root.Threads.ChError <- err
			}
		} else {
			thread.changeStatus(ThFinish)
		}
		rt.Root.Threads.ChCount <- 1
	}()
	return thread.ThreadID
}

// GetConst returns the value of the constant
func (rt *RunTime) GetConst(cmd ICmd) (err error) {
	name := cmd.GetObject().GetName()

	rt.Root.Threads.ConstMutex.RLock()
	v, ok := rt.Root.Consts[name]
	if ok {
		rt.Stack = append(rt.Stack, v)
		rt.Root.Threads.ConstMutex.RUnlock()
	} else {
		// TODO: Insert redefinition of constants here
		constObj := cmd.GetObject().(*ConstObject)
		if constObj.Iota != NotIota {
			rt.Root.Consts[ConstIota] = constObj.Iota
		}
		rt.Root.Threads.ConstMutex.RUnlock()
		if err = rt.runCmd(constObj.Exp); err != nil {
			return err
		}
		rt.Root.Threads.ConstMutex.Lock()
		rt.Root.Consts[name] = rt.Stack[len(rt.Stack)-1]
		rt.Root.Threads.ConstMutex.Unlock()
	}
	return nil
}
