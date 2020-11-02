// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

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
	// ErrCustom is generated when there is an invalid custom declaration
	ErrCustom
	// ErrCRC is returned when Exec was compiled with different stdlib or custom functions
	ErrCRC
	// ErrMainThread is returned when the function is called in go thread
	ErrMainThread
	// ErrThread is returned when the function is called in main thread
	ErrThread
	// ErrObjNil is returned when the object value is undefined
	ErrObjNil
	// ErrObjArr is returned when the object must contains array
	ErrObjArr
	// ErrObjMap is returned when the object must contains map
	ErrObjMap
	// ErrObjType is returned when the value has incompatible type
	ErrObjType
	// ErrTerminated is returned when the script has been terminated by owner
	ErrTerminated
	// ErrExit is returned when exit function has been called
	ErrExit
	// ErrPlayCycle is returned when maximum cycle count has been reached in Playground mode
	ErrPlayCycle
	// ErrPlayRun is returned on starting any processes in Playground mode
	ErrPlayRun
	// ErrPlayEnv is returned on setting environment in Playground mode
	ErrPlayEnv
	// ErrPlayAccess is returned on access denied error in Playground mode
	ErrPlayAccess
	// ErrPlayCount is returned when files count limit is exceeded in Playground mode
	ErrPlayCount
	// ErrPlaySize is returned when the file size limit reached in Playground mode
	ErrPlaySize
	// ErrPlayAllSize is returned when the summary files size limit reached in Playground mode
	ErrPlayAllSize
	// ErrPlayDepth is returned when maximum depth of recursion has been reached om Playground mode
	ErrPlayDepth
	// ErrPlayFunc is returned on calling disabled function in Playground mode
	ErrPlayFunc

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
		ErrObjValue:     `value of the object has wrong type`,
		ErrCustom:       `invalid custom declaration`,
		ErrCRC:          `different CRC of stdlib or custom functions`,
		ErrMainThread:   `%s must be called in the main thread`,
		ErrThread:       `%s cannot be called in the main thread`,
		ErrObjNil:       `obj is undefined (nil)`,
		ErrObjArr:       `obj is not array`,
		ErrObjMap:       `obj is not map`,
		ErrObjType:      `type is incompatible to object`,
		ErrTerminated:   `code execution has been terminated`,
		ErrExit:         `exit`,
		ErrPlayCycle:    `[Playground] maximum cycle count has been reached`,
		ErrPlayRun:      `[Playground] starting any processes is disabled`,
		ErrPlayEnv:      `[Playground] setting the environment variable is disabled`,
		ErrPlayAccess:   `[Playground] access denied`,
		ErrPlayCount:    `[Playground] file limit reached`,
		ErrPlaySize:     `[Playground] file size limit reached`,
		ErrPlayAllSize:  `[Playground] summary files size limit reached`,
		ErrPlayDepth:    `[Playground] maximum depth of recursion has been reached`,
		ErrPlayFunc:     `[Playground] calling the %s function is prohibited`,

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
	if rt.ThreadID == 0 {
		entry = `run`
	} else {
		entry = `thread`
	}
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
				entry = rt.Owner.Exec.Strings[ipos.Name]
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
	if pos >= 0 {
		newTrace(int32(pos))
	}

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
