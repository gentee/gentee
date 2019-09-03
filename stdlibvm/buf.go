// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// AssignAddºBufBuf appends buffer to buffer
func AssignAddºBufBuf(buf interface{}, value interface{}) (interface{}, error) {
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data, value.(*core.Buffer).Data...)
	return buf, nil
}

// AssignAddºBufStr appends string to buffer
func AssignAddºBufStr(buf interface{}, value interface{}) (interface{}, error) {
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data, []byte(value.(string))...)
	return buf, nil
}

// AssignAddºBufChar appends rune to buffer
func AssignAddºBufChar(buf interface{}, value interface{}) (interface{}, error) {
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data,
		[]byte(string([]rune{rune(value.(int64))}))...)
	return buf, nil
}

// AssignAddºBufInt appends one byte to buffer
func AssignAddºBufInt(buf interface{}, value interface{}) (interface{}, error) {
	if uint64(value.(int64)) > 255 {
		return nil, fmt.Errorf(core.ErrorText(core.ErrByteOut))
	}
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data, byte(value.(int64)))
	return buf, nil
}
