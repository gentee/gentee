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
func InitStr(ws *core.Workspace) {
	for _, item := range []interface{}{
		core.Link{AddºStrStr, core.Bcode(core.TYPESTR<<16) | core.ADDSTR}, // binary +
		core.Link{EqualºStrStr, core.EQSTR},                               // binary ==
		core.Link{GreaterºStrStr, core.GTSTR},                             // binary >
		core.Link{LenºStr, core.Bcode(core.TYPESTR<<16) | core.LEN},       // the length of str
		core.Link{LessºStrStr, core.LTSTR},                                // binary <
		core.Link{intºStr, 2<<16 | core.EMBED},                            // int( str )
		core.Link{floatºStr, 20<<16 | core.EMBED},                         // float( str )
		boolºStr,                                                           // bool( str )
		core.Link{ExpStrºStr, core.ADDSTR},                                 // expression in string
		core.Link{AssignºStrStr, core.ASSIGN},                              // str = str
		core.Link{AssignAddºStrStr, core.ASSIGN + 1},                       // str += str
		core.Link{AssignºStrBool, core.ASSIGN + 2 /*12<<16 | core.EMBED*/}, // str = bool
		core.Link{AssignºStrInt, core.ASSIGN + 3 /*13<<16 | core.EMBED*/},  // str = int
		FindºStrStr,       // Find( str, str ) int
		FormatºStr,        // Format( str, ... ) str
		HasPrefixºStrStr,  // HasPrefix( str, str ) bool
		HasSuffixºStrStr,  // HasSuffix( str, str ) bool
		LeftºStrInt,       // Left( str, int ) str
		LowerºStr,         // Lower( str ) str
		RepeatºStrInt,     // Repeat( str, int )
		ReplaceºStrStrStr, // Replace( str, str, str )
		ShiftºStr,         // unary bitwise OR
		SubstrºStrIntInt,  // Substr( str, int, int ) str
		core.Link{TrimSpaceºStr, 35<<16 | core.EMBED}, // TrimSpace( str ) str
		TrimRightºStr, // TrimRight( str, str ) str
		UpperºStr,     // Upper( str ) str
	} {
		ws.StdLib().NewEmbed(item)
	}

	for _, item := range []embedInfo{
		{core.Link{LinesºStr, 36<<16 | core.EMBED}, `str`, `arr.str`}, // Lines( str ) arr
		{SplitºStrStr, `str,str`, `arr.str`},                          // Split( str, str ) arr
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
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

// SplitºStrStr splits a string to a array of strings
func SplitºStrStr(in, sep string) *core.Array {
	out := core.NewArray()
	items := strings.Split(in, sep)
	for _, item := range items {
		out.Data = append(out.Data, item)
	}
	return out
}

// ShiftºStr trims white spaces characters in the each line of the string.
func ShiftºStr(par string) string {
	lines := strings.Split(par, "\n")
	for i, v := range lines {
		lines[i] = strings.TrimSpace(v)
	}
	return strings.Join(lines, "\n")
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

// TrimRightºStr trims white space in a string
func TrimRightºStr(in string, set string) string {
	return strings.TrimRight(in, set)
}

// UpperºStr converts a copy of the string to their upper case and returns it.
func UpperºStr(s string) string {
	return strings.ToUpper(s)
}
