// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"strconv"

	"github.com/gentee/gentee/core"
)

// InitFloat appends stdlib float functions to the virtual machine
func InitFloat(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		AddºFloatFloat,       // binary +
		AddºFloatInt,         // binary +
		AddºIntFloat,         // binary +
		MulºFloatFloat,       // binary *
		MulºFloatInt,         // binary *
		MulºIntFloat,         // binary *
		SubºFloatFloat,       // binary -
		SubºFloatInt,         // binary -
		SubºIntFloat,         // binary -
		DivºFloatFloat,       // binary /
		DivºFloatInt,         // binary /
		DivºIntFloat,         // binary /
		EqualºFloatFloat,     // binary ==
		GreaterºFloatFloat,   // binary >
		LessºFloatFloat,      // binary <
		EqualºFloatInt,       // binary ==
		GreaterºFloatInt,     // binary >
		LessºFloatInt,        // binary <
		boolºFloat,           // bool( float )
		intºFloat,            // int( float )
		SignºFloat,           // unary sign -*/
		strºFloat,            // str( float )
		ExpStrºFloat,         // expression in string
		AssignºFloatFloat,    // float = float
		AssignAddºFloatFloat, // float += float
		AssignDivºFloatFloat, // float /= float
		AssignMulºFloatFloat, // float *= float
		AssignSubºFloatFloat, // float -= float
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// AddºFloatFloat adds two float values
func AddºFloatFloat(left, right float64) float64 {
	return left + right
}

// AddºFloatInt adds float and int
func AddºFloatInt(left float64, right int64) float64 {
	return left + float64(right)
}

// AddºIntFloat adds int and float
func AddºIntFloat(left int64, right float64) float64 {
	return float64(left) + right
}

// AssignºFloatFloat assign one float to another
func AssignºFloatFloat(ptr *interface{}, value float64) float64 {
	*ptr = value
	return (*ptr).(float64)
}

// AssignAddºFloatFloat adds one float to another
func AssignAddºFloatFloat(ptr *interface{}, value float64) float64 {
	*ptr = (*ptr).(float64) + value
	return (*ptr).(float64)
}

// AssignDivºFloatFloat does float /= float
func AssignDivºFloatFloat(ptr *interface{}, value float64) (float64, error) {
	if value == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	*ptr = (*ptr).(float64) / value
	return (*ptr).(float64), nil
}

// AssignMulºFloatFloat equals float *= float
func AssignMulºFloatFloat(ptr *interface{}, value float64) float64 {
	*ptr = (*ptr).(float64) * value
	return (*ptr).(float64)
}

// AssignSubºFloatFloat equals float *= float
func AssignSubºFloatFloat(ptr *interface{}, value float64) float64 {
	*ptr = (*ptr).(float64) - value
	return (*ptr).(float64)
}

// DivºFloatFloat divides one float by another
func DivºFloatFloat(left, right float64) (float64, error) {
	if right == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	return left / right, nil
}

// DivºFloatInt divides one float by int
func DivºFloatInt(left float64, right int64) (float64, error) {
	if right == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	return left / float64(right), nil
}

// DivºIntFloat divides one int by float
func DivºIntFloat(left int64, right float64) (float64, error) {
	if right == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	return float64(left) / right, nil
}

// ExpStrºFloat adds string and integer in string expression
func ExpStrºFloat(left string, right float64) string {
	return left + strºFloat(right)
}

// intºFloat converts float value to int
func intºFloat(val float64) int64 {
	return int64(val)
}

// MulºFloatFloat multiplies two float values
func MulºFloatFloat(left, right float64) float64 {
	return left * right
}

// MulºFloatInt multiplies float and int
func MulºFloatInt(left float64, right int64) float64 {
	return left * float64(right)
}

// MulºIntFloat multiplies int and float
func MulºIntFloat(left int64, right float64) float64 {
	return float64(left) * right
}

// SignºFloat changes the sign of the float value
func SignºFloat(val float64) float64 {
	return -val
}

// strºFloat converts float value to string
func strºFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

// SubºFloatFloat subtracts two float values
func SubºFloatFloat(left, right float64) float64 {
	return left - right
}

// SubºFloatInt subtracts float and int
func SubºFloatInt(left float64, right int64) float64 {
	return left - float64(right)
}

// SubºIntFloat subtracts int and float
func SubºIntFloat(left int64, right float64) float64 {
	return float64(left) - right
}

// EqualºFloatFloat returns true if left == right
func EqualºFloatFloat(left, right float64) bool {
	return left == right
}

// GreaterºFloatFloat returns true if left > right
func GreaterºFloatFloat(left, right float64) bool {
	return left > right
}

// LessºFloatFloat returns true if left < right
func LessºFloatFloat(left, right float64) bool {
	return left < right
}

// EqualºFloatInt returns true if left == right
func EqualºFloatInt(left float64, right int64) bool {
	return left == float64(right)
}

// GreaterºFloatInt returns true if left > right
func GreaterºFloatInt(left float64, right int64) bool {
	return left > float64(right)
}

// LessºFloatInt returns true if left < right
func LessºFloatInt(left float64, right int64) bool {
	return left < float64(right)
}

// boolºFloat converts integer value to boolean 0->false, not 0 -> true
func boolºFloat(val float64) bool {
	return val != 0.0
}
