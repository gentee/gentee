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
		Add,                 // binary +
		Div,                 // binary /
		Equal,               // binary ==
		Greater,             // binary >
		Less,                // binary <
		Mod,                 // binary %
		Mul,                 // binary *
		Sign,                // unary sign -
		Sub,                 // binary -
		BitOr,               // bitwise OR
		BitXor,              // bitwise XOR
		BitAnd,              // bitwise AND
		LShift,              // binary <<
		RShift,              // binary >>
		BitNot,              // unary bitwise NOT
		strºInt,             // str( int )
		boolºInt,            // bool( int )
		ExpStrºInt,          // expression in string
		AssignºIntInt,       // int = int
		AssignAddºIntInt,    // int += int
		AssignBitAndºIntInt, // int &= int
		AssignBitOrºIntInt,  // int |= int
		AssignBitXorºIntInt, // int ^= int
		AssignDivºIntInt,    // int /= int
		AssignModºIntInt,    // int %= int
		AssignMulºIntInt,    // int *= int
		AssignSubºIntInt,    // int -= int
		AssignLShiftºIntInt, // int <<= int
		AssignRShiftºIntInt, // int >>= int
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// AssignºIntInt assign one integer to another
func AssignºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) int64 {
	vars[cmdVar.Index] = value
	return vars[cmdVar.Index].(int64)
}

// AssignAddºIntInt adds one integer to another
func AssignAddºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) int64 {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) + value
	return vars[cmdVar.Index].(int64)
}

// AssignBitAndºIntInt equals int &= int
func AssignBitAndºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) int64 {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) & value
	return vars[cmdVar.Index].(int64)
}

// AssignBitOrºIntInt equals int |= int
func AssignBitOrºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) int64 {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) | value
	return vars[cmdVar.Index].(int64)
}

// AssignBitXorºIntInt equals int ^= int
func AssignBitXorºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) int64 {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) ^ value
	return vars[cmdVar.Index].(int64)
}

// AssignDivºIntInt does int /= int
func AssignDivºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) (int64, error) {
	if value == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) / value
	return vars[cmdVar.Index].(int64), nil
}

// AssignModºIntInt equals int %= int
func AssignModºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) (int64, error) {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) % value
	return vars[cmdVar.Index].(int64), nil
}

// AssignMulºIntInt equals int *= int
func AssignMulºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) (int64, error) {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) * value
	return vars[cmdVar.Index].(int64), nil
}

// AssignSubºIntInt equals int *= int
func AssignSubºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) (int64, error) {
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) - value
	return vars[cmdVar.Index].(int64), nil
}

// AssignLShiftºIntInt does int <<= int
func AssignLShiftºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) (int64, error) {
	if value < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) << uint64(value)
	return vars[cmdVar.Index].(int64), nil
}

// AssignRShiftºIntInt does int >>= int
func AssignRShiftºIntInt(vars []interface{}, cmdVar *core.CmdVar, value int64) (int64, error) {
	if value < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	vars[cmdVar.Index] = vars[cmdVar.Index].(int64) >> uint64(value)
	return vars[cmdVar.Index].(int64), nil
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

// ExpStrºInt adds string and integer in string expression
func ExpStrºInt(left string, right int64) string {
	return left + strºInt(right)
}
