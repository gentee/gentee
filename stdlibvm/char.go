// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

// AddºCharChar adds two rune values
func AddºCharChar(left, right int64) string {
	return string(rune(left)) + string(rune(right))
}

// AddºCharStr appends rune to string
func AddºCharStr(left int64, right string) string {
	return string(rune(left)) + right
}

// AddºStrChar appends rune to string
func AddºStrChar(left string, right int64) string {
	return left + string(rune(right))
}

// AssignAddºStrChar appends one rune to string
func AssignAddºStrChar(ptr *string, value int64) string {
	*ptr += string(rune(value))
	return *ptr
}

// ExpStrºChar adds string and char in string expression
func ExpStrºChar(left string, right int64) string {
	return left + string(rune(right))
}

// GreaterºCharChar returns true if left > right
func GreaterºCharChar(left, right int64) int64 {
	if rune(left) > rune(right) {
		return 1
	}
	return 0
}

// LessºCharChar returns true if left < right
func LessºCharChar(left, right int64) int64 {
	if rune(left) < rune(right) {
		return 1
	}
	return 0
}

// strºChar converts char value to string
func strºChar(val int64) string {
	return string(rune(val))
}
