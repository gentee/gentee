// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/gentee/gentee/core"
)

// AddºBufBuf adds two buffers
func AddºBufBuf(left *core.Buffer, right *core.Buffer) (out *core.Buffer) {
	out = core.NewBuffer()
	out.Data = left.Data
	out.Data = append(out.Data, right.Data...)
	return out
}

// AssignAddºBufBuf appends buffer to buffer
func AssignAddºBufBuf(buf interface{}, value interface{}) (interface{}, error) {
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data, value.(*core.Buffer).Data...)
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
		return nil, fmt.Errorf(ErrorText(ErrByteOut))
	}
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data, byte(value.(int64)))
	return buf, nil
}

// AssignAddºBufStr appends string to buffer
func AssignAddºBufStr(buf interface{}, value interface{}) (interface{}, error) {
	buf.(*core.Buffer).Data = append(buf.(*core.Buffer).Data, []byte(value.(string))...)
	return buf, nil
}

// Base64ºBuf encodes buf to base64 string
func Base64ºBuf(buf *core.Buffer) string {
	return base64.StdEncoding.EncodeToString(buf.Data)
}

// bufºStr converts string to buffer
func bufºStr(value string) *core.Buffer {
	b := core.NewBuffer()
	b.Data = []byte(value)
	return b
}

// DelºBufIntInt deletes part of the buffer
func DelºBufIntInt(buf *core.Buffer, off, length int64) (*core.Buffer, error) {
	size := int64(len(buf.Data))
	if off < 0 || off > size {
		return buf, fmt.Errorf(core.ErrorText(core.ErrInvalidParam))
	}
	if length < 0 {
		off += length
		length = -length
	}
	if off < 0 {
		off = 0
	}
	if off+length > size {
		length = size - off
	}
	buf.Data = append(buf.Data[:off], buf.Data[off+length:]...)
	return buf, nil
}

// HexºBuf encodes buf to hex string
func HexºBuf(buf *core.Buffer) string {
	return hex.EncodeToString(buf.Data)
}

// InsertºBufIntBuf inserts one buf object into another one
func InsertºBufIntBuf(buf *core.Buffer, off int64, b *core.Buffer) (*core.Buffer, error) {
	size := int64(len(buf.Data))
	if off < 0 || off > size {
		return buf, fmt.Errorf(core.ErrorText(core.ErrInvalidParam))
	}
	buf.Data = append(buf.Data[:off], append(b.Data, buf.Data[off:]...)...)
	return buf, nil
}

// strºBuf converts buffer to string
func strºBuf(buf *core.Buffer) string {
	return string(buf.Data)
}

// sysBufNil return nil buffer
func sysBufNil() *core.Buffer {
	b := core.NewBuffer()
	b.Data = nil
	return b
}

// UnBase64ºStr decodes base64 string to buf
func UnBase64ºStr(value string) (buf *core.Buffer, err error) {
	buf = core.NewBuffer()
	buf.Data, err = base64.StdEncoding.DecodeString(value)
	return
}

// UnHexºStr decodes hex string to the buffer
func UnHexºStr(value string) (*core.Buffer, error) {
	var err error
	buf := core.NewBuffer()
	buf.Data, err = hex.DecodeString(value)
	return buf, err
}
