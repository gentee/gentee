// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

// RunTimeBlock is a structure for storing variables
type RunTimeBlock struct {
	Vars  []interface{}
	Block *CmdBlock
}

// RunTime is the structure for running compiled functions
type RunTime struct {
	VM      *VirtualMachine
	Stack   []interface{} // the stack of values
	Calls   []ICmd        // the stack of calling functions
	Blocks  []RunTimeBlock
	Result  interface{} // result value
	Command uint32
	Consts  map[string]interface{}
}

func newRunTime(vm *VirtualMachine) *RunTime {
	rt := &RunTime{
		VM:     vm,
		Stack:  make([]interface{}, 0, 1024),
		Calls:  make([]ICmd, 0, 64),
		Consts: make(map[string]interface{}),
	}

	for _, item := range []string{ConstDepth, ConstCycle} {
		// TODO: Insert redefinition of constants here
		rt.runCmd(vm.StdLib().Names[item].(*ConstObject).Exp)
		rt.Consts[item] = rt.Stack[len(rt.Stack)-1]
	}
	return rt
}

func (rt *RunTime) getVars(block *CmdBlock) ([]interface{}, error) {
	var i int
	for i = len(rt.Blocks) - 1; i >= 0; i-- {
		if rt.Blocks[i].Block == block {
			return rt.Blocks[i].Vars, nil
		}
	}
	return nil, runtimeError(rt, block, ErrRuntime, `getVars`)
}

func (rt *RunTime) callFunc(cmd ICmd) (err error) {
	var (
		result []reflect.Value
	)
	if int64(len(rt.Calls)) == rt.Consts[ConstDepth].(int64) {
		return runtimeError(rt, cmd, ErrDepth)
	}
	pars := make([]reflect.Value, 0)
	lenStack := len(rt.Stack)
	switch cmd.GetType() {
	case CtFunc:
		for _, param := range cmd.(*CmdAnyFunc).Children {
			if err = rt.runCmd(param); err != nil {
				return
			}
		}
	case CtBinary:
		if err = rt.runCmd(cmd.(*CmdBinary).Left); err != nil {
			return
		}
		if err = rt.runCmd(cmd.(*CmdBinary).Right); err != nil {
			return
		}
	case CtUnary:
		if err = rt.runCmd(cmd.(*CmdUnary).Operand); err != nil {
			return
		}
	}
	switch cmd.GetObject().GetType() {
	case ObjEmbedded:
		for i := lenStack; i < len(rt.Stack); i++ {
			pars = append(pars, reflect.ValueOf(rt.Stack[i]))
		}
		rt.Stack = rt.Stack[:lenStack]
		result = reflect.ValueOf(cmd.GetObject().(*EmbedObject).Func).Call(pars)
		last := result[len(result)-1].Interface()
		if last != nil {
			if _, isError := last.(error); isError {
				return runtimeError(rt, cmd, result[len(result)-1].Interface().(error))
			}
		}
		rt.Stack = append(rt.Stack, result[0].Interface())
	case ObjFunc:
		if err = rt.runCmd(&cmd.GetObject().(*FuncObject).Block); err != nil {
			return
		}
	default:
		return runtimeError(rt, cmd, ErrRuntime, `callFunc`)
	}
	return
}

func (rt *RunTime) runCmd(cmd ICmd) (err error) {
	var vars []interface{}

	rt.Calls = append(rt.Calls, cmd)
	switch cmd.GetType() {
	case CtCommand:
		rt.Command = cmd.(*CmdCommand).ID
	case CtFunc, CtBinary, CtUnary:
		err = rt.callFunc(cmd)
	case CtValue:
		rt.Stack = append(rt.Stack, cmd.(*CmdValue).Value)
	case CtConst:
		name := cmd.GetObject().GetName()
		if v, ok := rt.Consts[name]; ok {
			rt.Stack = append(rt.Stack, v)
		} else {
			// TODO: Insert redefinition of constants here
			constObj := cmd.GetObject().(*ConstObject)
			if constObj.Iota != NotIota {
				rt.Consts[ConstIota] = constObj.Iota
			}
			if err = rt.runCmd(constObj.Exp); err != nil {
				return err
			}
			rt.Consts[name] = rt.Stack[len(rt.Stack)-1]
		}
	case CtVar:
		if err = getVar(rt, cmd.(*CmdVar)); err != nil {
			return err
		}
	case CtStack:
		cmdStack := cmd.(*CmdBlock)
		lenStack := len(rt.Stack)
		switch cmd.(*CmdBlock).ID {
		case StackNew:
			switch cmd.(*CmdBlock).Result.Original {
			case reflect.TypeOf(Array{}):
				parr := NewArray()
				for _, icmd := range cmdStack.Children {
					if err = rt.runCmd(icmd); err != nil {
						return err
					}
					var ptr interface{}
					CopyVar(&ptr, rt.Stack[len(rt.Stack)-1])
					parr.Data = append(parr.Data, ptr)
				}
				if lenStack >= len(rt.Stack) {
					rt.Stack = append(rt.Stack, parr)
				} else {
					rt.Stack[lenStack] = parr
				}
			case reflect.TypeOf(Buffer{}):
				pbuf := NewBuffer()
				for _, icmd := range cmdStack.Children {
					if err = rt.runCmd(icmd); err != nil {
						return err
					}
					switch v := rt.Stack[len(rt.Stack)-1].(type) {
					case int64:
						if uint64(v) > 255 {
							return runtimeError(rt, icmd, ErrByteOut)
						}
						pbuf.Data = append(pbuf.Data, byte(v))
					case string:
						pbuf.Data = append(pbuf.Data, []byte(v)...)
					case rune:
						pbuf.Data = append(pbuf.Data, []byte(string([]rune{v}))...)
					case *Buffer:
						pbuf.Data = append(pbuf.Data, v.Data...)
					default:
						return runtimeError(rt, icmd, ErrRuntime, `init buf`)
					}
				}
				rt.Stack[lenStack] = pbuf
			case reflect.TypeOf(Map{}):
				pmap := NewMap()
				for _, icmd := range cmdStack.Children {
					if err = rt.runCmd(icmd); err != nil {
						return err
					}
					var ptr interface{}
					CopyVar(&ptr, rt.Stack[len(rt.Stack)-1])
					keyValue := ptr.(KeyValue)
					pmap.Data[keyValue.Key.(string)] = keyValue.Value
					pmap.Keys = append(pmap.Keys, keyValue.Key.(string))
				}
				if lenStack >= len(rt.Stack) {
					rt.Stack = append(rt.Stack, pmap)
				} else {
					rt.Stack[lenStack] = pmap
				}
			case reflect.TypeOf(Struct{}):
				pstruct := NewStruct(cmd.(*CmdBlock).Result)
				for _, icmd := range cmdStack.Children {
					if err = rt.runCmd(icmd); err != nil {
						return err
					}
					var ptr interface{}
					CopyVar(&ptr, rt.Stack[len(rt.Stack)-1])
					keyValue := ptr.(KeyValue)
					pstruct.Values[keyValue.Key.(int64)] = keyValue.Value
				}
				if lenStack >= len(rt.Stack) {
					rt.Stack = append(rt.Stack, pstruct)
				} else {
					rt.Stack[lenStack] = pstruct
				}
			default:
				return runtimeError(rt, cmd, ErrRuntime, `init arr`)
			}
			lenStack++
		case StackInit:
			cmdVar := cmdStack.Children[0].(*CmdVar)
			if vars, err = rt.getVars(cmdVar.Block); err != nil {
				return err
			}
			if err = rt.runCmd(cmdStack.Children[1]); err != nil {
				return err
			}
			vars[cmdVar.Index] = rt.Stack[len(rt.Stack)-1]
		case StackQuestion:
			if err = rt.runCmd(cmdStack.Children[0]); err != nil {
				return err
			}
			iExp := 2
			if rt.Stack[len(rt.Stack)-1].(bool) {
				iExp = 1
			}
			if err = rt.runCmd(cmdStack.Children[iExp]); err != nil {
				return err
			}
			rt.Stack[lenStack] = rt.Stack[len(rt.Stack)-1]
			lenStack++
		case StackAnd, StackOr:
			if err = rt.runCmd(cmdStack.Children[0]); err != nil {
				return err
			}
			if (rt.Stack[len(rt.Stack)-1].(bool) && cmd.(*CmdBlock).ID == StackAnd) ||
				(!rt.Stack[len(rt.Stack)-1].(bool) && cmd.(*CmdBlock).ID == StackOr) {
				if err = rt.runCmd(cmdStack.Children[1]); err != nil {
					return err
				}
			}
			rt.Stack[lenStack] = rt.Stack[len(rt.Stack)-1]
			lenStack++
		case StackIncDec:
			cmdVar := cmdStack.Children[0].(*CmdVar)
			if vars, err = rt.getVars(cmdVar.Block); err != nil {
				return err
			}
			var post bool
			if err = getVar(rt, cmdVar); err != nil {
				return err
			}
			val := rt.Stack[len(rt.Stack)-1].(int64)
			rt.Stack = rt.Stack[:len(rt.Stack)-1]
			shift := int64(cmdStack.ParCount)
			if (shift & 0x1) == 0 {
				post = true
				shift /= 2
			}
			if err = setVar(rt, &CmdBlock{Children: []ICmd{
				cmdVar, &CmdValue{Value: val + shift},
			}, Object: rt.VM.StdLib().Names[DefAssignIntInt]}); err != nil {
				return err
			}
			rt.Stack = rt.Stack[:len(rt.Stack)-1]
			if !post {
				val += shift
			}
			rt.Stack = append(rt.Stack, val)
			lenStack++
		case StackAssign:
			if err = setVar(rt, cmdStack); err != nil {
				return err
			}
			lenStack++
		case StackIf:
			var i int
			lenIf := len(cmdStack.Children) >> 1
			for i = 0; i < lenIf; i++ {
				if err = rt.runCmd(cmdStack.Children[i<<1]); err != nil {
					return err
				}
				if rt.Stack[len(rt.Stack)-1].(bool) {
					if err = rt.runCmd(cmdStack.Children[(i<<1)+1]); err != nil {
						return err
					}
					break
				}
			}
			// Calling else
			if i == lenIf && len(cmdStack.Children)&1 == 1 {
				if err = rt.runCmd(cmdStack.Children[len(cmdStack.Children)-1]); err != nil {
					return err
				}
			}
		case StackWhile:
			cycle := rt.Consts[ConstCycle].(int64)
			for rt.Result == nil {
				if err = rt.runCmd(cmdStack.Children[0]); err != nil {
					return err
				}
				if rt.Stack[len(rt.Stack)-1].(bool) {
					rt.Stack = rt.Stack[:len(rt.Stack)-1]
					if err = rt.runCmd(cmdStack.Children[1]); err != nil {
						return err
					}
					if rt.Command == RcBreak {
						rt.Command = 0
						break
					}
					if rt.Command == RcContinue {
						rt.Command = 0
					}
					cycle--
					if cycle == 0 {
						return runtimeError(rt, cmdStack, ErrCycle)
					}
					continue
				}
				break
			}
		case StackFor:
			if err = rt.runCmd(cmdStack.Children[0]); err != nil {
				return err
			}
			value := rt.Stack[len(rt.Stack)-1]
			rt.Stack = rt.Stack[:len(rt.Stack)-1]
			var index int64
			length := getLength(value)
			lenStack -= initVars(rt, cmdStack)
			if vars, err = rt.getVars(cmdStack); err != nil {
				return err
			}
			cycle := rt.Consts[ConstCycle].(int64)
			for ; index < length; index++ {
				vars[0] = getIndex(value, index)
				vars[1] = index
				if err = rt.runCmd(cmdStack.Children[1]); err != nil {
					return err
				}
				if rt.Result != nil {
					break
				}
				if rt.Command == RcBreak {
					rt.Command = 0
					break
				}
				if rt.Command == RcContinue {
					rt.Command = 0
				}
				length = getLength(value)
				if index > cycle {
					return runtimeError(rt, cmdStack, ErrCycle)
				}
			}
			deleteVars(rt)
		case StackBlock:
			rt.Result = nil
			lenStack -= initVars(rt, cmdStack)
			for _, item := range cmdStack.Children {
				if err = rt.runCmd(item); err != nil {
					return err
				}
				if rt.Result != nil {
					if cmdStack.Result != nil {
						rt.Stack = rt.Stack[:lenStack]
						rt.Stack = append(rt.Stack, rt.Result)
						lenStack++
						rt.Result = nil
					}
					break
				}
				if rt.Command != 0 {
					break
				}
			}
			deleteVars(rt)
		case StackReturn:
			if cmdStack.Children != nil {
				if err = rt.runCmd(cmdStack.Children[0]); err != nil {
					return err
				}
				rt.Result = rt.Stack[len(rt.Stack)-1]
			} else { // return from the function without result value
				rt.Result = true
			}
		}
		rt.Stack = rt.Stack[:lenStack]
	}
	if err == nil {
		rt.Calls = rt.Calls[:len(rt.Calls)-1]
	}
	return err
}

func getLength(value interface{}) (length int64) {
	switch reflect.TypeOf(value).String() {
	case `string`:
		length = int64(len([]rune(value.(string))))
	case `core.Range`:
		rangeVal := value.(Range)
		length = rangeVal.To - rangeVal.From
		if length < 0 {
			length = -length
		}
		length++
	case `*core.Buffer`:
		length = int64(len(value.(*Buffer).Data))
	case `*core.Array`:
		length = int64(len(value.(*Array).Data))
	case `*core.Map`:
		length = int64(len(value.(*Map).Keys))
	}
	return
}

func getIndex(value interface{}, index int64) interface{} {
	switch reflect.TypeOf(value).String() {
	case `string`:
		return []rune(value.(string))[index]
	case `core.Range`:
		rangeVal := value.(Range)
		if rangeVal.From < rangeVal.To {
			return rangeVal.From + index
		}
		return rangeVal.From - index
	case `*core.Buffer`:
		return int64(value.(*Buffer).Data[index])
	case `*core.Array`:
		return value.(*Array).Data[index]
	case `*core.Map`:
		return value.(*Map).Data[value.(*Map).Keys[index]]
	}
	return nil
}
