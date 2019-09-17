// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
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
}

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

func initVar(ptype *TypeObject) interface{} {
	var value interface{}
	if ptype.GetName() == `char` {
		value = ' '
	} else {
		original := ptype.Original
		switch original {
		case reflect.TypeOf(Set{}):
			value = NewSet()
		case reflect.TypeOf(Buffer{}):
			value = NewBuffer()
		case reflect.TypeOf(Array{}):
			value = NewArray()
		case reflect.TypeOf(Map{}):
			value = NewMap()
		case reflect.TypeOf(Struct{}):
			value = NewStruct(ptype)
		case reflect.TypeOf(Fn{}):
			value = NewFn(ptype)
		default:
			value = reflect.New(ptype.Original).Elem().Interface()
		}
	}
	return value
}

func getVar(rt *RunTime, cmdVar *CmdVar) error {
	var (
		vars []interface{}
		err  error
	)

	if vars, err = rt.getVars(cmdVar.Block); err != nil {
		return err
	}
	value := vars[cmdVar.Index]
	if cmdVar.Indexes != nil {
		typeValue := cmdVar.Block.Vars[cmdVar.Index]
		for _, ival := range cmdVar.Indexes {
			if typeValue == nil {
				return runtimeError(rt, cmdVar, ErrRuntime, `getVar.typeValue`)
			}
			if err = rt.runCmd(ival.Cmd); err != nil {
				return err
			}
			var (
				index    int64
				mapIndex string
				custom   *Struct
			)
			if typeValue.Original == reflect.TypeOf(Map{}) {
				mapIndex = rt.Stack[len(rt.Stack)-1].(string)
			} else {
				index = rt.Stack[len(rt.Stack)-1].(int64)
			}
			rt.Stack = rt.Stack[:len(rt.Stack)-1]
			switch typeValue.GetName() {
			case `str`:
				runes := []rune(value.(string))
				if index < 0 || index >= int64(len(runes)) {
					return runtimeError(rt, ival.Cmd, ErrIndexOut)
				}
				value = runes[index]
			default:
				switch typeValue.Original {
				case reflect.TypeOf(Struct{}):
					custom = value.(*Struct)
					value = custom.Values[index]
					if value == nil {
						return runtimeError(rt, cmdVar, ErrUndefined)
					}
				case reflect.TypeOf(Set{}):
					var set *Set
					set = value.(*Set)
					if index < 0 || index >= MaxSet {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					value = set.IsSet(index)
				case reflect.TypeOf(Buffer{}):
					var buf *Buffer
					buf = value.(*Buffer)
					if index < 0 || index >= int64(len(buf.Data)) {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					value = int64(buf.Data[index])
				case reflect.TypeOf(Array{}):
					var arr *Array
					arr = value.(*Array)
					if index < 0 || index >= int64(len(arr.Data)) {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					value = arr.Data[index]
				case reflect.TypeOf(Map{}):
					var (
						pmap *Map
						ok   bool
					)
					pmap = value.(*Map)
					if value, ok = pmap.Data[mapIndex]; !ok {
						return runtimeError(rt, ival.Cmd, ErrMapIndex, mapIndex)
					}
				default:
					return runtimeError(rt, cmdVar, ErrRuntime, `getVar.default`)
				}
			}
			if typeValue.Original == reflect.TypeOf(Struct{}) {
				typeValue = custom.Type.Custom.Types[index]
			} else {
				typeValue = typeValue.IndexOf
			}
		}
	}
	rt.Stack = append(rt.Stack, value)
	return nil
}

func setVar(rt *RunTime, cmdStack *CmdBlock) error {
	var (
		vars []interface{}
		err  error
		ptr  *interface{}
	)
	if err = rt.runCmd(cmdStack.Children[1]); err != nil {
		return err
	}
	cmdVar := cmdStack.Children[0].(*CmdVar)
	if vars, err = rt.getVars(cmdVar.Block); err != nil {
		return err
	}
	ptr = &vars[cmdVar.Index]
	var (
		runes                           []rune
		strIndex, setIndex              int64
		prev                            *interface{}
		arr                             *Array
		buf                             *Buffer
		set                             *Set
		arrIndex, structIndex, bufIndex int64
		pmap                            *Map
		pstruct                         *Struct
		mapIndex                        string
	)
	value := rt.Stack[len(rt.Stack)-1]
	typeValue := cmdVar.Block.Vars[cmdVar.Index]
	if value == vars[cmdVar.Index] && (typeValue.Original == reflect.TypeOf(Struct{}) ||
		typeValue.Original == reflect.TypeOf(Array{}) || typeValue.Original == reflect.TypeOf(Buffer{}) ||
		typeValue.Original == reflect.TypeOf(Set{}) || typeValue.Original == reflect.TypeOf(Map{})) {
		return runtimeError(rt, cmdStack, ErrAssignment)
	}
	if cmdVar.Indexes != nil {
		for _, ival := range cmdVar.Indexes {
			if typeValue == nil {
				return runtimeError(rt, cmdVar, ErrRuntime, `setVar.typeValue`)
			}
			if err = rt.runCmd(ival.Cmd); err != nil {
				return err
			}
			index := rt.Stack[len(rt.Stack)-1]
			rt.Stack = rt.Stack[:len(rt.Stack)-1]
			switch typeValue.GetName() {
			case `str`:
				var strRune interface{}
				prev = ptr
				strIndex = index.(int64)
				runes = []rune((*ptr).(string))
				if strIndex < 0 || strIndex >= int64(len(runes)) {
					return runtimeError(rt, ival.Cmd, ErrIndexOut)
				}
				strRune = runes[strIndex]
				ptr = &strRune
			default:
				switch typeValue.Original {
				case reflect.TypeOf(Struct{}):
					var structPtr interface{}
					structIndex = index.(int64)
					pstruct = (*ptr).(*Struct)
					pmap = nil
					arr = nil
					buf = nil
					set = nil
					structPtr = pstruct.Values[structIndex]
					ptr = &structPtr
				case reflect.TypeOf(Set{}):
					var (
						setPtr interface{}
						vset   bool
					)

					setIndex = index.(int64)
					set = (*ptr).(*Set)
					pmap = nil
					pstruct = nil
					arr = nil
					buf = nil
					if setIndex < 0 || setIndex >= MaxSet {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					setPtr = vset
					ptr = &setPtr
				case reflect.TypeOf(Buffer{}):
					var bufPtr interface{}
					bufIndex = index.(int64)
					buf = (*ptr).(*Buffer)
					pmap = nil
					pstruct = nil
					arr = nil
					set = nil
					if bufIndex < 0 || bufIndex >= int64(len(buf.Data)) {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					bufPtr = buf.Data[bufIndex]
					ptr = &bufPtr
				case reflect.TypeOf(Array{}):
					var arrPtr interface{}
					arrIndex = index.(int64)
					arr = (*ptr).(*Array)
					pmap = nil
					pstruct = nil
					buf = nil
					set = nil
					if arrIndex < 0 || arrIndex >= int64(len(arr.Data)) {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					arrPtr = arr.Data[arrIndex]
					ptr = &arrPtr
				case reflect.TypeOf(Map{}):
					var (
						mapPtr interface{}
						ok     bool
					)
					mapIndex = index.(string)
					arr = nil
					pstruct = nil
					buf = nil
					set = nil
					pmap = (*ptr).(*Map)
					if mapPtr, ok = pmap.Data[mapIndex]; !ok {
						pmap.Keys = append(pmap.Keys, mapIndex)
						if typeValue.IndexOf.Original == reflect.TypeOf(Array{}) {
							mapPtr = NewArray()
						} else if typeValue.IndexOf.Original == reflect.TypeOf(Map{}) {
							mapPtr = NewMap()
						} else if typeValue.IndexOf.Original == reflect.TypeOf(Buffer{}) {
							mapPtr = NewBuffer()
						} else if typeValue.IndexOf.Original == reflect.TypeOf(Set{}) {
							mapPtr = NewSet()
						}
						pmap.Data[mapIndex] = mapPtr
					}
					ptr = &mapPtr
				default:
					return runtimeError(rt, cmdVar, ErrRuntime, `setVar.default`)
				}
			}
			if typeValue.Original == reflect.TypeOf(Struct{}) {
				typeValue = pstruct.Type.Custom.Types[structIndex]
			} else {
				typeValue = typeValue.IndexOf
			}
			if value == *ptr && (typeValue.Original == reflect.TypeOf(Struct{}) ||
				typeValue.Original == reflect.TypeOf(Array{}) ||
				typeValue.Original == reflect.TypeOf(Map{}) ||
				typeValue.Original == reflect.TypeOf(Buffer{}) ||
				typeValue.Original == reflect.TypeOf(Set{})) {
				return runtimeError(rt, cmdStack, ErrAssignment)
			}
		}
	}
	pars := []reflect.Value{reflect.ValueOf(ptr), reflect.ValueOf(rt.Stack[len(rt.Stack)-1])}
	result := reflect.ValueOf(cmdStack.GetObject().(*EmbedObject).Func).Call(pars)
	last := result[len(result)-1].Interface()
	if last != nil && reflect.TypeOf(last).String() == `*errors.errorString` {
		return runtimeError(rt, cmdStack, result[len(result)-1].Interface().(error))
	}
	if prev != nil {
		runes[strIndex] = result[0].Interface().(rune)
		*prev = string(runes)
		*ptr = *prev
	}
	if set != nil {
		set.Set(setIndex, result[0].Interface().(bool))
	}
	if buf != nil {
		val := result[0].Interface().(int64)
		if uint64(val) > 255 {
			return runtimeError(rt, cmdStack, ErrByteOut)
		}
		buf.Data[bufIndex] = byte(val)
	}
	if arr != nil {
		arr.Data[arrIndex] = *ptr
	}
	if pmap != nil {
		pmap.Data[mapIndex] = *ptr
	}
	if pstruct != nil {
		pstruct.Values[structIndex] = *ptr
	}
	rt.Stack[len(rt.Stack)-1] = result[0].Interface()
	return nil
}

func initVars(rt *RunTime, cmdStack *CmdBlock) (count int) {
	var variadicCount int

	rtBlock := RunTimeBlock{Block: cmdStack}
	if cmdStack.Variadic {
		variadicCount = rt.AllCount - cmdStack.ParCount
	}
	count = variadicCount + cmdStack.ParCount

	if len(cmdStack.Vars) > 0 {
		for i := 0; i < cmdStack.ParCount; i++ {
			rtBlock.Vars = append(rtBlock.Vars, rt.Stack[len(rt.Stack)-count+i])
		}
		for i := cmdStack.ParCount; i < len(cmdStack.Vars); i++ {
			var v interface{}
			if rt.Optional != nil && rt.Optional[i] != nil {
				v = rt.Optional[i]
				rtBlock.Optional = append(rtBlock.Optional, i)
			} else {
				v = initVar(cmdStack.Vars[i])
			}
			rtBlock.Vars = append(rtBlock.Vars, v)
		}
		if cmdStack.Variadic {
			aVar := rtBlock.Vars[cmdStack.ParCount].(*Array)
			for i := 0; i < variadicCount; i++ {
				value := rt.Stack[len(rt.Stack)-variadicCount+i]
				if v, ok := value.(*Array); ok &&
					cmdStack.Vars[cmdStack.ParCount].IndexOf.Original != reflect.TypeOf(Array{}) {
					aVar.Data = append(aVar.Data, v.Data...)
				} else {
					aVar.Data = append(aVar.Data, value)
				}
			}
		}
		rt.Stack = rt.Stack[:len(rt.Stack)-count]
	}
	rt.Optional = nil
	rt.Blocks = append(rt.Blocks, rtBlock)
	return
}

func deleteVars(rt *RunTime) {
	rt.Blocks = rt.Blocks[:len(rt.Blocks)-1]
}

// CopyVar copies one object to another one
func CopyVar(ptr *interface{}, value interface{}) {
	switch vItem := value.(type) {
	case *Fn:
		var pfn *Fn
		if ptr == nil || *ptr == nil {
			pfn = &Fn{}
		} else {
			pfn = (*ptr).(*Fn)
		}
		pfn.Func = vItem.Func
		*ptr = pfn
	case *Struct:
		var pstruct *Struct
		if ptr == nil || *ptr == nil {
			pstruct = NewStruct(vItem.Type)
		} else {
			pstruct = (*ptr).(*Struct)
		}
		pstruct.Values = make([]interface{}, len(vItem.Values))
		for i, v := range vItem.Values {
			CopyVar(&pstruct.Values[i], v)
		}
		*ptr = pstruct
	case *Set:
		var pset *Set
		if ptr == nil || *ptr == nil {
			pset = NewSet()
		} else {
			pset = (*ptr).(*Set)
		}
		pset.Data = make([]uint64, len(vItem.Data))
		copy(pset.Data, vItem.Data)
		*ptr = pset
	case *Buffer:
		var pbuf *Buffer
		if ptr == nil || *ptr == nil {
			pbuf = NewBuffer()
		} else {
			pbuf = (*ptr).(*Buffer)
		}
		pbuf.Data = make([]byte, len(vItem.Data))
		copy(pbuf.Data, vItem.Data)
		*ptr = pbuf
	case *Array:
		var parr *Array
		if ptr == nil || *ptr == nil {
			parr = NewArray()
		} else {
			parr = (*ptr).(*Array)
		}
		parr.Data = make([]interface{}, len(vItem.Data))
		for i, v := range vItem.Data {
			CopyVar(&parr.Data[i], v)
		}
		*ptr = parr
	case *Map:
		var pmap *Map
		if ptr == nil || *ptr == nil {
			pmap = NewMap()
		} else {
			pmap = (*ptr).(*Map)
		}
		pmap.Keys = make([]string, len(vItem.Keys))
		var mapPtr interface{}
		for i, v := range vItem.Keys {
			mapPtr = pmap.Data[v]
			pmap.Keys[i] = v
			CopyVar(&mapPtr, vItem.Data[v])
			pmap.Data[v] = mapPtr
		}
		*ptr = pmap
	default:
		*ptr = value
	}
}
