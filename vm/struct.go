// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/gentee/gentee/core"
)

const (
	TRACESTRUCT = iota
	TIMESTRUCT
	FINFOSTRUCT
	HINFOSTRUCT
)

const (
	BININT = binary.MaxVarintLen64 + iota
	BINCHAR
	BINSTR
	BINFLOAT
	BINBUF
)

// Struct is used for custom struct types
type Struct struct {
	Type   *core.StructInfo
	Values []interface{} // Values of fields
}

// NewStruct creates a new struct object
func NewStruct(rt *Runtime, sInfo *core.StructInfo) *Struct {
	values := make([]interface{}, len(sInfo.Fields))
	for i, v := range sInfo.Fields {
		if v < core.TYPESTRUCT ||
			&rt.Owner.Exec.Structs[(v-core.TYPESTRUCT)>>8] != sInfo {
			values[i] = newValue(rt, int(v))
		}
	}
	return &Struct{
		Type:   sInfo,
		Values: values,
	}
}

// String interface for Struct
func (pstruct Struct) String() string {
	name := pstruct.Type.Name
	list := make([]string, len(pstruct.Values))
	for i, v := range pstruct.Values {
		list[i] = fmt.Sprintf(`%s:%v`, pstruct.Type.Keys[i], fmt.Sprint(v))
	}
	return name + `[` + strings.Join(list, ` `) + `]`
}

// Len is part of Indexer interface.
func (pstruct *Struct) Len() int {
	return len(pstruct.Values)
}

// GetIndex is part of Indexer interface.
func (pstruct *Struct) GetIndex(index interface{}) (interface{}, bool) {
	sindex := int(index.(int64))
	if sindex < 0 || sindex >= len(pstruct.Values) {
		return nil, false
	}
	return pstruct.Values[sindex], true
}

// SetIndex is part of Indexer interface.
func (pstruct *Struct) SetIndex(index, value interface{}) int {
	sindex := int(index.(int64))
	if sindex < 0 || sindex >= len(pstruct.Values) {
		return core.ErrIndexOut
	}
	pstruct.Values[sindex] = value
	return 0
}

// StructDecode decodes buffer to struct variable
func StructDecode(input *core.Buffer, pstruct *Struct) (err error) {
	var (
		num   int64
		dtype uint8
	)
	bint := make([]byte, 0, binary.MaxVarintLen64+1)
	buf := bytes.NewReader(input.Data)
	getInt := func(size uint8) (x int64, err error) {
		var (
			n int
		)
		if size == 0 {
			if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
				return 0, err
			}
		}
		if size > BININT {
			return 0, fmt.Errorf(`OOOPS`)
		}
		bint = bint[:size]
		if err = binary.Read(buf, binary.LittleEndian, &bint); err != nil {
			return 0, err
		}
		x, n = binary.Varint(bint)
		if n != int(size) {
			return 0, fmt.Errorf("Varint did not consume all of in")
		}
		return
	}
	for i, p := range pstruct.Type.Fields {
		var unsupported bool
		if dtype == 0 {
			if err = binary.Read(buf, binary.LittleEndian, &dtype); err != nil {
				return err
			}
		}
		switch p {
		case core.TYPEINT, core.TYPEBOOL, core.TYPECHAR:
			if num, err = getInt(dtype); err != nil {
				return err
			}
			pstruct.Values[i] = num
		case core.TYPESTR:
			if dtype != BINSTR {
				return fmt.Errorf("OOPS")
			}
			if num, err = getInt(0); err != nil {
				return err
			}
			tmp := make([]byte, num)
			if err = binary.Read(buf, binary.LittleEndian, &tmp); err != nil {
				return err
			}
			pstruct.Values[i] = string(tmp)
		case core.TYPEFLOAT:
			var f float64
			if dtype != BINFLOAT {
				return fmt.Errorf("OOPS")
			}
			if err = binary.Read(buf, binary.LittleEndian, &f); err != nil {
				return err
			}
			pstruct.Values[i] = f
		case core.TYPEBUF:
			if dtype != BINBUF {
				return fmt.Errorf("OOPS")
			}
			if num, err = getInt(0); err != nil {
				return err
			}
			tmp := make([]byte, num)
			if err = binary.Read(buf, binary.LittleEndian, &tmp); err != nil {
				return err
			}
			b := core.NewBuffer()
			b.Data = tmp
			pstruct.Values[i] = b
		default:
			unsupported = true
		}
		if !unsupported {
			dtype = 0
		}
	}
	return
}

// StructEncode encodes struct variable to buffer
func StructEncode(pstruct *Struct) (ret *core.Buffer, err error) {
	num := make([]byte, binary.MaxVarintLen64+1)
	buf := new(bytes.Buffer)

	putInt := func(i int64) error {
		n := binary.PutVarint(num[1:], i)
		num[0] = byte(n)
		return binary.Write(buf, binary.LittleEndian, num[:n+1])
	}
	for i, p := range pstruct.Type.Fields {
		var (
			data  interface{}
			dtype uint8
		)
		switch p {
		case core.TYPEINT, core.TYPEBOOL, core.TYPECHAR:
			if err = putInt(pstruct.Values[i].(int64)); err != nil {
				return
			}
		case core.TYPESTR:
			dtype = BINSTR
			data = []byte(pstruct.Values[i].(string))
		case core.TYPEFLOAT:
			dtype = BINFLOAT
			data = pstruct.Values[i].(float64)
		case core.TYPEBUF:
			dtype = BINBUF
			data = pstruct.Values[i].(*core.Buffer).Data
		}
		if dtype > BININT {
			if err = binary.Write(buf, binary.LittleEndian, dtype); err != nil {
				return
			}
			if dtype == BINBUF || dtype == BINSTR {
				if err = putInt(int64(len(data.([]byte)))); err != nil {
					return
				}
			}
			if err = binary.Write(buf, binary.LittleEndian, data); err != nil {
				return
			}
		}
	}
	ret = core.NewBuffer()
	ret.Data = buf.Bytes()
	return
}
