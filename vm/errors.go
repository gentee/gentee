// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee/core"
)

const (
	// The list of errors

	// ErrRunIndex is returned when Run has been executed with wrong index
	ErrRunIndex = iota + 1
	// ErrDepth is returned when maximum depth of recursion has been reached
	ErrDepth
	// ErrDivZero is returned when there is division by zero
	ErrDivZero
	// ErrCycle is returned when maximum cycle count has been reached
	ErrCycle
	// ErrShift is returned when << or >> are used with the negative right operand
	ErrShift
	// ErrStrToInt is returned when the conversion string to integer is invalid
	ErrStrToInt
	// ErrStrToFloat is returned when the conversion string to float is invalid
	ErrStrToFloat
	// ErrEmptyCommand is returned if empty command is specified in $
	ErrEmptyCommand
	// ErrQuoteCommand is returned if there is an unclosed quotation mark in $ command
	ErrQuoteCommand
	// ErrIndexOut means that int index is out of the length of the array
	ErrIndexOut
	// ErrMapIndex is returned when there is not the key in the map
	ErrMapIndex
	// ErrAssignment is returned when there is a recursive assignment
	ErrAssignment
	// ErrUndefined means that the value of the variable is undefined
	ErrUndefined
	// ErrByteOut is returned when value for buf is greater 255
	ErrByteOut
	// ErrInvalidParam is returned when the function gets invalid parameter(s)
	ErrInvalidParam
	// ErrNotRun is returned when the executing unit doesn't have run function
	ErrNotRun
	// ErrFnEmpty is returned in case of calling undefined fn variable
	ErrFnEmpty
	// ErrThreadIndex is returned when the index of the thread is out of range
	ErrThreadIndex
	// ErrThreadClosed is generated when the thread has been closed
	ErrThreadClosed
	// ErrPlatform is generated when the function is not supported on the current platform
	ErrPlatform

	// ErrEmbedded means golang error in embedded functions
	ErrEmbedded = 254
	// ErrRuntime error. It means bug
	ErrRuntime = 255
)

// TraceInfo is a structure for stack func info
type TraceInfo struct {
	Path  string // the full path name of the source
	Entry string // the entry function name
	Func  string // the called function
	Line  int64  // line position in the source
	Pos   int64  // column position in the line
}

// RuntimeError is a runtime error type
type RuntimeError struct {
	ID      int
	Message string
	Trace   []TraceInfo
}

func (re *RuntimeError) Error() string {
	if len(re.Trace) == 0 {
		return re.Message
	}
	si := re.Trace[len(re.Trace)-1]
	return ErrFormat(si.Path, si.Line, si.Pos, re.Message)
}

var (
	errText = map[int]string{
		ErrRunIndex:     `invalid name of Run`,
		ErrDepth:        `maximum depth of recursion has been reached`,
		ErrDivZero:      `divided by zero`,
		ErrCycle:        `maximum cycle count has been reached`,
		ErrShift:        `right operand of shift cannot be less than zero`,
		ErrStrToInt:     `converting string to integer is invalid`,
		ErrStrToFloat:   `converting string to float is invalid`,
		ErrEmptyCommand: `empty $ command`,
		ErrQuoteCommand: `unclosed quotation mark in $ command`,
		ErrIndexOut:     `index out of range`,
		ErrMapIndex:     `there is not key in the map`,
		ErrAssignment:   `there is a recursive or self assignment`,
		ErrUndefined:    `undefined value`,
		ErrByteOut:      `byte value is out of range`,
		ErrInvalidParam: `invalid value of parameter(s)`,
		ErrNotRun:       `there is not run function`,
		ErrFnEmpty:      `fn variable has not been defined`,
		ErrThreadIndex:  `invalid thread`,
		ErrThreadClosed: `thread has been closed`,
		ErrPlatform:     `unsupported platform`,

		ErrRuntime: `you have found a runtime bug. Let us know, please`,
	}
)

// ErrFormat is a function for formating error message
func ErrFormat(path string, line, pos int64, message string) string {
	dirs := strings.Split(filepath.ToSlash(path), `/`)
	if len(dirs) > 3 {
		path = `...` + path[len(path)-len(strings.Join(dirs[len(dirs)-3:], `/`))-1:]
	}
	return strings.TrimSpace(fmt.Sprintf(`%s [%d:%d] %s`, path, line, pos, message))

}

// ErrorText returns the text of the error message
func ErrorText(id int) string {
	return errText[id]
}

// GetTrace returns information about called functions
func GetTrace(rt *Runtime, pos int64) []TraceInfo {
	var (
		entry string
	)

	ret := make([]TraceInfo, 0, 16)
	newTrace := func(offset int32) {
		for _, ipos := range rt.Owner.Exec.Pos {
			if ipos.Offset >= offset {
				ret = append(ret, TraceInfo{
					Path:  rt.Owner.Exec.Strings[ipos.Path],
					Entry: entry,
					Func:  rt.Owner.Exec.Strings[ipos.Name],
					Line:  int64(ipos.Line),
					Pos:   int64(ipos.Column),
				})
				break
			}
		}
	}
	for _, call := range rt.Calls {
		if !call.IsFunc {
			continue
		}
		newTrace(call.Offset)
	}
	newTrace(int32(pos))

	/*	for _, cmd := range rt.Calls {
		if cmd == nil {
			continue
		}
		obj := cmd.GetObject()
		if obj == nil {
			continue
		}
		if plex := obj.GetLex(); plex != nil {
			lex = plex
		}

		if obj.GetType() == ObjFunc || obj.GetType() == ObjEmbedded {
			if len(entry) == 0 && obj.GetType() == ObjFunc {
				entry = obj.GetName()
				if entry[0] == '*' {
					entry = `thread`
				}
				continue
			}
			line, column = lex.LineColumn(cmd.GetToken())
			ret = append(ret, TraceInfo{
				Path:  lex.Path,
				Entry: entry,
				Func:  obj.GetName(),
				Line:  line,
				Pos:   column,
			})
			last = cmd
			if obj.GetType() == ObjFunc {
				entry = ``
			}
		}
	}*/
	/*	if cmd != nil && cmd != last {
		line, column = lex.LineColumn(cmd.GetToken())
		ret = append(ret, TraceInfo{
			Path:  lex.Path,
			Entry: entry,
			Func:  ``,
			Line:  line,
			Pos:   column,
		})
	}*/
	return ret
}

func runtimeError(rt *Runtime, pos int64, err interface{}, labels ...interface{}) error {
	var (
		errText string
		idError int
	)
	switch v := err.(type) {
	case int:
		errText = ErrorText(v)
		idError = v
	case *RuntimeError:
		errText = v.Message
		idError = v.ID
	case error:
		errText = v.Error()
		idError = ErrEmbedded
	}
	for _, item := range labels {
		errText += fmt.Sprintf(` [%v]`, item)
	}
	return &RuntimeError{
		ID:      idError,
		Message: errText,
		Trace:   GetTrace(rt, pos),
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
