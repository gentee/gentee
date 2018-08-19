// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
)

// RunTimeBlock is a structure for storing variables
type RunTimeBlock struct {
	Vars  []interface{}
	Block *CmdBlock
}

// RunTime is the structure for running compiled functions
type RunTime struct {
	VM     *VirtualMachine
	Stack  []interface{} // the stack of values
	Calls  []ICmd        // the stack of calling functions
	Blocks []RunTimeBlock
	Result interface{} // result value
	Consts map[string]interface{}
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

func (rt *RunTime) callFunc(cmd ICmd) (err error) {
	var (
		result []reflect.Value
	)
	if int64(len(rt.Calls)) == rt.Consts[ConstDepth].(int64) {
		var (
			line, column int
			lex          *Lex
		)
		for i := len(rt.Calls) - 1; i >= 0 && lex == nil; i-- {
			if rt.Calls[i].GetObject() != nil {
				if lex = rt.Calls[i].GetObject().GetLex(); lex != nil {
					line, column = lex.LineColumn(cmd.GetToken())
					break
				}
			}
		}
		return fmt.Errorf(`%d:%d: %s`, line, column, ErrorText(ErrDepth))
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
		if len(result) > 1 && result[len(result)-1].Interface() != nil {
			var (
				line, column int
				lex          *Lex
			)
			for i := len(rt.Calls) - 1; i >= 0 && lex == nil; i-- {
				if rt.Calls[i].GetObject() == nil {
					continue
				}
				if lex = rt.Calls[i].GetObject().GetLex(); lex != nil {
					line, column = lex.LineColumn(cmd.GetToken())
					break
				}
			}
			return fmt.Errorf(`%d:%d: %s`, line, column, result[len(result)-1].Interface().(error))
		}
		rt.Stack = append(rt.Stack, result[0].Interface())
	case ObjFunc:
		if err = rt.runCmd(&cmd.GetObject().(*FuncObject).Block); err != nil {
			return
		}
	default:
		return runtimeError(rt, ErrRuntime)
	}
	return
}

func (rt *RunTime) runCmd(cmd ICmd) (err error) {
	rt.Calls = append(rt.Calls, cmd)
	switch cmd.GetType() {
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
		cmdVar := cmd.(*CmdVar)
		if !cmdVar.LValue {
			var i int
			for i = len(rt.Blocks) - 1; i >= 0; i-- {
				if rt.Blocks[i].Block == cmdVar.Block {
					rt.Stack = append(rt.Stack, rt.Blocks[i].Vars[cmdVar.Index])
					break
				}
			}
			if i < 0 {
				return runtimeError(rt, ErrRuntime)
			}
		}
	case CtStack:
		cmdStack := cmd.(*CmdBlock)
		lenStack := len(rt.Stack)
		switch cmd.(*CmdBlock).ID {
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
			var i int
			for i = len(rt.Blocks) - 1; i >= 0; i-- {
				if rt.Blocks[i].Block == cmdVar.Block {
					var post bool
					val := rt.Blocks[i].Vars[cmdVar.Index].(int64)
					shift := int64(cmdStack.ParCount)
					if (shift & 0x1) == 0 {
						post = true
						shift /= 2
					}
					rt.Blocks[i].Vars[cmdVar.Index] = val + shift
					if !post {
						val += shift
					}
					rt.Stack = append(rt.Stack, val)
					lenStack++
					break
				}
			}
		case StackAssign:
			if err = rt.runCmd(cmdStack.Children[1]); err != nil {
				return err
			}
			cmdVar := cmdStack.Children[0].(*CmdVar)
			var i int
			for i = len(rt.Blocks) - 1; i >= 0; i-- {
				if rt.Blocks[i].Block == cmdVar.Block {
					rt.Blocks[i].Vars[cmdVar.Index] = rt.Stack[len(rt.Stack)-1]
					break
				}
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
			for true {
				if err = rt.runCmd(cmdStack.Children[0]); err != nil {
					return err
				}
				if rt.Stack[len(rt.Stack)-1].(bool) {
					rt.Stack = rt.Stack[:len(rt.Stack)-1]
					if err = rt.runCmd(cmdStack.Children[1]); err != nil {
						return err
					}
					cycle--
					if cycle == 0 {
						var (
							line, column int
							lex          *Lex
						)
						for i := len(rt.Calls) - 1; i >= 0 && lex == nil; i-- {
							if rt.Calls[i].GetObject() != nil {
								if lex = rt.Calls[i].GetObject().GetLex(); lex != nil {
									line, column = lex.LineColumn(cmdStack.GetToken())
									break
								}
							}
						}
						return fmt.Errorf(`%d:%d: %s`, line, column, ErrorText(ErrCycle))
					}
					continue
				}
				break
			}
		case StackBlock:
			rt.Result = nil
			rtBlock := RunTimeBlock{Block: cmdStack}
			if len(cmdStack.Vars) > 0 {
				for i := 0; i < cmdStack.ParCount; i++ {
					rtBlock.Vars = append(rtBlock.Vars, rt.Stack[len(rt.Stack)-cmdStack.ParCount+i])
					lenStack--
				}
				rt.Stack = rt.Stack[:len(rt.Stack)-cmdStack.ParCount]
				for i := cmdStack.ParCount; i < len(cmdStack.Vars); i++ {
					rtBlock.Vars = append(rtBlock.Vars,
						reflect.New(cmdStack.Vars[i].Original).Elem().Interface())
				}
			}
			rt.Blocks = append(rt.Blocks, rtBlock)
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
			}
			rt.Blocks = rt.Blocks[:len(rt.Blocks)-1]
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
