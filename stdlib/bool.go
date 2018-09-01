// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitBool appends stdlib bool functions to the virtual machine
func InitBool(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		strºBool,        // str( bool )
		intºBool,        // int( bool )
		Not,             // unary boolean not
		ExpStrºBool,     // expression in string
		AssignºBoolBool, // bool = bool
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// AssignºBoolBool assign one boolean to another
func AssignºBoolBool(vars []interface{}, cmdVar *core.CmdVar, value bool) bool {
	vars[cmdVar.Index] = value
	return vars[cmdVar.Index].(bool)
}

// Not changes true to false or false to true
func Not(val bool) bool {
	return !val
}

// strºBool converts boolean value to string
func strºBool(val bool) string {
	if val {
		return `true`
	}
	return `false`
}

// intºBool converts boolean value to int false -> 0, true -> 1
func intºBool(val bool) int64 {
	if val {
		return 1
	}
	return 0
}

// ExpStrºBool adds string and boolean in string expression
func ExpStrºBool(left string, right bool) string {
	return left + strºBool(right)
}
