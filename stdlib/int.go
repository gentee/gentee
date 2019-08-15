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
func InitInt(ws *core.Workspace) {
	for _, item := range []interface{}{
		core.Link{AddºIntInt, core.ADD},                   // binary +
		core.Link{DivºIntInt, core.DIV},                   // binary /
		core.Link{EqualºIntInt, core.EQ},                  // binary ==
		core.Link{GreaterºIntInt, core.GT},                // binary >
		core.Link{LessºIntInt, core.LT},                   // binary <
		core.Link{ModºIntInt, core.MOD},                   // binary %
		core.Link{MulºIntInt, core.MUL},                   // binary *
		core.Link{SignºInt, core.SIGN},                    // unary sign -
		core.Link{SubºIntInt, core.SUB},                   // binary -
		core.Link{BitOrºIntInt, core.BITOR},               // bitwise OR
		core.Link{BitXorºIntInt, core.BITXOR},             // bitwise XOR
		core.Link{BitAndºIntInt, core.BITAND},             // bitwise AND
		core.Link{LShiftºIntInt, core.LSHIFT},             // binary <<
		core.Link{RShiftºIntInt, core.RSHIFT},             // binary >>
		core.Link{BitNotºInt, core.BITNOT},                // unary bitwise NOT
		floatºInt,                                         // float( int )
		core.Link{strºInt, 1<<16 | core.EMBED},            // str( int )
		core.Link{boolºInt, 3<<16 | core.EMBED},           // bool( int )
		core.Link{ExpStrºInt, 4<<16 | core.EMBED},         // expression in string
		core.Link{AssignºIntInt, core.ASSIGN},             // int = int
		AssignºIntChar,                                    // int = char
		core.Link{AssignAddºIntInt, core.ASSIGNADD},       // int += int
		core.Link{AssignBitAndºIntInt, core.ASSIGNBITAND}, // int &= int
		core.Link{AssignBitOrºIntInt, core.ASSIGNBITOR},   // int |= int
		core.Link{AssignBitXorºIntInt, core.ASSIGNBITXOR}, // int ^= int
		core.Link{AssignDivºIntInt, core.ASSIGNDIV},       // int /= int
		core.Link{AssignModºIntInt, core.ASSIGNMOD},       // int %= int
		core.Link{AssignMulºIntInt, core.ASSIGNMUL},       // int *= int
		core.Link{AssignSubºIntInt, core.ASSIGNSUB},       // int -= int
		core.Link{AssignLShiftºIntInt, core.ASSIGNLSHIFT}, // int <<= int
		core.Link{AssignRShiftºIntInt, core.ASSIGNRSHIFT}, // int >>= int
		MaxºIntInt, // Max(int, int)
		MinºIntInt, // Min(int, int)
	} {
		ws.StdLib().NewEmbed(item)
	}
}

// AssignºIntInt assign one integer to another
func AssignºIntInt(ptr *interface{}, value int64) int64 {
	*ptr = value
	return (*ptr).(int64)
}

// AssignºIntChar assign a rune to integer
func AssignºIntChar(ptr *interface{}, value rune) int64 {
	*ptr = int64(value)
	return (*ptr).(int64)
}

// AssignAddºIntInt adds one integer to another
func AssignAddºIntInt(ptr *interface{}, value int64) (int64, error) {
	switch v := (*ptr).(type) {
	case uint8:
		value += int64(v)
		if uint64(value) > 255 {
			return 0, fmt.Errorf(core.ErrorText(core.ErrByteOut))
		}
		*ptr = value
	default:
		*ptr = v.(int64) + value
	}
	return (*ptr).(int64), nil
}

// AssignBitAndºIntInt equals int &= int
func AssignBitAndºIntInt(ptr *interface{}, value int64) int64 {
	*ptr = (*ptr).(int64) & value
	return (*ptr).(int64)
}

// AssignBitOrºIntInt equals int |= int
func AssignBitOrºIntInt(ptr *interface{}, value int64) int64 {
	*ptr = (*ptr).(int64) | value
	return (*ptr).(int64)
}

// AssignBitXorºIntInt equals int ^= int
func AssignBitXorºIntInt(ptr *interface{}, value int64) int64 {
	*ptr = (*ptr).(int64) ^ value
	return (*ptr).(int64)
}

// AssignDivºIntInt does int /= int
func AssignDivºIntInt(ptr *interface{}, value int64) (int64, error) {
	if value == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	*ptr = (*ptr).(int64) / value
	return (*ptr).(int64), nil
}

// AssignModºIntInt equals int %= int
func AssignModºIntInt(ptr *interface{}, value int64) (int64, error) {
	if value == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	*ptr = (*ptr).(int64) % value
	return (*ptr).(int64), nil
}

// AssignMulºIntInt equals int *= int
func AssignMulºIntInt(ptr *interface{}, value int64) int64 {
	*ptr = (*ptr).(int64) * value
	return (*ptr).(int64)
}

// AssignSubºIntInt equals int *= int
func AssignSubºIntInt(ptr *interface{}, value int64) int64 {
	*ptr = (*ptr).(int64) - value
	return (*ptr).(int64)
}

// AssignLShiftºIntInt does int <<= int
func AssignLShiftºIntInt(ptr *interface{}, value int64) (int64, error) {
	if value < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	*ptr = (*ptr).(int64) << uint64(value)
	return (*ptr).(int64), nil
}

// AssignRShiftºIntInt does int >>= int
func AssignRShiftºIntInt(ptr *interface{}, value int64) (int64, error) {
	if value < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	*ptr = (*ptr).(int64) >> uint64(value)
	return (*ptr).(int64), nil
}

// AddºIntInt add two integer value
func AddºIntInt(left, right int64) int64 {
	return left + right
}

// BitAndºIntInt is bitwise AND
func BitAndºIntInt(left, right int64) int64 {
	return left & right
}

// BitNotºInt is bitwise NOT
func BitNotºInt(val int64) int64 {
	return ^val
}

// BitOrºIntInt is bitwise OR
func BitOrºIntInt(left, right int64) int64 {
	return left | right
}

// BitXorºIntInt is bitwise XOR
func BitXorºIntInt(left, right int64) int64 {
	return left ^ right
}

// DivºIntInt divides one number by another
func DivºIntInt(left, right int64) (int64, error) {
	if right == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	return left / right, nil
}

// EqualºIntInt returns true if left == right
func EqualºIntInt(left, right int64) bool {
	return left == right
}

// GreaterºIntInt returns true if left > right
func GreaterºIntInt(left, right int64) bool {
	return left > right
}

// LessºIntInt returns true if left < right
func LessºIntInt(left, right int64) bool {
	return left < right
}

// LShiftºIntInt returns left << right
func LShiftºIntInt(left, right int64) (int64, error) {
	if right < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	return left << uint64(right), nil
}

// MaxºIntInt returns the maximum of two integers
func MaxºIntInt(left, right int64) int64 {
	if left < right {
		return right
	}
	return left
}

// MinºIntInt returns the minimum of two integers
func MinºIntInt(left, right int64) int64 {
	if left > right {
		return right
	}
	return left
}

// ModºIntInt returns the remainder after division of one number by another
func ModºIntInt(left, right int64) (int64, error) {
	if right == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	return left % right, nil
}

// MulºIntInt multiplies one number by another
func MulºIntInt(left, right int64) int64 {
	return left * right
}

// RShiftºIntInt returns left >> right
func RShiftºIntInt(left, right int64) (int64, error) {
	if right < 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrShift))
	}
	return left >> uint64(right), nil
}

// SignºInt changes the sign of the integer value
func SignºInt(val int64) int64 {
	return -val
}

// SubºIntInt subtracts one number from another
func SubºIntInt(left, right int64) int64 {
	return left - right
}

// floatºInt converts integer value to float
func floatºInt(val int64) float64 {
	return float64(val)
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
