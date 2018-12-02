// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
	"strings"
)

// Range is the type for operator ..
type Range struct {
	From int64
	To   int64
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

// Map is a map
type Map struct {
	Keys []string // it is required for 'for' statement and String interface
	Data map[string]interface{}
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

func initVar(ptype *TypeObject) interface{} {
	var value interface{}
	if ptype.GetName() == `char` {
		value = ' '
	} else {
		original := ptype.Original
		switch original {
		case reflect.TypeOf(Buffer{}):
			value = NewBuffer()
		case reflect.TypeOf(Array{}):
			value = NewArray()
		case reflect.TypeOf(Map{}):
			value = NewMap()
		case reflect.TypeOf(Struct{}):
			value = NewStruct(ptype)
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
		strIndex                        int64
		prev                            *interface{}
		arr                             *Array
		buf                             *Buffer
		arrIndex, structIndex, bufIndex int64
		pmap                            *Map
		pstruct                         *Struct
		mapIndex                        string
	)
	value := rt.Stack[len(rt.Stack)-1]
	typeValue := cmdVar.Block.Vars[cmdVar.Index]
	if value == vars[cmdVar.Index] && (typeValue.Original == reflect.TypeOf(Struct{}) ||
		typeValue.Original == reflect.TypeOf(Array{}) || typeValue.Original == reflect.TypeOf(Buffer{}) ||
		typeValue.Original == reflect.TypeOf(Map{})) {
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
					structPtr = pstruct.Values[structIndex]
					ptr = &structPtr
				case reflect.TypeOf(Buffer{}):
					var bufPtr interface{}
					bufIndex = index.(int64)
					buf = (*ptr).(*Buffer)
					pmap = nil
					pstruct = nil
					arr = nil
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
					pmap = (*ptr).(*Map)
					if mapPtr, ok = pmap.Data[mapIndex]; !ok {
						pmap.Keys = append(pmap.Keys, mapIndex)
						if typeValue.IndexOf.Original == reflect.TypeOf(Array{}) {
							mapPtr = NewArray()
						} else if typeValue.IndexOf.Original == reflect.TypeOf(Map{}) {
							mapPtr = NewMap()
						} else if typeValue.IndexOf.Original == reflect.TypeOf(Buffer{}) {
							mapPtr = NewBuffer()
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
				typeValue.Original == reflect.TypeOf(Buffer{})) {
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
	rtBlock := RunTimeBlock{Block: cmdStack}
	if len(cmdStack.Vars) > 0 {
		for i := 0; i < cmdStack.ParCount; i++ {
			rtBlock.Vars = append(rtBlock.Vars, rt.Stack[len(rt.Stack)-cmdStack.ParCount+i])
			count++
		}
		rt.Stack = rt.Stack[:len(rt.Stack)-cmdStack.ParCount]
		for i := cmdStack.ParCount; i < len(cmdStack.Vars); i++ {
			rtBlock.Vars = append(rtBlock.Vars, initVar(cmdStack.Vars[i]))
		}
	}
	rt.Blocks = append(rt.Blocks, rtBlock)
	return
}

func deleteVars(rt *RunTime) {
	rt.Blocks = rt.Blocks[:len(rt.Blocks)-1]
}

// CopyVar copies one object to another one
func CopyVar(ptr *interface{}, value interface{}) {
	switch vItem := value.(type) {
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
