// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gentee/gentee/core"
)

// InitFloat appends stdlib float functions to the virtual machine
func InitFloat(ws *core.Workspace) {
	for _, item := range []interface{}{
		core.Link{AddºFloatFloat, core.ADDFLOAT},         // binary +
		core.Link{AddºFloatInt, 21<<16 | core.EMBED},     // binary +
		core.Link{AddºIntFloat, 22<<16 | core.EMBED},     // binary +
		core.Link{MulºFloatFloat, core.MULFLOAT},         // binary *
		core.Link{MulºFloatInt, 25<<16 | core.EMBED},     // binary *
		core.Link{MulºIntFloat, 30<<16 | core.EMBED},     // binary *
		core.Link{SubºFloatFloat, core.SUBFLOAT},         // binary -
		core.Link{SubºFloatInt, 26<<16 | core.EMBED},     // binary -
		core.Link{SubºIntFloat, 27<<16 | core.EMBED},     // binary -
		core.Link{DivºFloatFloat, core.DIVFLOAT},         // binary /
		core.Link{DivºFloatInt, 28<<16 | core.EMBED},     // binary /
		core.Link{DivºIntFloat, 29<<16 | core.EMBED},     // binary /
		core.Link{EqualºFloatFloat, core.EQFLOAT},        // binary ==
		core.Link{GreaterºFloatFloat, core.GTFLOAT},      // binary >
		core.Link{LessºFloatFloat, core.LTFLOAT},         // binary <
		core.Link{EqualºFloatInt, 32<<16 | core.EMBED},   // binary ==
		core.Link{GreaterºFloatInt, 33<<16 | core.EMBED}, // binary >
		core.Link{LessºFloatInt, 34<<16 | core.EMBED},    // binary <
		boolºFloat, // bool( float )
		core.Link{intºFloat, 23<<16 | core.EMBED},        // int( float )
		core.Link{SignºFloat, core.SIGNFLOAT},            // unary sign -*/
		core.Link{strºFloat, 24<<16 | core.EMBED},        // str( float )
		core.Link{ExpStrºFloat, 31<<16 | core.EMBED},     // expression in string
		core.Link{AssignºFloatFloat, core.ASSIGN},        // float = float
		core.Link{AssignAddºFloatFloat, core.ASSIGN + 1}, // float += float
		core.Link{AssignDivºFloatFloat, core.ASSIGN + 4}, // float /= float
		core.Link{AssignMulºFloatFloat, core.ASSIGN + 3}, // float *= float
		core.Link{AssignSubºFloatFloat, core.ASSIGN + 2}, // float -= float
		RoundºFloat,    // Round( float ) int
		FloorºFloat,    // Floor( float ) int
		CeilºFloat,     // Ceil( float ) int
		RoundºFloatInt, // Round( float, int ) float
		MinºFloatFloat, // Min(float, float)
		MaxºFloatFloat, // Max(float, float)
	} {
		ws.StdLib().NewEmbed(item)
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

// ExpStrºFloat adds string and float in string expression
func ExpStrºFloat(left string, right float64) string {
	return left + strºFloat(right)
}

// intºFloat converts float value to int
func intºFloat(val float64) int64 {
	return int64(val)
}

// MaxºFloatFloat returns the maximum of two float numbers
func MaxºFloatFloat(left, right float64) float64 {
	if left < right {
		return right
	}
	return left
}

// MinºFloatFloat returns the minimum of two float numbers
func MinºFloatFloat(left, right float64) float64 {
	if left > right {
		return right
	}
	return left
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

// RoundºFloat returns the nearest integer, rounding half away from zero.
func RoundºFloat(val float64) int64 {
	return int64(math.Round(val))
}

// FloorºFloat returns the greatest integer value less than or equal to val.
func FloorºFloat(val float64) int64 {
	return int64(math.Floor(val))
}

// CeilºFloat returns the least integer value greater than or equal to val.
func CeilºFloat(val float64) int64 {
	return int64(math.Ceil(val))
}

// RoundºFloatInt returns a number with the specified number of decimal places.
func RoundºFloatInt(val float64, digits int64) float64 {
	dec := float64(1)
	for ; digits > 0; digits-- {
		dec *= 10
	}
	return math.Round(val*dec) / dec
}
