// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"encoding/base64"
	"encoding/binary"
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

// boolºBuf converts buf to bool
func boolºBuf(buf *core.Buffer) int64 {
	if len(buf.Data) == 0 {
		return 0
	}
	return 1
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
		return buf, fmt.Errorf(ErrorText(ErrInvalidParam))
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

// EncodeºBufInt encodes int to buf
func EncodeºBufInt(buf *core.Buffer, i int64) *core.Buffer {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	buf.Data = append(buf.Data, b...)
	return buf
}

// DecodeºBufInt decodes int from buf
func DecodeºBufInt(buf *core.Buffer, offset int64) (int64, error) {
	if offset < 0 || offset+8 > int64(len(buf.Data)) {
		return 0, fmt.Errorf(ErrorText(ErrDecode))
	}
	return int64(binary.LittleEndian.Uint64(buf.Data[offset : offset+8])), nil
}

// HexºBuf encodes buf to hex string
func HexºBuf(buf *core.Buffer) string {
	return hex.EncodeToString(buf.Data)
}

// InsertºBufIntBuf inserts one buf object into another one
func InsertºBufIntBuf(buf *core.Buffer, off int64, b *core.Buffer) (*core.Buffer, error) {
	size := int64(len(buf.Data))
	if off < 0 || off > size {
		return buf, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	buf.Data = append(buf.Data[:off], append(b.Data, buf.Data[off:]...)...)
	return buf, nil
}

// SetLenºBuf sets the length of the buffer
func SetLenºBuf(buf *core.Buffer, size int64) (*core.Buffer, error) {
	if size < 0 {
		return buf, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	length := int64(len(buf.Data))
	if size < length {
		buf.Data = buf.Data[:size]
	} else if size > length {
		buf.Data = append(buf.Data, make([]byte, size-length)...)
	}
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

// WriteºBuf writes one buffer to another
func WriteºBuf(buf *core.Buffer, offset int64, input *core.Buffer) (*core.Buffer, error) {
	length := int64(len(buf.Data))
	ilen := int64(len(input.Data))
	if offset < 0 || offset > length {
		return buf, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	count := ilen
	if offset+ilen > length {
		count = length - offset
	}
	for i := int64(0); i < count; i++ {
		buf.Data[i+offset] = input.Data[i]
	}
	if count < ilen {
		buf.Data = append(buf.Data, input.Data[count:]...)
	}
	return buf, nil
}
