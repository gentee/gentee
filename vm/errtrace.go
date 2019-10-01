// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// errorºIntStr throws an error
func errorºIntStr(code int64, text string, pars ...interface{}) error {
	if len(pars) > 0 {
		text = fmt.Sprintf(text, pars...)
	}
	return &RuntimeError{
		ID:      int(code),
		Message: text,
	}
}

// ErrID returns the id of the error
func ErrID(err *RuntimeError) int64 {
	return int64(err.ID)
}

// ErrText returns the text of the error
func ErrText(err *RuntimeError) string {
	return err.Message
}

func getTrace(rt *Runtime, list []TraceInfo, it *core.Array) *core.Array {
	for _, item := range list {
		trace := NewStruct(rt, &rt.Owner.Exec.Structs[TRACESTRUCT])
		trace.Values[0] = item.Path
		trace.Values[1] = item.Entry
		trace.Values[2] = item.Func
		trace.Values[3] = int64(item.Line)
		trace.Values[4] = int64(item.Pos)
		it.Data = append(it.Data, trace)
	}
	return it
}

// ErrTrace returns the trace of the error
func ErrTrace(rt *Runtime, err *RuntimeError) *core.Array {
	return getTrace(rt, err.Trace, core.NewArray())
}

// Trace gets trace information
func Trace(rt *Runtime) *core.Array {
	return getTrace(rt, GetTrace(rt, -1), core.NewArray())
}
