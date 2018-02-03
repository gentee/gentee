// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
)

const (
	// The list of errors

	// ErrLetter returns when an unknown character has been found
	ErrLetter = iota + 1
	// ErrWord returns when a sequence of characters is wrong
	ErrWord
	// ErrDecl returns when the unexpexted token has been found on the top level
	ErrDecl
	// ErrCurly returns when the unexpexted token, expecting {
	ErrCurly
	// ErrExp returns when the unexpected token, expecting expression or statement {
	ErrExp
	// ErrValue returns when the unexpected token, expecting value, identifier or calling func
	ErrValue
	// ErrRun returns when the compiler has found the second run function.
	ErrRun

	// ErrCompiler error. It means a bug.
	ErrCompiler
	// ErrRuntime error. It means bug
	ErrRuntime

	// ErrNoRun is returned when there is not run function
	ErrNoRun
)

var (
	errText = map[int]string{
		ErrLetter: `unknown character`,
		ErrWord:   `wrong sequence of characters`,
		ErrCurly:  `unexpected token, expecting {`,
		ErrDecl:   `expected declaration: func, run etc`,
		ErrExp:    `unexpected token, expecting expression or statement`,
		ErrRun:    `run function has already been defined`,
		ErrValue:  `unexpected token, expecting value, identifier or calling func`,

		ErrCompiler: `you have found a compiler bug. Let us know, please`,
		ErrRuntime:  `you have found a runtime bug. Let us know, please`,

		ErrNoRun: `there is not run function`,
	}
)

func errorText(id int) string {
	return errText[id]
}

func compileError(lp *Lex, idError, cur int, ext ...string) error {
	var more string
	line, column := lp.LineColumn(cur)
	if len(ext) > 0 {
		more = fmt.Sprintf(` (%s)`, ext[0])
	}
	return fmt.Errorf(` %d:%d: %s%s`, line, column, errorText(idError), more)
}

func runtimeError(rt *RunTime, idError int) error {
	var line, column int
	return fmt.Errorf(` %d:%d: %s`, line, column, errorText(idError))
}
