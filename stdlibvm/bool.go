// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

// StrºBool converts boolean value to string
func StrºBool(val int64) string {
	if val != 0 {
		return `true`
	}
	return `false`
}

// ExpStrºBool adds string and boolean in string expression
func ExpStrºBool(left string, right int64) string {
	return left + StrºBool(right)
}
