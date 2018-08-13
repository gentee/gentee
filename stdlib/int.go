// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"strconv"

	"github.com/gentee/gentee/core"
)

// InitInt appends stdlib int functions to the virtual machine
func InitInt(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		Add,      // binary +
		Div,      // binary /
		Equal,    // binary ==
		Greater,  // binary >
		Less,     // binary <
		Mod,      // binary %
		Mul,      // binary *
		Sign,     // unary sign -
		Sub,      // binary -
		BitOr,    // bitwise OR
		BitXor,   // bitwise XOR
		BitAnd,   // bitwise AND
		LShift,   // binary <<
		RShift,   // binary >>
		BitNot,   // unary bitwise NOT
		strºInt,  // str( int )
		boolºInt, // bool( int )
	} {
		vm.Units[core.DefName].NewEmbed(item)
	}
}

// Add add two integer value
func Add(left, right int64) int64 {
	return left + right
}

// BitAnd is bitwise AND
func BitAnd(left, right int64) int64 {
	return left & right
}

// BitNot is bitwise NOT
func BitNot(val int64) int64 {
	return ^val
}

// BitOr is bitwise OR
func BitOr(left, right int64) int64 {
	return left | right
}

// BitXor is bitwise XOR
func BitXor(left, right int64) int64 {
	return left ^ right
}

// Div divides one number by another
func Div(left, right int64) (int64, error) {
	if right == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	return left / right, nil
}

// Equal returns true if left == right
func Equal(left, right int64) bool {
	return left == right
}

// Greater returns true if left > right
func Greater(left, right int64) bool {
	return left > right
}

// Less returns true if left < right
func Less(left, right int64) bool {
	return left < right
}

// LShift returns left << right
func LShift(left, right int64) (int64, error) {
	if right < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	return left << uint64(right), nil
}

// Mod returns the remainder after division of one number by another
func Mod(left, right int64) int64 {
	return left % right
}

// Mul multiplies one number by another
func Mul(left, right int64) int64 {
	return left * right
}

// RShift returns left >> right
func RShift(left, right int64) (int64, error) {
	if right < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	return left >> uint64(right), nil
}

// Sign changes the sign of the integer value
func Sign(val int64) int64 {
	return -val
}

// Sub subtracts one number from another
func Sub(left, right int64) int64 {
	return left - right
}

// strºInt converts integer value to string
func strºInt(val int64) string {
	return strconv.FormatInt(val, 10)
}

// boolºInt converts integer value to boolean 0->false, not 0 -> true
func boolºInt(val int64) bool {
	return val != 0
}
