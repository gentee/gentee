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
