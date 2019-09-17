// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// floatºStr converts string value to float
func floatºStr(val string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(val, 64)
	if err != nil {
		err = errors.New(core.ErrorText(core.ErrStrToFloat))
	}
	return
}

// intºStr converts string value to int
func intºStr(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(val, 0, 64)
	if err != nil {
		err = errors.New(core.ErrorText(core.ErrStrToInt))
	}
	return
}

// AssignºStrBool assigns boolean to string
func AssignºStrBool(ptr *string, value interface{}) (string, error) {
	*ptr = StrºBool(value.(int64))
	return *ptr, nil
}

// AssignºStrInt assigns integer to string
func AssignºStrInt(ptr *string, value interface{}) (string, error) {
	*ptr = fmt.Sprint(value)
	return *ptr, nil
}

// AssignºStrStr assigns one string to another
func AssignºStrStr(ptr *string, value interface{}) (string, error) {
	*ptr = value.(string)
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

// LinesºStr splits a string to a array of strings
func LinesºStr(in string) *core.Array {
	out := core.NewArray()
	items := strings.Split(in, "\n")
	for _, item := range items {
		out.Data = append(out.Data, strings.Trim(item, "\r"))
	}
	return out
}

// LeftºStrInt cuts the string.
func LeftºStrInt(s string, count int64) string {
	r := []rune(s)
	if int(count) > len(r) {
		count = int64(len(r))
	}
	return string(r[:count])
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

func replaceArr(in string, old, new []string) string {
	input := []rune(in)
	out := make([]rune, 0, len(input))
	lin := len(input)
	for i := 0; i < lin; i++ {
		eq := -1
		maxLen := lin - i
		for k, item := range old {
			litem := len([]rune(item))
			if maxLen >= litem && string(input[i:i+litem]) == item {
				eq = k
				break
			}
		}
		if eq >= 0 {
			out = append(out, []rune(new[eq])...)
			i += len([]rune(old[eq])) - 1
		} else {
			out = append(out, input[i])
		}
	}
	return string(out)
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
		return ``, fmt.Errorf(core.ErrorText(core.ErrInvalidParam))
	}
	if length == 0 {
		length = rlen - off
	}
	return string(rin[off : off+length]), nil
}

// TrimSpaceºStr trims white space in a string
func TrimSpaceºStr(in string) string {
	return strings.TrimSpace(in)
}

func ReplaceArr(in string, old, new []string) string {
	input := []rune(in)
	out := make([]rune, 0, len(input))
	lin := len(input)
	for i := 0; i < lin; i++ {
		eq := -1
		maxLen := lin - i
		for k, item := range old {
			litem := len([]rune(item))
			if maxLen >= litem && string(input[i:i+litem]) == item {
				eq = k
				break
			}
		}
		if eq >= 0 {
			out = append(out, []rune(new[eq])...)
			i += len([]rune(old[eq])) - 1
		} else {
			out = append(out, input[i])
		}
	}
	return string(out)
}

// TrimRightºStr trims white space in a string
func TrimRightºStr(in string, set string) string {
	return strings.TrimRight(in, set)
}

// UpperºStr converts a copy of the string to their upper case and returns it.
func UpperºStr(s string) string {
	return strings.ToUpper(s)
}
