// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"fmt"
	"path/filepath"
	"strings"
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
	// ErrObjValue is returned when obj has wrong type
	ErrObjValue

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
	Line  int    // line position in the source
	Pos   int    // column position in the line
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

// ErrFormat is a function for formating error message
func ErrFormat(path string, line, pos int, message string) string {
	dirs := strings.Split(filepath.ToSlash(path), `/`)
	if len(dirs) > 3 {
		path = `...` + path[len(path)-len(strings.Join(dirs[len(dirs)-3:], `/`))-1:]
	}
	return strings.TrimSpace(fmt.Sprintf(`%s [%d:%d] %s`, path, line, pos, message))

}
