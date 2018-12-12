// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// InitStr appends stdlib string functions to the virtual machine
func InitStr(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		AddºStrStr,        // binary +
		EqualºStrStr,      // binary ==
		GreaterºStrStr,    // binary >
		LenºStr,           // the length of str
		LessºStrStr,       // binary <
		intºStr,           // int( str )
		floatºStr,         // float( str )
		boolºStr,          // bool( str )
		ExpStrºStr,        // expression in string
		AssignºStrStr,     // str = str
		AssignAddºStrStr,  // str += str
		AssignºStrBool,    // str = bool
		AssignºStrInt,     // str = int
		FindºStrStr,       // Find( str, str ) int
		FormatºStr,        // Format( str, ... ) str
		HasPrefixºStrStr,  // HasPrefix( str, str ) bool
		HasSuffixºStrStr,  // HasSuffix( str, str ) bool
		ReplaceºStrStrStr, // Replace( str, str, str )
		LinesºStrArr,      // Lines( str, arr )
		SplitºStrStrArr,   // Split( str, str, arr )
		SubstrºStrIntInt,  // Substr( str, int, int ) str
		TrimSpaceºStr,     // TrimSpace( str )
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// AssignºStrStr assigns one string to another
func AssignºStrStr(ptr *interface{}, value string) string {
	*ptr = value
	return (*ptr).(string)
}

// AssignAddºStrStr appends one string to another
func AssignAddºStrStr(ptr *interface{}, value string) string {
	*ptr = (*ptr).(string) + value
	return (*ptr).(string)
}

// AssignºStrBool assigns boolean to string
func AssignºStrBool(ptr *interface{}, value bool) string {
	*ptr = fmt.Sprint(value)
	return (*ptr).(string)
}

// AssignºStrInt assigns integer to string
func AssignºStrInt(ptr *interface{}, value int64) string {
	*ptr = fmt.Sprint(value)
	return (*ptr).(string)
}

// ExpStrºStr adds two strings in string expression
func ExpStrºStr(left, right string) string {
	return left + right
}

// AddºStrStr adds two integer value
func AddºStrStr(left, right string) string {
	return left + right
}

// EqualºStrStr returns true if left == right
func EqualºStrStr(left, right string) bool {
	return left == right
}

// GreaterºStrStr returns true if left > right
func GreaterºStrStr(left, right string) bool {
	return left > right
}

// LenºStr returns the length of the string
func LenºStr(param string) int64 {
	return int64(len([]rune(param)))
}

// LessºStrStr returns true if left < right
func LessºStrStr(left, right string) bool {
	return left < right
}

// intºStr converts string value to int
func intºStr(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(val, 0, 64)
	if err != nil {
		err = errors.New(core.ErrorText(core.ErrStrToInt))
	}
	return
}

// floatºStr converts string value to float
func floatºStr(val string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(val, 64)
	if err != nil {
		err = errors.New(core.ErrorText(core.ErrStrToFloat))
	}
	return
}

// intºBool converts boolean value to int false -> 0, true -> 1
func boolºStr(val string) bool {
	return len(val) != 0 && val != `0` && strings.ToLower(val) != `false`
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
func HasPrefixºStrStr(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffixºStrStr returns true if the string s ends with suffix
func HasSuffixºStrStr(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// ReplaceºStrStrStr replaces strings in a string
func ReplaceºStrStrStr(in, old, new string) string {
	return strings.Replace(in, old, new, -1)
}

// LinesºStrArr splits a string to a array of strings
func LinesºStrArr(in string, out *core.Array) *core.Array {
	items := strings.Split(in, "\n")
	for _, item := range items {
		out.Data = append(out.Data, strings.Trim(item, "\r"))
	}
	return out
}

// SplitºStrStrArr splits a string to a array of strings
func SplitºStrStrArr(in, sep string, out *core.Array) *core.Array {
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
