// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitChar appends stdlib int functions to the virtual machine
func InitChar(ws *core.Workspace) {
	for _, item := range []interface{}{
		core.Link{AddºCharChar, 12<<16 | core.EMBED},     // binary +
		core.Link{AddºStrChar, 14<<16 | core.EMBED},      // binary str + char
		core.Link{AddºCharStr, 13<<16 | core.EMBED},      // binary char + str
		core.Link{AssignAddºStrChar, core.ASSIGN + 4},    // str += char
		core.Link{AssignºCharChar, core.ASSIGN},          // char = char
		core.Link{ExpStrºChar, 15<<16 | core.EMBED},      // expression in string
		core.Link{EqualºCharChar, core.EQ},               // binary ==
		core.Link{GreaterºCharChar, 16<<16 | core.EMBED}, // binary >
		core.Link{LessºCharChar, 17<<16 | core.EMBED},    // binary <
		core.Link{intºChar, core.NOP},                    // int( char )
		core.Link{strºChar, 18<<16 | core.EMBED},         // str( char )
	} {
		ws.StdLib().NewEmbed(item)
	}
}

// AddºCharChar adds two rune values
func AddºCharChar(left, right rune) string {
	return string(left) + string(right)
}

// AddºStrChar appends rune to string
func AddºStrChar(left string, right rune) string {
	return left + string(right)
}

// AddºCharStr appends rune to string
func AddºCharStr(left rune, right string) string {
	return string(left) + right
}

// AssignºCharChar assigns one rune to another
func AssignºCharChar(ptr *interface{}, value rune) rune {
	*ptr = value
	return (*ptr).(rune)
}

// AssignAddºStrChar appends one rune to string
func AssignAddºStrChar(ptr *interface{}, value rune) string {
	*ptr = (*ptr).(string) + string(value)
	return (*ptr).(string)
}

// ExpStrºChar adds string and char in string expression
func ExpStrºChar(left string, right rune) string {
	return left + string(right)
}

// EqualºCharChar returns true if left == right
func EqualºCharChar(left, right rune) bool {
	return left == right
}

// GreaterºCharChar returns true if left > right
func GreaterºCharChar(left, right rune) bool {
	return left > right
}

// LessºCharChar returns true if left < right
func LessºCharChar(left, right rune) bool {
	return left < right
}

// intºChar converts char value to int64
func intºChar(val rune) int64 {
	return int64(val)
}

// strºChar converts char value to string
func strºChar(val rune) string {
	return string(val)
}
