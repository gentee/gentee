// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"encoding/hex"
	"fmt"

	"github.com/gentee/gentee/core"
)

// InitBuffer appends stdlib buffer functions to the virtual machine
func InitBuffer(vm *core.VirtualMachine) {
	for _, item := range []embedInfo{
		{bufºStr, `str`, `buf`},                // buf( str ) buf
		{strºBuf, `buf`, `str`},                // str( buf ) str
		{LenºBuf, `buf`, `int`},                // the length of the buffer
		{AssignºBufBuf, `buf,buf`, `buf`},      // buf = buf
		{AssignAddºBufInt, `buf,int`, `buf`},   // buf += int
		{AssignAddºBufStr, `buf,str`, `buf`},   // buf += str
		{AssignAddºBufChar, `buf,char`, `buf`}, // buf += char
		{AssignAddºBufBuf, `buf,buf`, `buf`},   // buf += buf
		{HexºBuf, `buf`, `str`},                // Hex( buf ) str
		{UnHexºStr, `str`, `buf`},              // UnHex( str ) buf
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AssignºBufBuf copies one buf to another one
func AssignºBufBuf(ptr *interface{}, value *core.Buffer) *core.Buffer {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Buffer)
}

// AssignAddºBufChar appends rune to buffer
func AssignAddºBufChar(ptr *interface{}, value rune) *core.Buffer {
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, []byte(string([]rune{value}))...)
	return (*ptr).(*core.Buffer)
}

// AssignAddºBufInt appends one byte to buffer
func AssignAddºBufInt(ptr *interface{}, value int64) (*core.Buffer, error) {
	if uint64(value) > 255 {
		return nil, fmt.Errorf(core.ErrorText(core.ErrByteOut))
	}
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, byte(value))
	return (*ptr).(*core.Buffer), nil
}

// AssignAddºBufBuf appends buffer to buffer
func AssignAddºBufBuf(ptr *interface{}, value *core.Buffer) *core.Buffer {
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, value.Data...)
	return (*ptr).(*core.Buffer)
}

// AssignAddºBufStr appends string to buffer
func AssignAddºBufStr(ptr *interface{}, value string) *core.Buffer {
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, []byte(value)...)
	return (*ptr).(*core.Buffer)
}

// bufºStr converts string to buffer
func bufºStr(value string) *core.Buffer {
	b := core.NewBuffer()
	b.Data = []byte(value)
	return b
}

// strºBuf converts buffer to string
func strºBuf(buf *core.Buffer) string {
	return string(buf.Data)
}

// HexºBuf encodes buf to hex string
func HexºBuf(buf *core.Buffer) string {
	return hex.EncodeToString(buf.Data)
}

// LenºBuf returns the length of the buffer
func LenºBuf(buf *core.Buffer) int64 {
	return int64(len(buf.Data))
}

// UnHexºBufStr decodes hex string to the buffer
func UnHexºBufStr(buf *core.Buffer, value string) (*core.Buffer, error) {
	var err error
	buf.Data, err = hex.DecodeString(value)
	return buf, err
}

// UnHexºStr decodes hex string to the buffer
func UnHexºStr(value string) (*core.Buffer, error) {
	var err error
	buf := core.NewBuffer()
	buf.Data, err = hex.DecodeString(value)
	return buf, err
}
