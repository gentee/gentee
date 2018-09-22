// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"fmt"
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
	// ErrEmptyCommand is returned if empty command is specified in $
	ErrEmptyCommand
	// ErrQuoteCommand is returned if there is an unclosed quotation mark in $ command
	ErrQuoteCommand
	// ErrIndexOut means that int index is out of the length of the array
	ErrIndexOut

	// ErrRuntime error. It means bug
	ErrRuntime
)

var (
	errText = map[int]string{
		ErrRunIndex:     `invalid name of Run`,
		ErrDepth:        `maximum depth of recursion has been reached`,
		ErrDivZero:      `divided by zero`,
		ErrCycle:        `maximum cycle count has been reached`,
		ErrShift:        `right operand of shift cannot be less than zero`,
		ErrStrToInt:     `converting string to integer is invalid`,
		ErrEmptyCommand: `empty $ command`,
		ErrQuoteCommand: `unclosed quotation mark in $ command`,
		ErrIndexOut:     `index out of range`,

		ErrRuntime: `you have found a runtime bug. Let us know, please`,
	}
)

// ErrorText returns the text of the error message
func ErrorText(id int) string {
	return errText[id]
}

func runtimeError(rt *RunTime, cmd ICmd, idError interface{}, labels ...interface{}) error {
	var (
		line, column int
		lex          *Lex
		errText      string
	)
	for i := len(rt.Calls) - 1; i >= 0 && lex == nil; i-- {
		if cmd != nil && rt.Calls[i].GetObject() != nil {
			if lex = rt.Calls[i].GetObject().GetLex(); lex != nil {
				line, column = lex.LineColumn(cmd.GetToken())
				break
			}
		}
	}
	switch v := idError.(type) {
	case int:
		errText = ErrorText(v)
	case error:
		errText = v.Error()
	}
	for _, item := range labels {
		errText += fmt.Sprintf(` [%v]`, item)
	}
	return fmt.Errorf(`%d:%d: %s`, line, column, errText)
}
