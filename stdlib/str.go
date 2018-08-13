// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"errors"
	"strconv"
	"strings"

	"bitbucket.org/novostrim/go-gentee/core"
)

// InitStr appends stdlib int functions to the virtual machine
func InitStr(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		AddºStr,     // binary +
		EqualºStr,   // binary ==
		GreaterºStr, // binary >
		LessºStr,    // binary <
		intºStr,     // int( str )
		boolºStr,    // bool( str )
	} {
		vm.Units[core.DefName].NewEmbed(item)
	}
}

// AddºStr add two integer value
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
