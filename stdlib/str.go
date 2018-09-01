// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// InitStr appends stdlib int functions to the virtual machine
func InitStr(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		AddºStr,          // binary +
		EqualºStr,        // binary ==
		GreaterºStr,      // binary >
		LenºStr,          // the length of str
		LessºStr,         // binary <
		intºStr,          // int( str )
		boolºStr,         // bool( str )
		ExpStr,           // expression in string
		AssignºStrStr,    // str = str
		AssignAddºStrStr, // str += str
		AssignºStrBool,   // str = bool
		AssignºStrInt,    // str = int
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// AssignºStrStr assigns one string to another
func AssignºStrStr(vars []interface{}, cmdVar *core.CmdVar, value string) string {
	vars[cmdVar.Index] = value
	return vars[cmdVar.Index].(string)
}

// AssignAddºStrStr appends one string to another
func AssignAddºStrStr(vars []interface{}, cmdVar *core.CmdVar, value string) string {
	vars[cmdVar.Index] = vars[cmdVar.Index].(string) + value
	return vars[cmdVar.Index].(string)
}

// AssignºStrBool assigns boolean to string
func AssignºStrBool(vars []interface{}, cmdVar *core.CmdVar, value bool) string {
	vars[cmdVar.Index] = fmt.Sprint(value)
	return vars[cmdVar.Index].(string)
}

// AssignºStrInt assigns integer to string
func AssignºStrInt(vars []interface{}, cmdVar *core.CmdVar, value int64) string {
	vars[cmdVar.Index] = fmt.Sprint(value)
	return vars[cmdVar.Index].(string)
}

// ExpStr adds two strings in string expression
func ExpStr(left, right string) string {
	return left + right
}

// AddºStr adds two integer value
func AddºStr(left, right string) string {
	return left + right
}

// EqualºStr returns true if left == right
func EqualºStr(left, right string) bool {
	return left == right
}

// GreaterºStr returns true if left > right
func GreaterºStr(left, right string) bool {
	return left > right
}

// LenºStr returns teh length of the string
func LenºStr(param string) int64 {
	return int64(len(param))
}

// LessºStr returns true if left < right
func LessºStr(left, right string) bool {
	return left < right
}

// intºStr converts strings value to int64
func intºStr(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(val, 0, 64)
	if err != nil {
		err = errors.New(core.ErrorText(core.ErrStrToInt))
	}
	return
}

// intºBool converts boolean value to int false -> 0, true -> 1
func boolºStr(val string) bool {
	return len(val) != 0 && val != `0` && strings.ToLower(val) != `false`
}
