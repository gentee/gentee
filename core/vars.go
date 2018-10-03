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

// Map is a map
type Map struct {
	Keys []string // it is required for 'for' statement and String interface
	Data map[string]interface{}
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
				if typeValue.Original == reflect.TypeOf(Array{}) {
					var arr *Array
					arr = value.(*Array)
					if index < 0 || index >= int64(len(arr.Data)) {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					value = arr.Data[index]
				} else if typeValue.Original == reflect.TypeOf(Map{}) {
					var (
						pmap *Map
						ok   bool
					)
					pmap = value.(*Map)
					if value, ok = pmap.Data[mapIndex]; !ok {
						return runtimeError(rt, ival.Cmd, ErrMapIndex, mapIndex)
					}
				} else {
					return runtimeError(rt, cmdVar, ErrRuntime, `getVar.default`)
				}
			}
			typeValue = typeValue.IndexOf
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
		runes    []rune
		strIndex int64
		prev     *interface{}
		arr      *Array
		arrIndex int64
		pmap     *Map
		mapIndex string
	)
	if cmdVar.Indexes != nil {
		typeValue := cmdVar.Block.Vars[cmdVar.Index]
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
				if typeValue.Original == reflect.TypeOf(Array{}) {
					var arrPtr interface{}
					arrIndex = index.(int64)
					arr = (*ptr).(*Array)
					pmap = nil
					if arrIndex < 0 || arrIndex >= int64(len(arr.Data)) {
						return runtimeError(rt, ival.Cmd, ErrIndexOut)
					}
					arrPtr = arr.Data[arrIndex]
					ptr = &arrPtr
				} else if typeValue.Original == reflect.TypeOf(Map{}) {
					var (
						mapPtr interface{}
						ok     bool
					)
					mapIndex = index.(string)
					arr = nil
					pmap = (*ptr).(*Map)
					if mapPtr, ok = pmap.Data[mapIndex]; !ok {
						pmap.Keys = append(pmap.Keys, mapIndex)
						if typeValue.IndexOf.Original == reflect.TypeOf(Array{}) {
							mapPtr = NewArray()
						} else if typeValue.IndexOf.Original == reflect.TypeOf(Map{}) {
							mapPtr = NewMap()
						}
						pmap.Data[mapIndex] = mapPtr
					}
					ptr = &mapPtr
				} else {
					return runtimeError(rt, cmdVar, ErrRuntime, `setVar.default`)
				}
			}
			typeValue = typeValue.IndexOf
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
	if arr != nil {
		arr.Data[arrIndex] = *ptr
	}
	if pmap != nil {
		pmap.Data[mapIndex] = *ptr
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
			var value interface{}
			if cmdStack.Vars[i].GetName() == `char` {
				value = ' '
			} else {
				original := cmdStack.Vars[i].Original
				switch original {
				case reflect.TypeOf(Array{}):
					value = NewArray()
				case reflect.TypeOf(Map{}):
					value = NewMap()
				default:
					value = reflect.New(cmdStack.Vars[i].Original).Elem().Interface()
				}
			}
			rtBlock.Vars = append(rtBlock.Vars, value)
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
