// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
)

const (
	// MaxSet is the maximum size of the set
	MaxSet = int64(64000000)
)

// Range is the type for operator ..
type Range struct {
	From int64
	To   int64
}

type Indexer interface {
	Len() int
	GetIndex(interface{}) (interface{}, bool)
	SetIndex(interface{}, interface{}) int
}

// KeyValue is the type for key value :
type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// Array is an array
type Array struct {
	Data []interface{}
}

// Buffer is []byte
type Buffer struct {
	Data []byte
}

// Set is []bool
type Set struct {
	Data []uint64
}

// Map is a map
type Map struct {
	Keys []string // it is required for 'for' statement and String interface
	Data map[string]interface{}
}

// FnType is used for func types
type FnType struct {
	Params []*TypeObject // Types of parameters
	Result *TypeObject   // Type of return value
}

// Fn is used for custom func types
type Fn struct {
	Func IObject
}

// StructType is used for custom struct types
type StructType struct {
	Fields map[string]int64 // Names of fields with indexes of the order
	Types  []*TypeObject    // Types of fields
}

// Struct is used for custom struct types
type Struct struct {
	Type   *TypeObject
	Values []interface{} // Values of fields
}

// Object is an object
type Obj struct {
	Data interface{}
}

// Len is part of Indexer interface.
func (prange *Range) Len() int {
	if prange.From < prange.To {
		return int(prange.To - prange.From + 1)
	}
	return int(prange.From - prange.To + 1)
}

// GetIndex is part of Indexer interface.
func (prange *Range) GetIndex(index interface{}) (interface{}, bool) {
	if prange.From < prange.To {
		return prange.From + index.(int64), true
	}
	return prange.From - index.(int64), true
}

// SetIndex is part of Indexer interface.
func (prange *Range) SetIndex(index, value interface{}) int {
	return ErrIndexOut
}

// String interface for Map
func (pmap Map) String() string {
	list := make([]string, len(pmap.Keys))
	for i, v := range pmap.Keys {
		list[i] = fmt.Sprintf(`%s:%v`, v, fmt.Sprint(pmap.Data[v]))
	}
	return `map[` + strings.Join(list, ` `) + `]`
}

// NewMap creates a new map object
func NewMap() *Map {
	return &Map{
		Data: make(map[string]interface{}),
		Keys: make([]string, 0),
	}
}

// Len is part of Indexer interface.
func (pmap *Map) Len() int {
	return len(pmap.Keys)
}

// GetIndex is part of Indexer interface.
func (pmap *Map) GetIndex(index interface{}) (interface{}, bool) {
	if key, ok := index.(string); ok {
		var (
			value interface{}
			ok    bool
		)
		if value, ok = pmap.Data[key]; !ok {
			return nil, false
		}
		return value, true
	}
	return pmap.Data[pmap.Keys[index.(int64)]], true
}

// SetIndex is part of Indexer interface.
func (pmap *Map) SetIndex(index, value interface{}) int {
	if v, ok := index.(int64); ok {
		pmap.Data[pmap.Keys[v]] = value
		return 0
	}
	sindex := index.(string)
	if _, ok := pmap.Data[sindex]; !ok {
		pmap.Data[sindex] = value
		pmap.Keys = append(pmap.Keys, sindex)
	} else {
		pmap.Data[sindex] = value
	}
	return 0
}

// String interface for Array
func (arr Array) String() string {
	return fmt.Sprint(arr.Data)
}

// NewArray creates a new array object
func NewArray() *Array {
	return &Array{
		Data: make([]interface{}, 0),
	}
}

// Len is part of sort.Interface.
func (arr *Array) Len() int {
	return len(arr.Data)
}

// Swap is part of sort.Interface.
func (arr *Array) Swap(i, j int) {
	arr.Data[i], arr.Data[j] = arr.Data[j], arr.Data[i]
}

// Less is part of sort.Interface.
func (arr *Array) Less(i, j int) bool {
	return arr.Data[i].(string) < arr.Data[j].(string)
}

// GetIndex is part of Indexer interface.
func (arr *Array) GetIndex(index interface{}) (interface{}, bool) {
	aindex := int(index.(int64))
	if aindex < 0 || aindex >= len(arr.Data) {
		return nil, false
	}
	return arr.Data[aindex], true
}

// SetIndex is part of Indexer interface.
func (arr *Array) SetIndex(index, value interface{}) int {
	aindex := int(index.(int64))
	if aindex < 0 || aindex >= len(arr.Data) {
		return ErrIndexOut
	}
	arr.Data[aindex] = value
	return 0
}

// String interface for Buffer
func (buf Buffer) String() string {
	return fmt.Sprint(buf.Data)
}

// NewBuffer creates a new buffer object
func NewBuffer() *Buffer {
	return &Buffer{
		Data: make([]byte, 0, 32),
	}
}

// Len is part of sort.Interface.
func (buf *Buffer) Len() int {
	return len(buf.Data)
}

// GetIndex is part of Indexer interface.
func (buf *Buffer) GetIndex(index interface{}) (interface{}, bool) {
	bindex := int(index.(int64))
	if bindex < 0 || bindex >= len(buf.Data) {
		return nil, false
	}
	return int64(buf.Data[bindex]), true
}

// SetIndex is part of Indexer interface.
func (buf *Buffer) SetIndex(index, value interface{}) int {
	bindex := int(index.(int64))
	if bindex < 0 || bindex >= len(buf.Data) {
		return ErrIndexOut
	}
	v := value.(int64)
	if uint64(v) > 255 {
		return ErrByteOut
	}
	buf.Data[bindex] = byte(v)
	return 0
}

// String interface for Set
func (set Set) String() string {
	var ret string
	for _, v := range set.Data {
		for pos := uint64(0); pos < 64; pos++ {
			if v&(1<<pos) == 0 {
				ret += `0`
			} else {
				ret += `1`
			}
		}
	}
	return strings.TrimRight(ret, `0`)
}

// IsSet returns the value of set[index]
func (set *Set) IsSet(index int64) bool {
	shift := int(index >> 6)
	pos := uint64(index % 64)
	if len(set.Data) <= shift || set.Data[shift]&(1<<pos) == 0 {
		return false
	}
	return true
}

// Set sets the value of set[index]
func (set *Set) Set(index int64, b bool) bool {
	shift := int(index >> 6)
	pos := uint64(index % 64)
	if len(set.Data) <= shift {
		set.Data = append(set.Data, make([]uint64, shift-len(set.Data)+1)...)
	}
	if b {
		set.Data[shift] |= 1 << pos
	} else {
		set.Data[shift] &= ^(1 << pos)
	}
	return b
}

// Len is part of sort.Interface.
func (set *Set) Len() int {
	return len(set.Data) << 6
}

// GetIndex is part of Indexer interface.
func (set *Set) GetIndex(index interface{}) (interface{}, bool) {
	sindex := int(index.(int64))
	shift := sindex >> 6
	if sindex < 0 || sindex >= int(MaxSet) { //len(set.Data) <= shift {
		return nil, false
	}
	if shift >= len(set.Data) {
		return int64(0), true
	}
	pos := uint64(sindex % 64)
	if set.Data[shift]&(1<<pos) == 0 {
		return int64(0), true
	}
	return int64(1), true
}

// SetIndex is part of Indexer interface.
func (set *Set) SetIndex(index, value interface{}) int {
	sindex := int64(index.(int64))
	if sindex < 0 {
		return ErrIndexOut
	}
	set.Set(sindex, value.(int64) == 1)
	return 0
}

// NewSet creates a new set object
func NewSet() *Set {
	return &Set{
		Data: make([]uint64, 1),
	}
}

/*
// NewStruct creates a new struct object
func NewStruct(ptype *TypeObject) *Struct {
	values := make([]interface{}, len(ptype.Custom.Types))
	for i, v := range ptype.Custom.Types {
		if v != ptype {
			values[i] = initVar(v)
		}
	}
	return &Struct{
		Type:   ptype,
		Values: values,
	}
}

// NewStructObj creates a stdlib struct object
func NewStructObj(rt *RunTime, name string) *Struct {
	return NewStruct(rt.VM.StdLib().FindType(name).(*TypeObject))
}*/

// String interface for Struct
func (pstruct Struct) String() string {
	name := pstruct.Type.GetName()
	keys := make([]string, len(pstruct.Values))
	list := make([]string, len(pstruct.Values))
	for key, ind := range pstruct.Type.Custom.Fields {
		keys[ind] = key
	}
	for i, v := range pstruct.Values {
		list[i] = fmt.Sprintf(`%s:%v`, keys[i], fmt.Sprint(v))
	}
	return name + `[` + strings.Join(list, ` `) + `]`
}

// NewFn creates a new func var
func NewFn(ptype *TypeObject) *Fn {
	return &Fn{
		Func: nil,
	}
}

// String interface for Obj
func (pobj Obj) String() string {
	if pobj.Data == nil {
		return `nil`
	}
	return fmt.Sprint(pobj.Data)
}

// NewObj creates a new object
func NewObj() *Obj {
	return &Obj{
		Data: nil,
	}
}
