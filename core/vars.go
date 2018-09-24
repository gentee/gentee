// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

// Range is the type for operator ..
type Range struct {
	From int64
	To   int64
}

// Array is an array
type Array struct {
	Data []interface{}
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
			index := rt.Stack[len(rt.Stack)-1].(int64)
			rt.Stack = rt.Stack[:len(rt.Stack)-1]
			switch typeValue.GetName() {
			case `str`:
				runes := []rune(value.(string))
				if index < 0 || index >= int64(len(runes)) {
					return runtimeError(rt, ival.Cmd, ErrIndexOut)
				}
				value = runes[index]
			case `arr`:
				var arr *Array
				arr = value.(*Array)
				if index < 0 || index >= int64(len(arr.Data)) {
					return runtimeError(rt, ival.Cmd, ErrIndexOut)
				}
				value = arr.Data[index]
			default:
				return runtimeError(rt, cmdVar, ErrRuntime, `getVar.default`)
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
			case `arr`:
				var arrPtr interface{}
				arrIndex = index.(int64)
				arr = (*ptr).(*Array)
				if arrIndex < 0 || arrIndex >= int64(len(arr.Data)) {
					return runtimeError(rt, ival.Cmd, ErrIndexOut)
				}
				arrPtr = arr.Data[arrIndex]
				ptr = &arrPtr
			default:
				return runtimeError(rt, cmdVar, ErrRuntime, `setVar.default`)
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
				value = reflect.New(cmdStack.Vars[i].Original).Elem().Interface()
			}
			if cmdStack.Vars[i].GetName() == `arr` {
				var arr Array
				arr = value.(Array)
				arr.Data = make([]interface{}, 0)
				value = &arr
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
