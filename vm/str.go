// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// AssignºStrBool assigns boolean to string
func AssignºStrBool(ptr *string, value interface{}) (string, error) {
	*ptr = strºBool(value.(int64))
	return *ptr, nil
}

// AssignºStrInt assigns integer to string
func AssignºStrInt(ptr *string, value interface{}) (string, error) {
	*ptr = fmt.Sprint(value)
	return *ptr, nil
}

// AssignAddºStrStr appends one string to another
func AssignAddºStrStr(ptr *string, value interface{}) (string, error) {
	*ptr += value.(string)
	return *ptr, nil
}

// boolºStr converts string value to bool
func boolºStr(val string) int64 {
	if len(val) != 0 && val != `0` && strings.ToLower(val) != `false` {
		return 1
	}
	return 0
}

// FindºStrStr returns the index of the first instance of substr
func FindºStrStr(s, substr string) (off int64) {
	off = int64(strings.Index(s, substr))
	if off > 0 {
		off = int64(len([]rune(s[:off])))
	}
	return
}

// floatºStr converts string value to float
func floatºStr(val string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(val, 64)
	if err != nil {
		err = errors.New(ErrorText(ErrStrToFloat))
	}
	return
}

// FormatºStr formats according to a format specifier and returns the resulting string
func FormatºStr(pattern string, pars ...interface{}) string {
	return fmt.Sprintf(pattern, pars...)
}

// HasPrefixºStrStr returns true if the string s begins with prefix
func HasPrefixºStrStr(s, prefix string) int64 {
	if strings.HasPrefix(s, prefix) {
		return 1
	}
	return 0
}

// HasSuffixºStrStr returns true if the string s ends with suffix
func HasSuffixºStrStr(s, suffix string) int64 {
	if strings.HasSuffix(s, suffix) {
		return 1
	}
	return 0
}

// intºStr converts string value to int
func intºStr(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(val, 0, 64)
	if err != nil {
		err = errors.New(ErrorText(ErrStrToInt))
	}
	return
}

// LeftºStrInt cuts the string.
func LeftºStrInt(s string, count int64) string {
	r := []rune(s)
	if int(count) > len(r) {
		count = int64(len(r))
	}
	return string(r[:count])
}

// LinesºStr splits a string to a array of strings
func LinesºStr(in string) *core.Array {
	out := core.NewArray()
	items := strings.Split(in, "\n")
	for _, item := range items {
		out.Data = append(out.Data, strings.Trim(item, "\r"))
	}
	return out
}

// LowerºStr converts a copy of the string to their lower case and returns it.
func LowerºStr(s string) string {
	return strings.ToLower(s)
}

// RepeatºStrInt returns a new string consisting of count copies of the specified string.
func RepeatºStrInt(input string, count int64) string {
	return strings.Repeat(input, int(count))
}

// ReplaceºStrStrStr replaces strings in a string
func ReplaceºStrStrStr(in, old, new string) string {
	return strings.Replace(in, old, new, -1)
}

// RightºStrInt returns the right substring of the string.
func RightºStrInt(s string, count int64) string {
	r := []rune(s)
	off := len(r) - int(count)
	if off < 0 {
		off = 0
	}
	return string(r[off:])
}

// ShiftºStr trims white spaces characters in the each line of the string.
func ShiftºStr(par string) string {
	lines := strings.Split(par, "\n")
	for i, v := range lines {
		lines[i] = strings.TrimSpace(v)
	}
	return strings.Join(lines, "\n")
}

// SplitºStrStr splits a string to a array of strings
func SplitºStrStr(in, sep string) *core.Array {
	out := core.NewArray()
	items := strings.Split(in, sep)
	for _, item := range items {
		out.Data = append(out.Data, item)
	}
	return out
}

// SubstrºStrIntInt returns a substring with the specified offset and length
func SubstrºStrIntInt(in string, off, length int64) (string, error) {
	var rin []rune
	rin = []rune(in)
	rlen := int64(len(rin))
	if length < 0 {
		length = -length
		off -= length
	}
	if off < 0 || off >= rlen || off+length > rlen {
		return ``, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	if length == 0 {
		length = rlen - off
	}
	return string(rin[off : off+length]), nil
}

// TrimºStr returns a substring of the string with all leading and trailing characters in set removed.
func TrimºStr(in string, set string) string {
	return strings.Trim(in, set)
}

// TrimLeftºStr returns a substring of the string with all leading characters in set removed.
func TrimLeftºStr(in string, set string) string {
	return strings.TrimLeft(in, set)
}

// TrimRightºStr returns a substring of the string with all trailing characters in set removed.
func TrimRightºStr(in string, set string) string {
	return strings.TrimRight(in, set)
}

// TrimSpaceºStr trims white space in a string
func TrimSpaceºStr(in string) string {
	return strings.TrimSpace(in)
}

// UpperºStr converts a copy of the string to their upper case and returns it.
func UpperºStr(s string) string {
	return strings.ToUpper(s)
}
