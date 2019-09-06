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

// LinesºStr splits a string to a array of strings
func LinesºStr(in string) *core.Array {
	out := core.NewArray()
	items := strings.Split(in, "\n")
	for _, item := range items {
		out.Data = append(out.Data, strings.Trim(item, "\r"))
	}
	return out
}

// TrimSpaceºStr trims white space in a string
func TrimSpaceºStr(in string) string {
	return strings.TrimSpace(in)
}
