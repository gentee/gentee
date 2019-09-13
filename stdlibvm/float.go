// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gentee/gentee/core"
)

// boolºFloat converts integer value to boolean 0->false, not 0 -> true
func boolºFloat(val float64) int64 {
	if val != 0.0 {
		return 1
	}
	return 0
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
func AssignºFloatFloat(ptr *float64, value float64) (float64, error) {
	*ptr = value
	return *ptr, nil
}

// AssignAddºFloatFloat adds one float to another
func AssignAddºFloatFloat(ptr *float64, value float64) (float64, error) {
	*ptr += value
	return *ptr, nil
}

// AssignDivºFloatFloat does float /= float
func AssignDivºFloatFloat(ptr *float64, value float64) (float64, error) {
	if value == 0 {
		return 0, fmt.Errorf(core.ErrorText(core.ErrDivZero))
	}
	*ptr /= value
	return *ptr, nil
}

// AssignMulºFloatFloat equals float *= float
func AssignMulºFloatFloat(ptr *float64, value float64) (float64, error) {
	*ptr *= value
	return *ptr, nil
}

// AssignSubºFloatFloat equals float *= float
func AssignSubºFloatFloat(ptr *float64, value float64) (float64, error) {
	*ptr -= value
	return *ptr, nil
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
	return left + StrºFloat(right)
}

// MulºFloatInt multiplies float and int
func MulºFloatInt(left float64, right int64) float64 {
	return left * float64(right)
}

// MulºIntFloat multiplies int and float
func MulºIntFloat(left int64, right float64) float64 {
	return float64(left) * right
}

// SubºFloatInt subtracts float and int
func SubºFloatInt(left float64, right int64) float64 {
	return left - float64(right)
}

// SubºIntFloat subtracts int and float
func SubºIntFloat(left int64, right float64) float64 {
	return float64(left) - right
}

// EqualºFloatInt returns true if left == right
func EqualºFloatInt(left float64, right int64) int64 {
	if left == float64(right) {
		return 1
	}
	return 0
}

// GreaterºFloatInt returns true if left > right
func GreaterºFloatInt(left float64, right int64) int64 {
	if left > float64(right) {
		return 1
	}
	return 0
}

// LessºFloatInt returns true if left < right
func LessºFloatInt(left float64, right int64) int64 {
	if left < float64(right) {
		return 1
	}
	return 0
}

// intºFloat converts float value to int
func intºFloat(val float64) int64 {
	return int64(val)
}

// StrºFloat converts float value to string
func StrºFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

// CeilºFloat returns the least integer value greater than or equal to val.
func CeilºFloat(val float64) int64 {
	return int64(math.Ceil(val))
}

// FloorºFloat returns the greatest integer value less than or equal to val.
func FloorºFloat(val float64) int64 {
	return int64(math.Floor(val))
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

// RoundºFloat returns the nearest integer, rounding half away from zero.
func RoundºFloat(val float64) int64 {
	return int64(math.Round(val))
}

// RoundºFloatInt returns a number with the specified number of decimal places.
func RoundºFloatInt(val float64, digits int64) float64 {
	dec := float64(1)
	for ; digits > 0; digits-- {
		dec *= 10
	}
	return math.Round(val*dec) / dec
}
