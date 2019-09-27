// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"errors"
	"strconv"
	"strings"
)

// boolºStr converts string value to bool
func boolºStr(val string) int64 {
	if len(val) != 0 && val != `0` && strings.ToLower(val) != `false` {
		return 1
	}
	return 0
}

// floatºStr converts string value to float
func floatºStr(val string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(val, 64)
	if err != nil {
		err = errors.New(ErrorText(ErrStrToFloat))
	}
	return
}

// intºStr converts string value to int
func intºStr(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(val, 0, 64)
	if err != nil {
		err = errors.New(ErrorText(ErrStrToInt))
	}
	return
}
