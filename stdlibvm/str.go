// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gentee/gentee/core"
)

// intºStr converts string value to int
func intºStr(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(val, 0, 64)
	if err != nil {
		err = errors.New(core.ErrorText(core.ErrStrToInt))
	}
	return
}

// AssignºStrBool assigns boolean to string
func AssignºStrBool(ptr *string, value int64) string {
	*ptr = strºBool(value)
	return *ptr
}

// AssignºStrInt assigns integer to string
func AssignºStrInt(ptr *string, value int64) string {
	*ptr = fmt.Sprint(value)
	return *ptr
}
