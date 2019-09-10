// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

const (
	// ThQueue means that the thread is in the queue to start
	ThQueue = iota
	// ThWork means that the thread is running
	ThWork
	// ThPaused means that the thread has been suspended
	ThPaused
	// ThWait means that the thread is waiting for the end of another thread
	ThWait
	// ThFinished means that the thread finished
	ThFinished
	// ThError means that the thread has been closed with an error
	ThError
	// ThClosed means that the thread has been closed
	ThClosed
)

const (
	// ThCmdClose closes the thread
	ThCmdClose = iota
	// ThCmdResume resumes the thread
	ThCmdResume
	// ThCmdContinue continues the thread after waiting
	ThCmdContinue
)

// Thread contains information about a thread
type Thread struct {
	Status byte
	Sleep  int64
	Chan   chan int
	Notify []int64 // who waits the end
}

/*
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
	if rt.Root.Thread != nil && rt.Root.Thread.Status >= ThFinished {
		return false
	}
	rt.Thread = &Thread{
		Status: status,
		Chan:   make(chan int, 8),
	}
	root.Threads = append(root.Threads, rt.Thread)
	rt.ThreadID = int64(len(root.Threads) - 1)
	if status == ThQueue {
		root.Count++
	}
	return true
}

func (rt *RunTime) setStatus(status byte) {
	root := rt.Root.Threads
	root.ThreadMutex.Lock()
	rt.Thread.Status = status
	root.ThreadMutex.Unlock()
}
*/

func (vm *VM) newThread(status byte) *Runtime {
	if len(vm.Runtimes) > 0 && vm.Runtimes[0].Thread.Status >= ThFinished {
		return nil
	}
	rt := &Runtime{
		Owner: vm,
		Thread: Thread{
			Status: status,
			Chan:   make(chan int, 8),
		},
	}
	vm.ThreadMutex.Lock()
	defer vm.ThreadMutex.Unlock()
	vm.Runtimes = append(vm.Runtimes, rt)
	rt.ThreadID = int64(len(vm.Runtimes) - 1)
	if status == ThQueue {
		vm.Count++
	}
	return rt
}

func (vm *VM) closeAll() {
	vm.ThreadMutex.Lock()
	for i := range vm.Runtimes {
		if vm.Runtimes[i].Thread.Status < ThFinished {
			vm.Runtimes[i].Thread.Chan <- ThCmdClose
		}
	}
	vm.ThreadMutex.Unlock()
}

// GoThread executes a new thread
func (rt *Runtime) GoThread(offset int64) int64 {
	/*	for i := 0; i < funcObj.Block.ParCount; i++ {
		var par interface{}
		CopyVar(&par, rt.Stack[len(rt.Stack)-funcObj.Block.ParCount+i])
		thread.Stack = append(thread.Stack, par)
	}*/
	thread := rt.Owner.newThread(ThQueue)
	if thread == nil {
		return -1
	}
	go func() {
		thread.Thread.Status = ThWork

		_, err := thread.Run(offset)
		rt.Owner.ThreadMutex.Lock()
		if err != nil {
			if thread.Thread.Status != ThClosed {
				thread.Thread.Status = ThError
				rt.Owner.ChError <- err
			}
		} else {
			thread.Thread.Status = ThFinished
		}
		close(thread.Thread.Chan)
		for _, nfyid := range thread.Thread.Notify {
			if rt.Owner.Runtimes[nfyid].Thread.Status == ThWait {
				rt.Owner.Runtimes[nfyid].Thread.Chan <- ThCmdContinue
			}
		}
		rt.Owner.ThreadMutex.Unlock()
		rt.Owner.ChCount <- 1
	}()
	return thread.ThreadID
}
